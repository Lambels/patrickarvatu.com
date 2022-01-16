package sqlite_test

import (
	"context"
	"reflect"
	"testing"

	pa "github.com/Lambels/patrickarvatu.com"
	"github.com/Lambels/patrickarvatu.com/sqlite"
)

// TODO: add same comments and variable names.

func TestCreateComment(t *testing.T) {
	t.Parallel()
	t.Run("Ok Create Call", func(t *testing.T) {
		db := MustOpenTempDB(t)
		defer MustCloseDB(t, db)

		backgroundCtx := context.Background()

		commentService := sqlite.NewCommentService(db)
		userService := sqlite.NewUserService(db)
		blogService := sqlite.NewBlogService(db)
		subBlogService := sqlite.NewSubBlogService(db)

		user := &pa.User{
			Name:    "Lambels",
			Email:   "lamb@lambels.com",
			IsAdmin: true,
		}

		// create user.
		if err := userService.CreateUser(backgroundCtx, user); err != nil {
			t.Fatal(err)
		}

		// declare ctx with enriched user.
		adminUserCtx := pa.NewContextWithUser(backgroundCtx, user)

		blog := &pa.Blog{
			Title:       "Cool Title",
			Description: "Idk man",
		}

		// create blog.
		if err := blogService.CreateBlog(adminUserCtx, blog); err != nil {
			t.Fatal(err)
		}

		subBlog := &pa.SubBlog{
			BlogID:  blog.ID,
			Title:   "Cool Sub blog",
			Content: "idk",
		}

		// create sub blog.
		if err := subBlogService.CreateSubBlog(adminUserCtx, subBlog); err != nil {
			t.Fatal(err)
		}

		comment := &pa.Comment{
			SubBlogID: subBlog.ID,
			UserID:    user.ID,
			Content:   "Cool content",
		}

		// create comment.
		if err := commentService.CreateComment(adminUserCtx, comment); err != nil {
			t.Fatal(err)
		} else if comment.ID == 0 {
			t.Fatal("got id = 0")
		} else if comment.User == nil {
			t.Fatalf("expected user attachment: %v", *comment.User)
		} else if comment.CreatedAt.IsZero() {
			t.Fatal("expected comment time stamp creation (created AT)")
		}

		// assert creation.
		if gotComment, err := commentService.FindCommentByID(backgroundCtx, 1); err != nil {
			t.Fatal(err)
		} else if !reflect.DeepEqual(gotComment, comment) {
			t.Fatal("DeepEqual: gotComment != comment")
		}
	})
}

func TestCommentAttachments(t *testing.T) {

}
