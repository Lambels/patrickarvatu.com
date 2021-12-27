package pa

import (
	"context"
	"time"
)

type User struct {
	ID int `json:"id"`

	Name  string `json:"name"`
	Email string `json:"email"`

	APIKey string `json:"-"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Auths []*Auth `json:"auths"`
}

func (u *User) Validate() error {
	if u.Name == "" {
		return Errorf(EINVALID, "name is a required field.")
	}
	return nil
}

func (u *User) AvatarURL(size int) string {
	for _, auth := range u.Auths {
		if s := auth.AvatarURL(size); s != nil {
			return s
		}
	}

	return ""
}

type UserService interface {
	FindUserByID(ctx context.Context, id int) (*User, error)

	FindUsers(ctx context.Context, filter UserFilter) ([]*User, error)

	CreateUser(ctx context.Context, user *User) error

	UpdateUser(ctx context.Context, update UserUpdate) (*User, error)

	DeleteUser(ctx context.Context, id int) error
}

type UserFilter struct {
	ID     *int    `json:"id"`
	Email  *string `json:"email"`
	APIKey *string `json:"apiKey"`

	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type UserUpdate struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
}
