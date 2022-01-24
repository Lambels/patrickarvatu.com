package pa

import (
	"context"
	"time"
)

// User represents an user in the system.
type User struct {
	// the pk of the user.
	ID int `json:"id"`

	// name / email
	Name  string `json:"name"`
	Email string `json:"email"`

	// apikey for the user to communicate to the api on behalf of the user.
	APIKey string `json:"-"`

	// timestamps.
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	// assosciated auths.
	Auths []*Auth `json:"auths"`

	// field to identify the user as an admin.
	IsAdmin bool `json:"isAdmin"`
}

// Vlidate performs basic validation on User.
// returns EINVALID if any error is found.
func (u *User) Validate() error {
	if u.Name == "" {
		return Errorf(EINVALID, "name is a required field.")
	}
	return nil
}

// AvatarURL checks the first auth in .Auths and returns a URL to the users pfp on set auth source.
// returns an empty string if no avatar URL is found.
func (u *User) AvatarURL(size int) string {
	for _, auth := range u.Auths {
		if s := auth.AvatarURL(size); s != "" {
			return s
		}
	}

	return ""
}

// UserService represents a service which manages users in the system.
type UserService interface {
	// FindUserByID returns a user based on id.
	// returns ENOTFOUND if the user doesent exist.
	FindUserByID(ctx context.Context, id int) (*User, error)

	// FindUsers returns a range of users and the length of the range. If filter
	// is specified FindUsers will apply the filter to return set response.
	FindUsers(ctx context.Context, filter UserFilter) ([]*User, int, error)

	// CreateUser creates an user. To only be used in testing, the main pipeline when creating an user
	// starts over at CreateAuth -> ./auth.go
	CreateUser(ctx context.Context, user *User) error

	// UpdateUser updates a user based on the update field.
	// returns ENOTFOUND if the user doesent exist.
	// returns EUHATHORIZED if the caller isnt trying to update himself.
	UpdateUser(ctx context.Context, id int, update UserUpdate) (*User, error)

	// DeleteUser permanently deletes a user. This also permanently deletes all of the users assosiacions
	// such as: auths, comments.
	// returns ENOTFOUND if user doesent exist.
	// returns EUHATHORIZED if the caller isnt trying to delete himself.
	DeleteUser(ctx context.Context, id int) error
}

// UserFilter represents a filter used by FindUsers to filter the response.
type UserFilter struct {
	// fields to filter on.
	ID     *int    `json:"id"`
	Email  *string `json:"email"`
	APIKey *string `json:"apiKey"`

	// restrictions on the result set, used for pagination and set limits.
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

// UserUpdate represents an update used by UpdateUser to update a user.
type UserUpdate struct {
	// fields which can be updated.
	Name   *string `json:"name"`
	Email  *string `json:"email"`
	ApiKey *string `json:"apiKey"` // TODO: test api key update.
}
