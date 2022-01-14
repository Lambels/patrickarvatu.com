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
	t.Parallel() // run tests in parallel.
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

	t.Run("Bad Update Call", func(t *testing.T) {
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

		if err := authService.CreateAuth(backgroundCtx, auth); err == nil {
			t.Fatal("expected error")
		} else {
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
		if err := authService.CreateAuth(backgroundCtx, auth); err != nil {
			t.Fatal(err)
		} else if err := authService.CreateAuth(backgroundCtx, auth2); err != nil {
			t.Fatal(err)
		}

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

// TODO: add testing for all CRUD functions.
