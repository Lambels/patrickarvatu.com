package sqlite_test

import (
	"context"
	"reflect"
	"testing"

	pa "github.com/Lambels/patrickarvatu.com"
	"github.com/Lambels/patrickarvatu.com/sqlite"
)

func TestCreateSubBlog(t *testing.T) {
	t.Run("Ok Create Call", func(t *testing.T) {
		db := MustOpenTempDB(t)
		defer MustCloseDB(t, db)

		backgroundCtx := context.Background()

		subBlogService := sqlite.NewSubBlogService(db)

		user := &pa.User{
			Name:    "Jhon Doe",
			Email:   "jhon@doe.com",
			IsAdmin: true,
		} // no need to create user as CreateBlog doesent check any keys.

		adminUsrCtx := pa.NewContextWithUser(backgroundCtx, user)

		blog := &pa.Blog{
			Title:       "Epic Blog",
			Description: "Honestly the best blog ever.",
		}

		// create blog.
		MustCreateBlog(t, db, adminUsrCtx, blog)

		subBlog := &pa.SubBlog{
			BlogID:  blog.ID,
			Title:   "some title",
			Content: "some content",
		}

		// create sub blog.
		if err := subBlogService.CreateSubBlog(adminUsrCtx, subBlog); err != nil {
			t.Fatal(err)
		} else if subBlog.ID == 0 {
			t.Fatal("got id = 0")
		} else if subBlog.CreatedAt.IsZero() {
			t.Fatal("expected sub blog time stamp creation (created AT)")
		} else if subBlog.UpdatedAt.IsZero() {
			t.Fatal("expected sub blog time stamp creation (updated AT)")
		}

		// assert creation.
		if gotSubBlog, err := subBlogService.FindSubBlogByID(backgroundCtx, 1); err != nil {
			t.Fatal(err)
		} else if !reflect.DeepEqual(gotSubBlog, subBlog) {
			t.Fatal("DeepEqual: gotSubBlog != SubBlog")
		}
	})

	t.Run("Bad Create Call (Un Auth)", func(t *testing.T) {
		db := MustOpenTempDB(t)
		defer MustCloseDB(t, db)

		backgroundCtx := context.Background()

		subBlogService := sqlite.NewSubBlogService(db)

		user := &pa.User{
			Name:    "Jhon Doe",
			Email:   "jhon@doe.com",
			IsAdmin: true,
		} // no need to create user as CreateBlog doesent check any keys.

		adminUsrCtx := pa.NewContextWithUser(backgroundCtx, user)

		user2 := &pa.User{
			Name:  "Lamb",
			Email: "lamb@bels.com",
		} // no need to create user as CreateBlog doesent check any keys.

		usrCtx := pa.NewContextWithUser(backgroundCtx, user2)

		blog := &pa.Blog{
			Title:       "Epic Blog",
			Description: "Honestly the best blog ever.",
		}

		// create blog.
		MustCreateBlog(t, db, adminUsrCtx, blog)

		subBlog := &pa.SubBlog{
			BlogID:  blog.ID,
			Title:   "some title",
			Content: "some content",
		}

		// create sub blog (un auth).
		if err := subBlogService.CreateSubBlog(usrCtx, subBlog); pa.ErrorCode(err) != pa.EUNAUTHORIZED {
			t.Fatal("expected UnAuth error")
		} else if subBlog.ID != 0 {
			t.Fatal("got id != 0")
		}
	})
}

func TestDeleteSubBlog(t *testing.T) {
	t.Run("Ok Delete Call", func(t *testing.T) {
		db := MustOpenTempDB(t)
		defer MustCloseDB(t, db)

		backgroundCtx := context.Background()

		subBlogService := sqlite.NewSubBlogService(db)

		user := &pa.User{
			Name:    "Jhon Doe",
			Email:   "jhon@doe.com",
			IsAdmin: true,
		} // no need to create user as DeleteBlog doesent check any keys.

		adminUsrCtx := pa.NewContextWithUser(backgroundCtx, user)

		blog := &pa.Blog{
			Title:       "Epic Blog",
			Description: "Honestly the best blog ever.",
		}

		// create blog.
		MustCreateBlog(t, db, adminUsrCtx, blog)

		subBlog := &pa.SubBlog{
			BlogID:  blog.ID,
			Title:   "some title",
			Content: "some content",
		}

		// create sub blog.
		MustCreateSubBlog(t, db, adminUsrCtx, subBlog)

		// delete sub blog.
		if err := subBlogService.DeleteSubBlog(adminUsrCtx, 1); err != nil {
			t.Fatal(err)
		}

		// assert deletion.
		if _, err := subBlogService.FindSubBlogByID(backgroundCtx, subBlog.ID); pa.ErrorCode(err) != pa.ENOTFOUND {
			t.Fatal("err != ENOTFOUND")
		}
	})

	t.Run("Bad Delete Call (Un Auth)", func(t *testing.T) {
		db := MustOpenTempDB(t)
		defer MustCloseDB(t, db)

		backgroundCtx := context.Background()

		subBlogService := sqlite.NewSubBlogService(db)

		user := &pa.User{
			Name:    "Jhon Doe",
			Email:   "jhon@doe.com",
			IsAdmin: true,
		} // no need to create user as DeleteBlog doesent check any keys.

		adminUsrCtx := pa.NewContextWithUser(backgroundCtx, user)

		user2 := &pa.User{
			Name:  "Lambels",
			Email: "lambi@lam.com",
		} // no need to create user as DeleteBlog doesent check any keys.

		usrCtx := pa.NewContextWithUser(backgroundCtx, user2)

		blog := &pa.Blog{
			Title:       "Epic Blog",
			Description: "Honestly the best blog ever.",
		}

		// create blog.
		MustCreateBlog(t, db, adminUsrCtx, blog)

		subBlog := &pa.SubBlog{
			BlogID:  blog.ID,
			Title:   "some title",
			Content: "some content",
		}

		// create sub blog.
		MustCreateSubBlog(t, db, adminUsrCtx, subBlog)

		// delete sub blog (Un Auth).
		if err := subBlogService.DeleteSubBlog(usrCtx, subBlog.ID); pa.ErrorCode(err) != pa.EUNAUTHORIZED {
			t.Fatal("err != EUNAUTHORIZED")
		}
	})

	t.Run("Bad Delete Call (Not Found)", func(t *testing.T) {
		db := MustOpenTempDB(t)
		defer MustCloseDB(t, db)

		backgroundCtx := context.Background()

		subBlogService := sqlite.NewSubBlogService(db)

		user := &pa.User{
			Name:    "Jhon Doe",
			Email:   "jhon@doe.com",
			IsAdmin: true,
		} // no need to create user as DeleteBlog doesent check any keys.

		adminUsrCtx := pa.NewContextWithUser(backgroundCtx, user)

		// delete sub blog (Not Found).
		if err := subBlogService.DeleteSubBlog(adminUsrCtx, 134); pa.ErrorCode(err) != pa.ENOTFOUND {
			t.Fatal("err != ENOTFOUND")
		}
	})
}

func TestUpdateSubBlog(t *testing.T) {
	t.Run("Ok Update Call", func(t *testing.T) {
		db := MustOpenTempDB(t)
		defer MustCloseDB(t, db)

		backgroundCtx := context.Background()

		subBlogService := sqlite.NewSubBlogService(db)

		user := &pa.User{
			Name:    "Jhon Doe",
			Email:   "jhon@doe.com",
			IsAdmin: true,
		} // no need to create user as UpdateBlog doesent check any keys.

		adminUsrCtx := pa.NewContextWithUser(backgroundCtx, user)

		blog := &pa.Blog{
			Title:       "Epic Blog",
			Description: "Honestly the best blog ever.",
		}

		// create blog.
		MustCreateBlog(t, db, adminUsrCtx, blog)

		subBlog := &pa.SubBlog{
			BlogID:  blog.ID,
			Title:   "some title",
			Content: "some content",
		}

		// create sub blog.
		MustCreateSubBlog(t, db, adminUsrCtx, subBlog)

		subBlogUpdate := pa.SubBlogUpdate{
			Title:   NewStringPointer("other title"),
			Content: NewStringPointer("other content"),
		}

		// update sub blog.
		if updatedSubBlog, err := subBlogService.UpdateSubBlog(adminUsrCtx, subBlog.ID, subBlogUpdate); err != nil {
			t.Fatal(err)
		} else if gotSubBlog, err := subBlogService.FindSubBlogByID(backgroundCtx, subBlog.ID); err != nil { // assert update.
			t.Fatal(err)
		} else if !reflect.DeepEqual(updatedSubBlog, gotSubBlog) {
			t.Log(*updatedSubBlog, *gotSubBlog)
			t.Fatal("DeepEqual: updatedSubBlog != gotSubBlog")
		}
	})

	t.Run("Bad Update Call", func(t *testing.T) {
		db := MustOpenTempDB(t)
		defer MustCloseDB(t, db)

		backgroundCtx := context.Background()

		subBlogService := sqlite.NewSubBlogService(db)

		user := &pa.User{
			Name:    "Jhon Doe",
			Email:   "jhon@doe.com",
			IsAdmin: true,
		} // no need to create user as UpdateBlog doesent check any keys.

		adminUsrCtx := pa.NewContextWithUser(backgroundCtx, user)

		user2 := &pa.User{
			Name:  "Lambels",
			Email: "lambi@lam.com",
		} // no need to create user as DeleteBlog doesent check any keys.

		usrCtx := pa.NewContextWithUser(backgroundCtx, user2)

		blog := &pa.Blog{
			Title:       "Epic Blog",
			Description: "Honestly the best blog ever.",
		}

		// create blog.
		MustCreateBlog(t, db, adminUsrCtx, blog)

		subBlog := &pa.SubBlog{
			BlogID:  blog.ID,
			Title:   "some title",
			Content: "some content",
		}

		// create sub blog.
		MustCreateSubBlog(t, db, adminUsrCtx, subBlog)

		subBlogUpdate := pa.SubBlogUpdate{
			Title:   NewStringPointer("other title"),
			Content: NewStringPointer("other content"),
		}

		// update sub blog (Un Auth).
		if _, err := subBlogService.UpdateSubBlog(usrCtx, subBlog.ID, subBlogUpdate); pa.ErrorCode(err) != pa.EUNAUTHORIZED {
			t.Fatal("err != EUNAUTHORIZED")
		}
	})

	t.Run("Bad Update Call (Not Found)", func(t *testing.T) {
		db := MustOpenTempDB(t)
		defer MustCloseDB(t, db)

		backgroundCtx := context.Background()

		subBlogService := sqlite.NewSubBlogService(db)

		user := &pa.User{
			Name:    "Jhon Doe",
			Email:   "jhon@doe.com",
			IsAdmin: true,
		} // no need to create user as UpdateBlog doesent check any keys.

		adminUsrCtx := pa.NewContextWithUser(backgroundCtx, user)

		subBlogUpdate := pa.SubBlogUpdate{
			Title:   NewStringPointer("other title"),
			Content: NewStringPointer("other content"),
		}

		// update sub blog (Not Found).
		if _, err := subBlogService.UpdateSubBlog(adminUsrCtx, 3, subBlogUpdate); pa.ErrorCode(err) != pa.ENOTFOUND {
			t.Fatal("err != ENOTFOUND")
		}
	})
}

func TestFindSubBlogs(t *testing.T) {
	t.Run("Ok Find Call (filter - blogID)", func(t *testing.T) {
		db := MustOpenTempDB(t)
		defer MustCloseDB(t, db)

		backgroundCtx := context.Background()

		subBlogService := sqlite.NewSubBlogService(db)

		user := &pa.User{
			Name:    "Jhon Doe",
			Email:   "jhon@doe.com",
			IsAdmin: true,
		}

		adminUsrCtx := MustCreateUser(t, db, backgroundCtx, user)

		blog := &pa.Blog{
			Title:       "Epic Blog",
			Description: "Honestly the best blog ever.",
		}

		// create blog.
		MustCreateBlog(t, db, adminUsrCtx, blog)

		subBlog := &pa.SubBlog{
			BlogID:  blog.ID,
			Title:   "some title 1",
			Content: "some title 2",
		}

		// create sub blog.
		MustCreateSubBlog(t, db, adminUsrCtx, subBlog)

		subBlog2 := &pa.SubBlog{
			BlogID:  blog.ID,
			Title:   "some title 2",
			Content: "some title 2",
		}

		// create sub blog.
		MustCreateSubBlog(t, db, adminUsrCtx, subBlog2)

		comment := &pa.Comment{
			SubBlogID: subBlog.ID,
			UserID:    user.ID,
			Content:   "some content 1",
		}

		// create comment.
		MustCreateComment(t, db, adminUsrCtx, comment)

		comment2 := &pa.Comment{
			SubBlogID: subBlog.ID,
			UserID:    user.ID,
			Content:   "some content 2",
		}

		// create comment.
		MustCreateComment(t, db, adminUsrCtx, comment2)

		comment3 := &pa.Comment{
			SubBlogID: subBlog2.ID,
			UserID:    user.ID,
			Content:   "some content 3",
		}

		// create comment.
		MustCreateComment(t, db, adminUsrCtx, comment3)

		subBlogFilter := pa.SubBlogFilter{
			BlogID: NewIntPointer(blog.ID),
		}

		// find sub blogs.
		if gotSubBlogs, n, err := subBlogService.FindSubBlogs(backgroundCtx, subBlogFilter); err != nil {
			t.Fatal(err)
		} else if len(gotSubBlogs) != 2 { // assert find.
			t.Fatalf("len=%v != 2", gotSubBlogs)
		} else if n != 2 {
			t.Fatalf("n=%v != 2", n)
		} else if len(gotSubBlogs[0].Comments) != 2 {
			t.Fatalf("len comments 1=%v != 2", len(gotSubBlogs[0].Comments))
		} else if len(gotSubBlogs[1].Comments) != 1 {
			t.Fatalf("len comments 2=%v != 1", len(gotSubBlogs[1].Comments))
		}
	})

	t.Run("Bad Find Call (Not Found)", func(t *testing.T) {
		db := MustOpenTempDB(t)
		defer MustCloseDB(t, db)

		backgroundCtx := context.Background()

		subBlogService := sqlite.NewSubBlogService(db)

		// find sub blog (Not Found).
		if _, err := subBlogService.FindSubBlogByID(backgroundCtx, 1); pa.ErrorCode(err) != pa.ENOTFOUND {
			t.Fatal("err != ENOTFOUND")
		}
	})
}

func MustCreateSubBlog(t *testing.T, db *sqlite.DB, ctx context.Context, subBlog *pa.SubBlog) {
	t.Helper()
	if err := sqlite.NewSubBlogService(db).CreateSubBlog(ctx, subBlog); err != nil {
		t.Fatal(err)
	}
}
