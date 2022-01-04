package sqlite

import (
	"context"

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

	user, err := updateUser(ctx, tx, update)
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

func findUsers(ctx context.Context, tx *Tx, filter pa.UserFilter) ([]*pa.User, int, error) {

}

func attachAuthtoUser(ctx context.Context, tx *Tx, user *pa.User) error {

}

func createUser(ctx context.Context, tx *Tx, user *pa.User) error {

}

func updateUser(ctx context.Context, tx *Tx, update pa.UserUpdate) (*pa.User, error) {

}

func deleteUser(ctx context.Context, tx *Tx, id int) error {

}
