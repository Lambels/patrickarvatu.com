package pa

import "context"

type contextKey int

const (
	userContextKey = contextKey(iota + 1)
)

func NewContextWithUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

func UserFromContext(ctx context.Context) *User {
	return ctx.Value(userContextKey).(*User)
}

func UserIDFromContext(ctx context.Context) int {
	if usr := UserFromContext(ctx); usr != nil {
		return usr.ID
	}
	return 0
}
