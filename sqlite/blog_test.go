package sqlite_test

import (
	"context"
	"reflect"
	"testing"

	pa "github.com/Lambels/patrickarvatu.com"
	"github.com/Lambels/patrickarvatu.com/sqlite"
)

func TestCreateBlog(t *testing.T) {
	t.Run("Ok Call", func(t *testing.T) {
		db := MustOpenDB(t)
		defer MustCloseDB(t, db)

		backgroundCtx := context.Background()
		adminUsrContext := pa.NewContextWithUser(backgroundCtx, &pa.User{
			Name:    "jhon DOE",
			Email:   "jhon@doe.com",
			IsAdmin: true,
		})

		blogService := sqlite.NewBlogService(db)

		blog := &pa.Blog{
			Title:       "Epic Blog",
			Description: "Honestly the best blog ever.",
		}

		// create blog.
		if err := blogService.CreateBlog(adminUsrContext, blog); err != nil {
			t.Fatal(err)
		} else if blog.ID == 0 {
			t.Fatal("got id = 0")
		} else if blog.CreatedAt.IsZero() {
			t.Fatal("expected blog time stamp creation (created AT)")
		} else if blog.UpdatedAt.IsZero() {
			t.Fatal("expected blog time stamp creation (updated AT)")
		}

		// assert creation.
		if gotBlog, err := blogService.FindBlogByID(backgroundCtx, 1); err != nil {
			t.Fatal(err)
		} else if !reflect.DeepEqual(gotBlog, blog) {
			t.Fatal("DeepEqual: gotBlog != blog")
		}
	})

	t.Run("UnAuth", func(t *testing.T) {
		db := MustOpenDB(t)
		defer MustCloseDB(t, db)

		backgroundCtx := context.Background()
		usrContext := pa.NewContextWithUser(backgroundCtx, &pa.User{
			Name:  "jhon DOE",
			Email: "jhon@doe.com",
		})

		blogService := sqlite.NewBlogService(db)

		blog := &pa.Blog{
			Title:       "Epic Blog",
			Description: "Honestly the best blog ever.",
		}

		// create blog.
		if err := blogService.CreateBlog(usrContext, blog); pa.ErrorCode(err) != pa.EUNAUTHORIZED {
			t.Fatal("expected UnAuth error")
		} else if blog.ID != 0 {
			t.Fatal("got id != 0")
		}
	})
}

// TODO: add testing for rest of CRUD methods.
