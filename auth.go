package pa

import (
	"context"
	"fmt"
	"time"
)

// auth sources represent different OAuth providers, the system is currently supporting only
// github as a provider but its implemented with this issue taken in mind.
const (
	AuthSourceGitHub = "github"
)

// Auth represents an OAuth object in the system.
type Auth struct {
	// the pk of the auth.
	ID int `json:"id"`

	// fields linking the auth object back to the user.
	UserID int   `json:"userID"`
	User   *User `json:"user"`

	// the source from where the OAuth object comes from, ie: "github".
	Source   string `json:"source"`
	SourceID string `json:"sourceID"`

	// OAuth credentials provided by the OAuth source.
	AccessToken  string     `json:"-"`
	RefreshToken string     `json:"-"`
	Expiry       *time.Time `json:"-"`

	// timestamps.
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Validate performs basic validation on Auth.
// returns EINVALID if any error is found.
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

// AvatarURL returns a URL to the avatar image provided by the OAuth source.
// returns an emtpy string if no source is identified.
func (a *Auth) AvatarURL(size int) string {
	switch a.Source {
	case AuthSourceGitHub:
		return fmt.Sprintf("https://avatars1.githubusercontent.com/u/%s?s=%d", a.SourceID, size)
	default:
		return ""
	}
}

// AuthService represents a service which manages auth in the system.
type AuthService interface {
	// FindAuthByID returns a auth based on the id.
	// returns ENOTFOUND if the auth doesent exist.
	FindAuthByID(ctx context.Context, id int) (*Auth, error)

	// FindAuths returns a range of auths and the length of the range. If filter
	// is specified FindAuths will apply the filter to return set response.
	FindAuths(ctx context.Context, filter AuthFilter) ([]*Auth, int, error)

	// CreateAuth creates a auth. Main entry point when creating a user.
	// The creation will only go through if the linking fields are attached or the object passes the validation.
	// On creation the auth will be linked to a user if found, otherwise the user gets created and the auth gets
	// linked to the user and the user linked to the auth through the linking fields.
	CreateAuth(ctx context.Context, auth *Auth) error

	// DeleteAuth permanently deletes a subscription. The linked user wont be deleted bu will appear as
	// not validated.
	DeleteAuth(ctx context.Context, id int) error
}

// AuthFilter represents a filter used by FindAuths to filter the response.
type AuthFilter struct {
	// fields to filter on.
	ID       *int    `json:"id"`
	UserID   *int    `json:"userID"`
	Source   *string `json:"source"`
	SourceID *string `json:"sourceID"`

	// restrictions on the result set, used for pagination and set limits.
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}
