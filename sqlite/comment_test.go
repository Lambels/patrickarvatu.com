package sqlite_test

import (
	"context"
	"reflect"
	"testing"

	pa "github.com/Lambels/patrickarvatu.com"
	"github.com/Lambels/patrickarvatu.com/sqlite"
)

func TestCreateComment(t *testing.T) {
	t.Run("Ok Create Call", func(t *testing.T) {
		db := MustOpenTempDB(t)
		defer MustCloseDB(t, db)

		backgroundCtx := context.Background()

		commentService := sqlite.NewCommentService(db)

		user := &pa.User{
			Name:    "Lambels",
			Email:   "lamb@lambels.com",
			IsAdmin: true,
		}

		// create user.
		adminUsrCtx := MustCreateUser(t, db, backgroundCtx, user)

		blog := &pa.Blog{
			Title:       "Cool Title",
			Description: "Idk man",
		}

		// create blog.
		MustCreateBlog(t, db, adminUsrCtx, blog)

		subBlog := &pa.SubBlog{
			BlogID:  blog.ID,
			Title:   "Cool Sub blog",
			Content: "idk",
		}

		// create sub blog.
		MustCreateSubBlog(t, db, adminUsrCtx, subBlog)

		comment := &pa.Comment{
			SubBlogID: subBlog.ID,
			Content:   "Cool content",
		}

		// create comment.
		if err := commentService.CreateComment(adminUsrCtx, comment); err != nil {
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

func TestDeleteComment(t *testing.T) {
	t.Run("Ok Delete Call", func(t *testing.T) {
		db := MustOpenTempDB(t)
		defer MustCloseDB(t, db)

		backgroundCtx := context.Background()

		commentService := sqlite.NewCommentService(db)

		user := &pa.User{
			Name:    "Lambels",
			Email:   "lamb@lambels.com",
			IsAdmin: true,
		}

		// create user.
		adminUsrCtx := MustCreateUser(t, db, backgroundCtx, user)

		blog := &pa.Blog{
			Title:       "Cool Title",
			Description: "Idk man",
		}

		// create blog.
		MustCreateBlog(t, db, adminUsrCtx, blog)

		subBlog := &pa.SubBlog{
			BlogID:  blog.ID,
			Title:   "Cool Sub blog",
			Content: "idk",
		}

		// create sub blog.
		MustCreateSubBlog(t, db, adminUsrCtx, subBlog)

		comment := &pa.Comment{
			SubBlogID: subBlog.ID,
			Content:   "Cool content",
		}

		// create comment.
		MustCreateComment(t, db, adminUsrCtx, comment)

		// delete comment.
		if err := commentService.DeleteComment(adminUsrCtx, 1); err != nil {
			t.Fatal(err)
		}

		// assert deletion
		if err := commentService.DeleteComment(adminUsrCtx, 1); pa.ErrorCode(err) != pa.ENOTFOUND {
			t.Fatal("err != ENOTFOUND")
		}
	})

	t.Run("Bad Delete Call (Un Auth)", func(t *testing.T) {
		db := MustOpenTempDB(t)
		defer MustCloseDB(t, db)

		backgroundCtx := context.Background()

		commentService := sqlite.NewCommentService(db)

		user := &pa.User{
			Name:    "Lambels",
			Email:   "lamb@lambels.com",
			IsAdmin: true,
		}

		// create user.
		adminUsrCtx := MustCreateUser(t, db, backgroundCtx, user)

		blog := &pa.Blog{
			Title:       "Cool Title",
			Description: "Idk man",
		}

		// create blog.
		MustCreateBlog(t, db, adminUsrCtx, blog)

		subBlog := &pa.SubBlog{
			BlogID:  blog.ID,
			Title:   "Cool Sub blog",
			Content: "idk",
		}

		// create sub blog.
		MustCreateSubBlog(t, db, adminUsrCtx, subBlog)

		comment := &pa.Comment{
			SubBlogID: subBlog.ID,
			Content:   "Cool content",
		}

		// create comment.
		MustCreateComment(t, db, adminUsrCtx, comment)

		user2 := &pa.User{
			Name:  "Hakcer",
			Email: "dfsf@sdff.com",
		}

		// create user.
		usr2Ctx := MustCreateUser(t, db, backgroundCtx, user2)

		// delete comment (Un Auth).
		if err := commentService.DeleteComment(usr2Ctx, comment.ID); pa.ErrorCode(err) != pa.EUNAUTHORIZED {
			t.Fatal("err != EUNAUTHORIZED")
		}
	})

	t.Run("Bad Delete Call (Not Found)", func(t *testing.T) {
		db := MustOpenTempDB(t)
		defer MustCloseDB(t, db)

		backgroundCtx := context.Background()

		commentService := sqlite.NewCommentService(db)

		user := &pa.User{
			Name:  "Lambels",
			Email: "lambi@lambels.com",
		}

		// create user.
		usrCtx := MustCreateUser(t, db, backgroundCtx, user)

		// delete comment (Not Found).
		if err := commentService.DeleteComment(usrCtx, 1); pa.ErrorCode(err) != pa.ENOTFOUND {
			t.Fatal("err != ENOTFOUND")
		}
	})
}

func TestUpdateComment(t *testing.T) {
	t.Run("Ok Update Call", func(t *testing.T) {
		db := MustOpenTempDB(t)
		defer MustCloseDB(t, db)

		backgroundCtx := context.Background()

		commentService := sqlite.NewCommentService(db)

		user := &pa.User{
			Name:    "Patrick",
			Email:   "patrick@arvatu.com",
			IsAdmin: true,
		}

		// create user.
		adminUsrCtx := MustCreateUser(t, db, backgroundCtx, user)

		blog := &pa.Blog{
			Title:       "Cool Title",
			Description: "Idk man",
		}

		// create blog.
		MustCreateBlog(t, db, adminUsrCtx, blog)

		subBlog := &pa.SubBlog{
			BlogID:  blog.ID,
			Title:   "Cool Sub blog",
			Content: "idk",
		}

		// create sub blog.
		MustCreateSubBlog(t, db, adminUsrCtx, subBlog)

		user2 := &pa.User{
			Name:  "Lambels",
			Email: "lamb@lambels.com",
		}

		// create user.
		usr2Ctx := MustCreateUser(t, db, backgroundCtx, user2)

		comment := &pa.Comment{
			SubBlogID: subBlog.ID,
			Content:   "Cool content",
		}

		// create comment.
		MustCreateComment(t, db, usr2Ctx, comment)

		update := pa.CommentUpdate{
			Content: NewStringPointer("Cool new content"),
		}

		// update comment.
		if updatedComment, err := commentService.UpdateComment(adminUsrCtx, comment.ID, update); err != nil {
			t.Fatal(err)
		} else if gotComment, err := commentService.FindCommentByID(backgroundCtx, comment.ID); err != nil { // assert update.
			t.Fatal(err)
		} else if !reflect.DeepEqual(updatedComment, gotComment) {
			t.Fatal("DeepEqual: updatedComment != gotComment")
		}
	})

	t.Run("Bad Update Call (Un Auth)", func(t *testing.T) {
		db := MustOpenTempDB(t)
		defer MustCloseDB(t, db)

		backgroundCtx := context.Background()

		commentService := sqlite.NewCommentService(db)

		user := &pa.User{
			Name:    "Patrick",
			Email:   "patrick@arvatu.com",
			IsAdmin: true,
		}

		// create user.
		adminUsrCtx := MustCreateUser(t, db, backgroundCtx, user)

		blog := &pa.Blog{
			Title:       "Cool Title",
			Description: "Idk man",
		}

		// create blog.
		MustCreateBlog(t, db, adminUsrCtx, blog)

		subBlog := &pa.SubBlog{
			BlogID:  blog.ID,
			Title:   "Cool Sub blog",
			Content: "idk",
		}

		// create sub blog.
		MustCreateSubBlog(t, db, adminUsrCtx, subBlog)

		user2 := &pa.User{
			Name:  "Lambels",
			Email: "lamb@lambels.com",
		}

		// create user.
		usr2Ctx := MustCreateUser(t, db, backgroundCtx, user2)

		comment := &pa.Comment{
			SubBlogID: subBlog.ID,
			Content:   "Cool content",
		}

		// create comment.
		MustCreateComment(t, db, usr2Ctx, comment)

		update := pa.CommentUpdate{
			Content: NewStringPointer("Cool new content"),
		}

		// update comment (Un Auth).
		if _, err := commentService.UpdateComment(usr2Ctx, comment.ID, update); pa.ErrorCode(err) != pa.EUNAUTHORIZED {
			t.Fatal("err != EUNAUTHORIZED")
		}
	})

	t.Run("Bad Update Call (Not Found)", func(t *testing.T) {
		db := MustOpenTempDB(t)
		defer MustCloseDB(t, db)

		backgroundCtx := context.Background()

		commentService := sqlite.NewCommentService(db)

		user := &pa.User{
			Name:    "Jhon Doe",
			Email:   "jhon@doe.com",
			IsAdmin: true,
		}

		// create user.
		adminUsrCtx := MustCreateUser(t, db, backgroundCtx, user)

		update := pa.CommentUpdate{
			Content: NewStringPointer("New Content"),
		}

		if _, err := commentService.UpdateComment(adminUsrCtx, 1, update); pa.ErrorCode(err) != pa.ENOTFOUND {
			t.Fatal("err != ENOTFOUND")
		}
	})
}

func TestFindComments(t *testing.T) {
	t.Run("Ok Find Call (filter - sub blog)", func(t *testing.T) {
		db := MustOpenTempDB(t)
		defer MustCloseDB(t, db)

		backgroundCtx := context.Background()

		commentService := sqlite.NewCommentService(db)

		user := &pa.User{
			Name:    "Lambels",
			Email:   "lamb@lambels.com",
			IsAdmin: true,
		}

		// create user.
		adminUsrCtx := MustCreateUser(t, db, backgroundCtx, user)

		blog := &pa.Blog{
			Title:       "Cool Title",
			Description: "Idk man",
		}

		// create blog.
		MustCreateBlog(t, db, adminUsrCtx, blog)

		subBlog := &pa.SubBlog{
			BlogID:  blog.ID,
			Title:   "Cool Sub blog",
			Content: "idk",
		}

		// create sub blog.
		MustCreateSubBlog(t, db, adminUsrCtx, subBlog)

		comment := &pa.Comment{
			SubBlogID: subBlog.ID,
			Content:   "Cool content",
		}

		// create comment.
		MustCreateComment(t, db, adminUsrCtx, comment)

		comment2 := &pa.Comment{
			SubBlogID: subBlog.ID,
			Content:   "BAd Content",
		}

		// create comment.
		MustCreateComment(t, db, adminUsrCtx, comment2)

		// find comments.
		if gotComments, n, err := commentService.FindComments(backgroundCtx, pa.CommentFilter{SubBlogID: NewIntPointer(subBlog.ID)}); err != nil {
			t.Fatal(err)
		} else if len(gotComments) != 2 { // assert find
			t.Fatalf("len=%v != 2", len(gotComments))
		} else if n != 2 {
			t.Fatalf("n=%v != 2", n)
		}
	})

	t.Run("Ok Find Call (filter - user)", func(t *testing.T) {
		db := MustOpenTempDB(t)
		defer MustCloseDB(t, db)

		backgroundCtx := context.Background()

		commentService := sqlite.NewCommentService(db)

		user := &pa.User{
			Name:    "Lambels",
			Email:   "lamb@lambels.com",
			IsAdmin: true,
		}

		// create user.
		adminUsrCtx := MustCreateUser(t, db, backgroundCtx, user)

		blog := &pa.Blog{
			Title:       "Cool Title",
			Description: "Idk man",
		}

		// create blog.
		MustCreateBlog(t, db, adminUsrCtx, blog)

		subBlog := &pa.SubBlog{
			BlogID:  blog.ID,
			Title:   "Cool Sub blog",
			Content: "idk",
		}

		// create sub blog.
		MustCreateSubBlog(t, db, adminUsrCtx, subBlog)

		comment := &pa.Comment{
			SubBlogID: subBlog.ID,
			Content:   "Cool content",
		}

		// create comment.
		MustCreateComment(t, db, adminUsrCtx, comment)

		comment2 := &pa.Comment{
			SubBlogID: subBlog.ID,
			Content:   "BAd Content",
		}

		// create comment.
		MustCreateComment(t, db, adminUsrCtx, comment2)

		// find comments.
		if gotComments, n, err := commentService.FindComments(backgroundCtx, pa.CommentFilter{UserID: NewIntPointer(user.ID)}); err != nil {
			t.Fatal(err)
		} else if len(gotComments) != 2 { // assert find
			t.Fatalf("len=%v != 2", len(gotComments))
		} else if n != 2 {
			t.Fatalf("n=%v != 2", n)
		}
	})

	t.Run("Bad Find Call (Not Found)", func(t *testing.T) {
		db := MustOpenTempDB(t)
		defer MustCloseDB(t, db)

		backgroundCtx := context.Background()

		commentService := sqlite.NewCommentService(db)

		if _, err := commentService.FindCommentByID(backgroundCtx, 1); pa.ErrorCode(err) != pa.ENOTFOUND {
			t.Fatal("err != ENOTFOUND")
		}
	})
}

func MustCreateComment(t *testing.T, db *sqlite.DB, ctx context.Context, comment *pa.Comment) {
	t.Helper()
	if err := sqlite.NewCommentService(db).CreateComment(ctx, comment); err != nil {
		t.Fatal(err)
	}
}
