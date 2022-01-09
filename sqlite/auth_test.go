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
	t.Run("Ok Call", func(t *testing.T) {
		db := MustOpenDB(t)
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

		// create auth.
		if err := authService.CreateAuth(backgroundCtx, auth); err != nil {
			t.Fatal(err)
		} else if auth.ID == 0 {
			t.Fatal("got id = 0")
		}

		// assert creation.
		if gotAuth, err := authService.FindAuthByID(backgroundCtx, 1); err != nil {
			t.Fatal(err)
		} else if !reflect.DeepEqual(auth, gotAuth) {
			t.Fatal("DeepEqual: auth != gotAuth")
		}
	})

	t.Run("User Not Found", func(t *testing.T) {
		db := MustOpenDB(t)
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

		if err := authService.CreateAuth(backgroundCtx, auth); pa.ErrorCode(err) != pa.ENOTFOUND {
			t.Fatal("err != pa.ENOTFOUND")
		}
	})

	// TODO: add ok update auth test.
}
