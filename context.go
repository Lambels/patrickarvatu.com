package pa

import "context"

type contextKey int

const (
	// userContextKey holds the user inside a ctx.
	userContextKey = contextKey(iota + 1)
)

// NewContextWithUser enriches the context ctx with the user: user under the key userContextKey.
func NewContextWithUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

// UserFromContext pulls the user from context ctx.
func UserFromContext(ctx context.Context) *User {
	usr, _ := ctx.Value(userContextKey).(*User)
	return usr
}

// UserIDFromContext is a helper function which returns only the id of the user under ctx.
// To only be used when checking id with id.
func UserIDFromContext(ctx context.Context) int {
	if usr := UserFromContext(ctx); usr != nil {
		return usr.ID
	}
	return 0
}

// IsAdminContext is a helper function to check if context: ctx is an admin context.
func IsAdminContext(ctx context.Context) bool {
	if usr := UserFromContext(ctx); usr != nil {
		return usr.IsAdmin
	}
	return false
}
