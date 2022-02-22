package sqlite_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	pa "github.com/Lambels/patrickarvatu.com"
	"github.com/Lambels/patrickarvatu.com/sqlite"
)

func TestCreateAuth(t *testing.T) {
	t.Run("Ok Create Call", func(t *testing.T) {
		db := MustOpenTempDB(t)
		defer MustCloseDB(t, db)

		backgroundCtx := context.Background()

		authService := sqlite.NewAuthService(db)

		time := time.Date(2022, time.January, 8, 0, 0, 0, 0, time.UTC)
		auth := &pa.Auth{
			User: &pa.User{ // make sure we dont provide uID indicating that we want to create a user.
				Name:  "Jhon",
				Email: "jhon@doe.com",
			},
			Source:       "cool-source",
			SourceID:     "cooler-source-id",
			RefreshToken: "my-secret-token",
			AccessToken:  "my-very-secret-token",
			Expiry:       &time,
		}

		// create auth.
		if err := authService.CreateAuth(backgroundCtx, auth); err != nil {
			t.Fatal(err)
		} else if auth.ID == 0 {
			t.Fatal("got id = 0")
		} else if auth.UserID == 0 {
			t.Fatal("expected user creation")
		} else if auth.User.CreatedAt.IsZero() {
			t.Fatal("expected user time stamp creation (created AT)")
		} else if auth.User.UpdatedAt.IsZero() {
			t.Fatal("expected user time stamp creation (updated AT)")
		} else if auth.CreatedAt.IsZero() {
			t.Fatal("expected auth time stamp creation (created AT)")
		} else if auth.UpdatedAt.IsZero() {
			t.Fatal("expected auth time stamp creation (updated AT)")
		}

		// assert creation.
		if gotAuth, err := authService.FindAuthByID(backgroundCtx, 1); err != nil {
			t.Fatal(err)
		} else if !reflect.DeepEqual(auth, gotAuth) {
			t.Fatal("DeepEqual: auth != gotAuth")
		}
	})

	t.Run("Bad Update Call (Not Found)", func(t *testing.T) {
		db := MustOpenTempDB(t)
		defer MustCloseDB(t, db)

		backgroundCtx := context.Background()

		authService := sqlite.NewAuthService(db)

		time := time.Date(2022, time.January, 8, 0, 0, 0, 0, time.UTC)
		auth := &pa.Auth{
			UserID:       1, // provide uID indicating that we want to update auth on uID 1 which doesent exist.
			Source:       "cool-source",
			SourceID:     "cooler-source-id",
			RefreshToken: "my-secret-token",
			AccessToken:  "my-very-secret-token",
			Expiry:       &time,
		}

		// update auth (Not found).
		if err := authService.CreateAuth(backgroundCtx, auth); err == nil {
			t.Fatal("expected error")
		} else {
			// TODO: Parse foreign key err to not found.
			t.Log(err) // FOREIGN KEY constraint failed "user not found".
		}
	})

	t.Run("Ok Update Call", func(t *testing.T) {
		db := MustOpenTempDB(t)
		defer MustCloseDB(t, db)

		backgroundCtx := context.Background()

		authService := sqlite.NewAuthService(db)

		time := time.Date(2022, time.January, 8, 0, 0, 0, 0, time.UTC)
		auth := &pa.Auth{
			User: &pa.User{ // make sure we dont provide uID indicating that we want to create a user.
				ID:    1,
				Name:  "Jhon",
				Email: "jhon@doe.com",
			},
			Source:       "cool-source",
			SourceID:     "cooler-source-id",
			RefreshToken: "my-secret-token",
			AccessToken:  "my-very-secret-token",
			Expiry:       &time,
		}

		auth2 := &pa.Auth{
			UserID:       1,
			Source:       "cool-source",
			SourceID:     "cooler-source-id",
			RefreshToken: "my-secret-token-2",
			AccessToken:  "my-very-secret-token-2",
			Expiry:       &time,
		}

		// create + update auth.
		MustCreateAuth(t, db, backgroundCtx, auth)
		MustCreateAuth(t, db, backgroundCtx, auth2)

		// assert updating.
		if gotAuth, err := authService.FindAuthByID(backgroundCtx, 1); err != nil {
			t.Fatal(err)
		} else if auth2.AccessToken != gotAuth.AccessToken {
			t.Fatal("Access tokens dont match")
		} else if auth2.RefreshToken != gotAuth.RefreshToken {
			t.Fatal("Refresh tokens dont match")
		}
	})
}

func TestDeleteAuth(t *testing.T) {
	t.Run("Ok Delete Call", func(t *testing.T) {
		db := MustOpenTempDB(t)
		defer MustCloseDB(t, db)

		backgroundCtx := context.Background()

		authService := sqlite.NewAuthService(db)

		auth := &pa.Auth{
			User: &pa.User{
				Name:  "Mona Lisa",
				Email: "octo@cat.com",
			},
			Source:       "cool-source",
			SourceID:     "cooler-source-id",
			RefreshToken: "my-secret-token",
			AccessToken:  "my-very-secret-token",
		}

		// create auth.
		userCtx := MustCreateAuth(t, db, backgroundCtx, auth)

		// delete auth.
		if err := authService.DeleteAuth(userCtx, 1); err != nil {
			t.Fatal(err)
		}

		// assert deletion.
		if _, err := authService.FindAuthByID(backgroundCtx, 1); pa.ErrorCode(err) != pa.ENOTFOUND {
			t.Fatal("err != ENOTFOUND")
		}
	})

	t.Run("Bad Delete Call (Un Authorized)", func(t *testing.T) {
		db := MustOpenTempDB(t)
		defer MustCloseDB(t, db)

		backgroundCtx := context.Background()

		authService := sqlite.NewAuthService(db)

		auth := &pa.Auth{
			User: &pa.User{
				Name:  "Mona Lisa",
				Email: "octo@cat.com",
			},
			Source:       "cool-source",
			SourceID:     "cooler-source-id",
			RefreshToken: "my-secret-token",
			AccessToken:  "my-very-secret-token",
		}

		// create auth.
		MustCreateAuth(t, db, backgroundCtx, auth)

		auth2 := &pa.Auth{
			User: &pa.User{
				Name:  "Bad Man",
				Email: "hacker@hacker.com",
			},
			Source:       "cool-source-2",
			SourceID:     "cooler-source-id-2",
			RefreshToken: "my-secret-token-2",
			AccessToken:  "my-very-secret-token-2",
		}

		// create auth.
		userCtx := MustCreateAuth(t, db, backgroundCtx, auth2)

		// delete auth (Un Auth).
		if err := authService.DeleteAuth(userCtx, 1); pa.ErrorCode(err) != pa.EUNAUTHORIZED {
			t.Fatal("err != EUNAUTHORIZED")
		}
	})

	t.Run("Bad Delete Call (Not Found)", func(t *testing.T) {
		db := MustOpenTempDB(t)
		defer MustCloseDB(t, db)

		backgroundCtx := context.Background()

		authService := sqlite.NewAuthService(db)

		// delete auth. (Not Found).
		if err := authService.DeleteAuth(backgroundCtx, 1); pa.ErrorCode(err) != pa.ENOTFOUND {
			t.Fatal("err != ENOTFOUND")
		}
	})
}

func TestFindAuths(t *testing.T) {
	t.Run("Ok Find Call", func(t *testing.T) {
		db := MustOpenTempDB(t)
		defer MustCloseDB(t, db)

		backgroundCtx := context.Background()

		authService := sqlite.NewAuthService(db)

		user1 := &pa.User{
			Name:  "Lambels",
			Email: "lamb@lambels.com",
		}

		auth1 := &pa.Auth{
			Source:      "good_source",
			SourceID:    "good_source_id_1",
			AccessToken: "some_access_token_1",
			User:        user1,
		}

		MustCreateAuth(t, db, backgroundCtx, auth1)

		auth2 := &pa.Auth{
			Source:      "ok_source",
			SourceID:    "ok_source_id_1",
			AccessToken: "some_access_token_2",
			User:        user1,
		}

		MustCreateAuth(t, db, backgroundCtx, auth2)

		user2 := &pa.User{
			Name:  "Patrick",
			Email: "patrick@arvatu.com",
		}

		auth3 := &pa.Auth{
			Source:      "meh_source",
			SourceID:    "meh_source_id_1",
			AccessToken: "some_access_token_3",
			User:        user2,
		}

		MustCreateAuth(t, db, backgroundCtx, auth3)

		if gotAuths, n, err := authService.FindAuths(backgroundCtx, pa.AuthFilter{UserID: NewIntPointer(1)}); err != nil {
			t.Fatal(err)
		} else if len(gotAuths) != 2 {
			t.Fatalf("len=%v != 2", len(gotAuths))
		} else if n != 2 {
			t.Fatalf("n=%v != 2", n)
		} else if gotAuths[0].SourceID != auth1.SourceID {
			t.Fatalf("gotAuth SourceID=%v != auth1 SourceID=%v", gotAuths[0].SourceID, auth1.SourceID)
		} else if gotAuths[1].SourceID != auth2.SourceID {
			t.Fatalf("gotAuth SourceID=%v != auth1 SourceID=%v", gotAuths[1].SourceID, auth2.SourceID)
		}
	})

	t.Run("Bad Find Call (Not Found)", func(t *testing.T) {
		db := MustOpenTempDB(t)
		defer MustCloseDB(t, db)

		backgroundCtx := context.Background()

		authService := sqlite.NewAuthService(db)

		if _, err := authService.FindAuthByID(backgroundCtx, 1); pa.ErrorCode(err) != pa.ENOTFOUND {
			t.Fatal("err != ENOTFOUND")
		}
	})
}

func MustCreateAuth(t *testing.T, db *sqlite.DB, ctx context.Context, auth *pa.Auth) context.Context {
	t.Helper()
	if err := sqlite.NewAuthService(db).CreateAuth(ctx, auth); err != nil {
		t.Fatal(err)
	}
	return pa.NewContextWithUser(ctx, auth.User)
}
