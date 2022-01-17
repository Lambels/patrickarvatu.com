package sqlite_test

import (
	"context"
	"testing"

	pa "github.com/Lambels/patrickarvatu.com"
	"github.com/Lambels/patrickarvatu.com/sqlite"
)

// TODO: write test

func MustCreateUser(t *testing.T, db *sqlite.DB, ctx context.Context, user *pa.User) context.Context {
	t.Helper()
	if err := sqlite.NewUserService(db).CreateUser(ctx, user); err != nil {
		t.Fatal(err)
	}
	return pa.NewContextWithUser(ctx, user)
}
