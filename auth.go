package pa

import (
	"context"
	"fmt"
	"time"
)

const (
	AuthSourceGitHub = "github"
)

type Auth struct {
	ID int `json:"id"`

	UserID int   `json:"userID"`
	User   *User `json:"user"`

	Source   string `json:"source"`
	SourceID string `json:"sourceID"`

	AccessToken  string     `json:"-"`
	RefreshToken string     `json:"-"`
	Expiry       *time.Time `json:"-"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (a *Auth) Validate() error {
	if a.UserID == 0 {
		return Errorf(EINVALID, "User required.")
	} else if a.Source == "" {
		return Errorf(EINVALID, "Source required.")
	} else if a.SourceID == "" {
		return Errorf(EINVALID, "Source ID required.")
	} else if a.AccessToken == "" {
		return Errorf(EINVALID, "Access token required.")
	}
	return nil
}

func (a *Auth) AvatarURL(size int) string {
	switch a.Source {
	case AuthSourceGitHub:
		return fmt.Sprintf("https://avatars1.githubusercontent.com/u/%s?s=%d", a.SourceID, size)
	default:
		return ""
	}
}

type AuthService interface {
	FindAuthByID(ctx context.Context, id int) (*Auth, error)

	FindAuths(ctx context.Context, filter AuthFilter) ([]*Auth, int, error)

	CreateAuth(ctx context.Context, auth *Auth) error

	DeleteAuth(ctx context.Context, id int) error
}

type AuthFilter struct {
	ID       *int    `json:"id"`
	UserID   *int    `json:"userID"`
	Source   *string `json:"source"`
	SourceID *string `json:"sourceID"`

	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}
