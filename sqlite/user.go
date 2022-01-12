package sqlite

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"strings"

	pa "github.com/Lambels/patrickarvatu.com"
)

// check to see if *UserService object implements set interface.
var _ pa.UserService = (*UserService)(nil)

// UserService represents a service used to manage users.
type UserService struct {
	db *DB
}

// NewUserService returns a new instance of UserService attached to db.
func NewUserService(db *DB) *UserService {
	return &UserService{
		db: db,
	}
}

// FindUserByID returns a user based on id.
// returns ENOTFOUND if the user doesent exist.
func (s *UserService) FindUserByID(ctx context.Context, id int) (*pa.User, error) {
	tx, err := s.db.BeginTX(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	user, err := findUserByID(ctx, tx, id)
	if err != nil {
		return user, err

	} else if err := attachAuthtoUser(ctx, tx, user); err != nil {
		return nil, err
	}
	return user, nil
}

// FindUsers returns a range of user based on filter.
func (s *UserService) FindUsers(ctx context.Context, filter pa.UserFilter) ([]*pa.User, int, error) {
	tx, err := s.db.BeginTX(ctx, nil)
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback()

	users, n, err := findUsers(ctx, tx, filter)
	if err != nil {
		return users, n, err
	}

	for _, user := range users {
		if err := attachAuthtoUser(ctx, tx, user); err != nil {
			return users, n, err
		}
	}
	return users, n, nil
}

// CreateUser creates a new user. To only be used in testing as users are created through the
// create auth process AuthService.CreateAuth() -> ./auth.go
func (s *UserService) CreateUser(ctx context.Context, user *pa.User) error {
	tx, err := s.db.BeginTX(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := createUser(ctx, tx, user); err != nil {
		return err

	} else if err = attachAuthtoUser(ctx, tx, user); err != nil {
		return err
	}
	return tx.Commit()
}

// UpdateUser updates user with id: id.
// returns EUNAUTHORIZED if the user isnt trying to update himself.
// returns ENOTFOUND if the user doesent exist.
func (s *UserService) UpdateUser(ctx context.Context, id int, update pa.UserUpdate) (*pa.User, error) {
	tx, err := s.db.BeginTX(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	user, err := updateUser(ctx, tx, id, update)
	if err != nil {
		return nil, err

	} else if err := attachAuthtoUser(ctx, tx, user); err != nil {
		return nil, err

	} else if err := tx.Commit(); err != nil {
		return nil, err
	}
	return user, nil
}

// DeleteUser permanently deletes the user specified by id.
// returns EUNAUTHORIZED if the user isnt trying to delete himself.
// returns ENOTFOUND if the user doesent exist.
func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	tx, err := s.db.BeginTX(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := deleteUser(ctx, tx, id); err != nil {
		return err
	}
	return tx.Commit()
}

func findUserByID(ctx context.Context, tx *Tx, id int) (*pa.User, error) {
	filter := pa.UserFilter{
		ID: &id,
	}
	users, _, err := findUsers(ctx, tx, filter)
	if err != nil {
		return nil, err
	} else if len(users) == 0 {
		return nil, pa.Errorf(pa.ENOTFOUND, "user not found")
	}
	return users[0], nil
}

func findUserByEmail(ctx context.Context, tx *Tx, email string) (*pa.User, error) {
	filter := pa.UserFilter{
		Email: &email,
	}
	users, _, err := findUsers(ctx, tx, filter)
	if err != nil {
		return nil, err
	} else if len(users) == 0 {
		return nil, pa.Errorf(pa.ENOTFOUND, "user not found")
	}
	return users[0], nil
}

func findUsers(ctx context.Context, tx *Tx, filter pa.UserFilter) (_ []*pa.User, n int, err error) {
	// build where and args statement method.
	// not vulnerable to sql injection attack.
	where, args := []string{"1 = 1"}, []interface{}{}
	if v := filter.ID; v != nil {
		where, args = append(where, "id = ?"), append(args, *v)
	}
	if v := filter.Email; v != nil {
		where, args = append(where, "email = ?"), append(args, *v)
	}
	if v := filter.APIKey; v != nil {
		where, args = append(where, "api_key = ?"), append(args, *v)
	}

	rows, err := tx.QueryContext(ctx, `
		SELECT 
		    id,
		    name,
		    email,
		    api_key,
		    created_at,
		    updated_at,
		    COUNT(*) OVER()
		FROM users
		WHERE `+strings.Join(where, " AND ")+`
		ORDER BY id ASC
		`+FormatLimitOffset(filter.Limit, filter.Offset),
		args...,
	)

	if err != nil {
		return nil, n, err
	}
	defer rows.Close()

	// deserialize rows.
	users := []*pa.User{}
	for rows.Next() {
		var email sql.NullString
		var user *pa.User
		if err := rows.Scan(
			&user.ID,
			&user.Name,
			&email,
			&user.APIKey,
			(*NullTime)(&user.CreatedAt),
			(*NullTime)(&user.UpdatedAt),
			&n,
		); err != nil {
			return nil, 0, err
		}

		if email.Valid {
			user.Email = email.String
		}

		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return users, n, nil
}

func createUser(ctx context.Context, tx *Tx, user *pa.User) error {
	user.CreatedAt = tx.now
	user.UpdatedAt = user.CreatedAt

	if err := user.Validate(); err != nil {
		return err
	}

	// make sure to instantiate email to pass it as NULL if not provided in the struct by the OAuth providers.
	var email *string
	if user.Email != "" {
		email = &user.Email
	}

	// generate a random api-key.
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return err
	}
	// encode rand bytes.
	user.APIKey = base64.StdEncoding.EncodeToString(buf)

	result, err := tx.ExecContext(ctx, `
		INSERT INTO users (
			name,
			email,
			api_key,
			created_at,
			updated_at
		)
		VALUES (?, ?, ?, ?, ?)
	`,
		user.Name,
		email,
		user.APIKey,
		(*NullTime)(&user.CreatedAt),
		(*NullTime)(&user.UpdatedAt),
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	// set id from database to user obj.
	user.ID = int(id)

	return nil
}

func updateUser(ctx context.Context, tx *Tx, id int, update pa.UserUpdate) (*pa.User, error) {
	if pa.UserIDFromContext(ctx) != id {
		return nil, pa.Errorf(pa.EUNAUTHORIZED, "user not authorized")
	}

	user, err := findUserByID(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	if v := update.Name; v != nil {
		user.Name = *v
	}
	if v := update.Email; v != nil {
		user.Email = *v
	}

	user.UpdatedAt = tx.now

	if err := user.Validate(); err != nil {
		return nil, err
	}

	// make sure to instantiate email to pass it as NULL if not provided in the struct by the OAuth providers.
	var email *string
	if user.Email != "" {
		email = &user.Email
	}

	if _, err := tx.ExecContext(ctx, `
		UPDATE users
		SET name = ?,
		    email = ?,
		    updated_at = ?
		WHERE id = ?
	`,
		user.Name,
		email,
		(*NullTime)(&user.UpdatedAt),
		id,
	); err != nil {
		return user, err
	}

	return user, nil
}

func deleteUser(ctx context.Context, tx *Tx, id int) error {
	user, err := findUserByID(ctx, tx, id)
	if err != nil {
		return err
	}

	if pa.UserIDFromContext(ctx) != user.ID {
		return pa.Errorf(pa.EUNAUTHORIZED, "user not authorized")
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM users WHERE id = ?`, id); err != nil {
		return err
	}
	return nil
}

func attachAuthtoUser(ctx context.Context, tx *Tx, user *pa.User) (err error) {
	filter := pa.AuthFilter{
		UserID: &user.ID,
	}

	if user.Auths, _, err = findAuths(ctx, tx, filter); err != nil {
		return err
	}
	return nil
}
