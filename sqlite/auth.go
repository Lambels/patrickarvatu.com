package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	pa "github.com/Lambels/patrickarvatu.com"
)

// check to see if *AuthService object implements set interface.
var _ pa.AuthService = (*AuthService)(nil)

// AuthService represents a service used to manage OAuth.
type AuthService struct {
	db *DB
}

// NewAuthService returns a new instance of AuthService attached to db.
func NewAuthService(db *DB) *AuthService {
	return &AuthService{
		db: db,
	}
}

// FindAuthByID returns a auth based on the id.
// returns ENOTFOUND if the auth doesent exist.
func (s *AuthService) FindAuthByID(ctx context.Context, id int) (*pa.Auth, error) {
	tx, err := s.db.BeginTX(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	auth, err := findAuthByID(ctx, tx, id) // find auth obj
	if err != nil {
		return nil, err

	} else if err := attachUserToAuth(ctx, tx, auth); err != nil { // attach user obj to auth obj
		return nil, err
	}

	return auth, nil
}

// FindAuths returns a range of auth based on filter.
func (s *AuthService) FindAuths(ctx context.Context, filter pa.AuthFilter) ([]*pa.Auth, int, error) {
	tx, err := s.db.BeginTX(ctx, nil)
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback()

	auths, n, err := findAuths(ctx, tx, filter)
	if err != nil {
		return auths, n, err
	}

	// loops like this work good for SQLite database but when using a remote database
	// buffer queries to avoid high latency time loss.
	for _, auth := range auths {
		if err := attachUserToAuth(ctx, tx, auth); err != nil {
			return auths, n, err
		}
	}

	return auths, n, nil
}

// CreateAuth creates a new auth obj if a user is attached, the auth obj is linked back
// to only an existing user, if not existing the user is created and the auth is attached.
// A sucessful call will return an auth with the auth.UserID != 0
func (s *AuthService) CreateAuth(ctx context.Context, auth *pa.Auth) error {
	tx, err := s.db.BeginTX(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// check if the auth exists with the same source
	if other, err := findAuthBySourceID(ctx, tx, auth.Source, auth.SourceID); err == nil {

		// if found update the found one to new requirements
		if other, err = updateAuth(ctx, tx, other.ID, auth.AccessToken, auth.RefreshToken, auth.Expiry); err != nil {
			return fmt.Errorf("updateAuth: err=%w id=%d", err, other.ID)
		} else if err := attachUserToAuth(ctx, tx, other); err != nil { // attach user ob to auth obj
			return err
		}

		// refrence the other auth obj (updated) to the caller auth
		*auth = *other
		return tx.Commit()
	} else if code := pa.ErrorCode(err); code != pa.ENOTFOUND {
		return fmt.Errorf("cant find auth by source id: %w", err)
	}

	// the ID set to 0 indicates the creation of a new user under auth.User
	// the existance of the auth.User indicates that we want to attach the auth to the user
	if auth.UserID == 0 && auth.User != nil {

		if user, err := findUserByEmail(ctx, tx, auth.User.Email); err != nil {
			auth.User = user // user exists so we attach

		} else if pa.ErrorCode(err) == pa.ENOTFOUND {
			if err := createUser(ctx, tx, auth.User); err != nil {
				return fmt.Errorf("createUser: err=%w email=%s", err, auth.User.Email)
			}

		} else {
			return fmt.Errorf("findUserByEmail: err=%w", err)
		}

		// attach the user id to the newly created user
		auth.ID = auth.User.ID
	}

	if err := createAuth(ctx, tx, auth); err != nil {
		return err

	} else if err := attachUserToAuth(ctx, tx, auth); err != nil {
		return err
	}
	return tx.Commit()
}

// DeleteAuth permanently deletes the auth specified by the id.
// attached user wont be removed.
func (s *AuthService) DeleteAuth(ctx context.Context, id int) error {
	tx, err := s.db.BeginTX(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := deleteAuth(ctx, tx, id); err != nil {
		return err
	}
	return tx.Commit()
}

func findAuthByID(ctx context.Context, tx *Tx, id int) (*pa.Auth, error) {
	filter := pa.AuthFilter{
		ID: &id,
	}

	auths, _, err := findAuths(ctx, tx, filter)
	if err != nil {
		return nil, err

	} else if len(auths) == 0 {
		return nil, pa.Errorf(pa.ENOTFOUND, "auth not found.")
	}
	return auths[0], nil
}

func findAuthBySourceID(ctx context.Context, tx *Tx, source string, sourceID string) (*pa.Auth, error) {
	filter := pa.AuthFilter{
		Source:   &source,
		SourceID: &sourceID,
	}

	auths, _, err := findAuths(ctx, tx, filter)
	if err != nil {
		return nil, err

	} else if len(auths) == 0 {
		return nil, pa.Errorf(pa.ENOTFOUND, "auth not found.")
	}
	return auths[0], nil
}

func findAuths(ctx context.Context, tx *Tx, filter pa.AuthFilter) (_ []*pa.Auth, n int, err error) {
	// build where and args statement method.
	// not vulnerable to sql injection attack.
	where, args := []string{"1 = 1"}, []interface{}{}

	if v := filter.ID; v != nil {
		where = append(where, "id = ?")
		args = append(args, *v)
	}
	if v := filter.UserID; v != nil {
		where = append(where, "user_id = ?")
		args = append(args, *v)
	}
	if v := filter.Source; v != nil {
		where = append(where, "source = ?")
		args = append(args, *v)
	}
	if v := filter.SourceID; v != nil {
		where = append(where, "source_id = ?")
		args = append(args, *v)
	}

	rows, err := tx.QueryContext(ctx, `
		SELECT 
		    id,
		    user_id,
		    source,
		    source_id,
		    access_token,
		    refresh_token,
		    expiry,
		    created_at,
		    updated_at,
		    COUNT(*) OVER()
		FROM auths
		WHERE `+strings.Join(where, " AND ")+`
		ORDER BY id ASC
		`+FormatLimitOffset(filter.Limit, filter.Offset)+`
	`,
		args...,
	)

	if err != nil {
		return nil, n, err
	}
	defer rows.Close()

	// deserialize rows.
	auths := []*pa.Auth{}
	for rows.Next() {
		var auth *pa.Auth
		var expiry sql.NullString

		if err := rows.Scan(
			&auth.ID,
			&auth.UserID,
			&auth.Source,
			&auth.SourceID,
			&auth.AccessToken,
			&auth.RefreshToken,
			&expiry,
			(*NullTime)(&auth.CreatedAt),
			(*NullTime)(&auth.UpdatedAt),
			&n,
		); err != nil {
			return nil, 0, err
		}

		// different providers differ in providing expiry so we validate its existance and attach
		// it to auth if valid.
		if expiry.Valid {
			if v, _ := time.Parse(time.RFC3339, expiry.String); !v.IsZero() {
				auth.Expiry = &v
			}
		}

		auths = append(auths, auth)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return auths, n, nil
}

func createAuth(ctx context.Context, tx *Tx, auth *pa.Auth) error {
	auth.CreatedAt = tx.now
	auth.UpdatedAt = auth.CreatedAt

	if err := auth.Validate(); err != nil {
		return err
	}

	var expiry *string
	if auth.Expiry != nil {
		v := auth.Expiry.Format(time.RFC3339)
		expiry = &v
	}

	result, err := tx.ExecContext(ctx, `
		INSERT INTO auths (
			user_id,
			source,
		    source_id,
		    access_token,
		    refresh_token,
		    expiry,
		    created_at,
		    updated_at,
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`,
		auth.UserID,
		auth.Source,
		auth.SourceID,
		auth.AccessToken,
		auth.RefreshToken,
		expiry,
		(*NullTime)(&auth.CreatedAt),
		(*NullTime)(&auth.UpdatedAt),
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	// set id from database to auth obj.
	auth.ID = int(id)
	return nil
}

// updateAuth refreshes the auth represented by id with accesToken, refreshToken, expiry.
func updateAuth(ctx context.Context, tx *Tx, id int, accesToken, refreshToken string, expiry *time.Time) (*pa.Auth, error) {
	auth, err := findAuthByID(ctx, tx, id) // current auth.
	if err != nil {
		return nil, err
	}

	// refresh auth obj.
	auth.AccessToken = accesToken
	auth.RefreshToken = refreshToken
	auth.Expiry = expiry
	auth.UpdatedAt = tx.now

	if err := auth.Validate(); err != nil {
		return nil, err
	}

	var expiryStringFmt *string
	if auth.Expiry != nil {
		v := auth.Expiry.Format(time.RFC3339)
		expiryStringFmt = &v
	}

	// update db with new refreshed auth.
	if _, err := tx.ExecContext(ctx, `
		UPDATE auths
		SET access_token 	= ?,
			refresh_token 	= ?,
			expiry			= ?,
			updated_at		= ?,
		WHERE id = ?
	`,
		auth.AccessToken,
		auth.RefreshToken,
		expiryStringFmt,
		(*NullTime)(&auth.UpdatedAt),
		id,
	); err != nil {
		return nil, err
	}

	return auth, nil
}

func deleteAuth(ctx context.Context, tx *Tx, id int) error {
	userID := pa.UserIDFromContext(ctx)

	auth, err := findAuthByID(ctx, tx, id)
	if err != nil {
		return err
	}

	if auth.UserID != userID {
		return pa.Errorf(pa.EUNAUTHORIZED, "cannot delete someone else's auth.")
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM auths WHERE id = ?`, id); err != nil {
		return err
	}
	return nil
}

func attachUserToAuth(ctx context.Context, tx *Tx, auth *pa.Auth) (err error) {
	if auth.User, err = findUserByID(ctx, tx, auth.UserID); err != nil {
		return fmt.Errorf("attachUserToAuth: %w", err)
	}
	return nil
}
