package sqlite_test

import (
	"context"
	"reflect"
	"testing"

	pa "github.com/Lambels/patrickarvatu.com"
	"github.com/Lambels/patrickarvatu.com/sqlite"
)

func TestCreateBlog(t *testing.T) {
	t.Parallel()
	t.Run("Ok Create Call", func(t *testing.T) {
		db := MustOpenTempDB(t)
		defer MustCloseDB(t, db)

		backgroundCtx := context.Background()

		blogService := sqlite.NewBlogService(db)

		user := &pa.User{
			Name:    "Jhon Doe",
			Email:   "jhon@doe.com",
			IsAdmin: true,
		} // no need to create user as CreateBlog doesent check any keys.

		adminUsrContext := pa.NewContextWithUser(backgroundCtx, user)

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

	t.Run("Bad Create Call (Un Auth)", func(t *testing.T) {
		db := MustOpenTempDB(t)
		defer MustCloseDB(t, db)

		backgroundCtx := context.Background()
		usrContext := pa.NewContextWithUser(backgroundCtx, &pa.User{
			Name:  "jhon DOE",
			Email: "jhon@doe.com",
		}) // no need to create user as CreateBlog doesent check any keys.

		blogService := sqlite.NewBlogService(db)

		blog := &pa.Blog{
			Title:       "Epic Blog",
			Description: "Honestly the best blog ever.",
		}

		// create blog (Un Auth).
		if err := blogService.CreateBlog(usrContext, blog); pa.ErrorCode(err) != pa.EUNAUTHORIZED {
			t.Fatal("expected UnAuth error")
		} else if blog.ID != 0 {
			t.Fatal("got id != 0")
		}
	})
}

// TODO: add testing for rest of CRUD methods.
func TestDeleteBlog(t *testing.T) {
	t.Parallel()
	t.Run("Ok Delete Call", func(t *testing.T) {
		db := MustOpenTempDB(t)
		defer MustCloseDB(t, db)

		backgroundCtx := context.Background()

		blogService := sqlite.NewBlogService(db)

		user := &pa.User{
			Name:    "Jhon Doe",
			Email:   "jhon@doe.com",
			IsAdmin: true,
		} // no need to create user as DeleteBlog doesent check any keys.

		adminUsrContext := pa.NewContextWithUser(backgroundCtx, user)

		blog := &pa.Blog{
			Title:       "Epic Blog",
			Description: "Honestly the best blog ever.",
		}

		// create blog.
		MustCreateBlog(t, db, adminUsrContext, blog)

		// delete blog.
		if err := blogService.DeleteBlog(adminUsrContext, 1); err != nil {
			t.Fatal(err)
		}

		// assert deletion.
		if _, err := blogService.FindBlogByID(backgroundCtx, 1); pa.ErrorCode(err) != pa.ENOTFOUND {
			t.Fatal("err != ENOTFOUND")
		}
	})

	t.Run("Bad Delete Call (Un Auth)", func(t *testing.T) {
		db := MustOpenTempDB(t)
		defer MustCloseDB(t, db)

		backgroundCtx := context.Background()

		blogService := sqlite.NewBlogService(db)

		user := &pa.User{
			Name:    "Jhon Doe",
			Email:   "jhon@doe.com",
			IsAdmin: true,
		} // no need to create user as DeleteBlog doesent check any keys.

		user2 := &pa.User{
			Name:  "Lambels",
			Email: "Lamb@Lambels.com",
		}

		adminUsrContext := pa.NewContextWithUser(backgroundCtx, user)
		user2Context := pa.NewContextWithUser(backgroundCtx, user2)

		blog := &pa.Blog{
			Title:       "Epic Blog",
			Description: "Honestly the best blog ever.",
		}

		// create blog.
		MustCreateBlog(t, db, adminUsrContext, blog)

		// delete blog.
		if err := blogService.DeleteBlog(user2Context, 1); pa.ErrorCode(err) != pa.EUNAUTHORIZED {
			t.Fatal("err != EUNAUTHORIZED")
		}
	})

	t.Run("Bad Delete Call (Not Found)", func(t *testing.T) {
		db := MustOpenTempDB(t)
		defer MustCloseDB(t, db)

		backgroundCtx := context.Background()

		blogService := sqlite.NewBlogService(db)

		user := &pa.User{
			Name:    "Jhon Doe",
			Email:   "jhon@doe.com",
			IsAdmin: true,
		} // no need to create user as DeleteBlog doesent check any keys.

		adminUsrContext := pa.NewContextWithUser(backgroundCtx, user)

		// delete blog (Not Found).
		if err := blogService.DeleteBlog(adminUsrContext, 1); pa.ErrorCode(err) != pa.ENOTFOUND {
			t.Fatal("err != ENOTFOUND")
		}
	})
}

func TestUpdateBlog(t *testing.T) {
	t.Parallel()
	t.Run("Ok Update Call", func(t *testing.T) {
		db := MustOpenTempDB(t)
		defer MustCloseDB(t, db)

		backgroundCtx := context.Background()

		blogService := sqlite.NewBlogService(db)

		user := &pa.User{
			Name:    "Jhon Doe",
			Email:   "jhon@doe.com",
			IsAdmin: true,
		} // no need to create user as UpdateBlog doesent check any keys.

		adminUsrContext := pa.NewContextWithUser(backgroundCtx, user)

		blog := &pa.Blog{
			Title:       "Epic Blog",
			Description: "Honestly the best blog ever.",
		}
		update := pa.BlogUpdate{
			Title:       NewStringPointer("Bad Blog"),
			Description: NewStringPointer("Honestly worst blog ever"),
		}

		// create blog.
		MustCreateBlog(t, db, adminUsrContext, blog)

		// update blog.
		if updatedBlog, err := blogService.UpdateBlog(adminUsrContext, 1, update); err != nil {
			t.Fatal(err)
		} else if gotBlog, err := blogService.FindBlogByID(backgroundCtx, 1); err != nil { // assert update.
			t.Fatal(err)
		} else if !reflect.DeepEqual(updatedBlog, gotBlog) {
			t.Log(*updatedBlog, *gotBlog)
			t.Fatal("DeepEqual: updatedBlog != gotBlog")
		}
	})

	t.Run("Bad Update Call (Un Auth)", func(t *testing.T) {
		db := MustOpenTempDB(t)
		defer MustCloseDB(t, db)

		backgroundCtx := context.Background()

		blogService := sqlite.NewBlogService(db)

		user := &pa.User{
			Name:    "Jhon Doe",
			Email:   "jhon@doe.com",
			IsAdmin: true,
		} // no need to create user as DeleteBlog doesent check any keys.

		user2 := &pa.User{
			Name:  "Lambels",
			Email: "Lamb@Lambels.com",
		}

		adminUsrContext := pa.NewContextWithUser(backgroundCtx, user)
		user2Context := pa.NewContextWithUser(backgroundCtx, user2)

		blog := &pa.Blog{
			Title:       "Epic Blog",
			Description: "Honestly the best blog ever.",
		}
		update := pa.BlogUpdate{
			Title:       NewStringPointer("Bad Blog"),
			Description: NewStringPointer("Honestly worst blog ever"),
		}

		// create blog.
		MustCreateBlog(t, db, adminUsrContext, blog)

		// update blog (Un Auth).
		if _, err := blogService.UpdateBlog(user2Context, 1, update); pa.ErrorCode(err) != pa.EUNAUTHORIZED {
			t.Fatal("err != UNAUTHORIZED")
		}
	})

	t.Run("Bad Update Call (Not Found)", func(t *testing.T) {
		db := MustOpenTempDB(t)
		defer MustCloseDB(t, db)

		backgroundCtx := context.Background()

		blogService := sqlite.NewBlogService(db)

		user := &pa.User{
			Name:    "Jhon Doe",
			Email:   "jhon@doe.com",
			IsAdmin: true,
		} // no need to create user as DeleteBlog doesent check any keys.

		adminUsrContext := pa.NewContextWithUser(backgroundCtx, user)

		update := pa.BlogUpdate{
			Title:       NewStringPointer("Bad Blog"),
			Description: NewStringPointer("Honestly worst blog ever"),
		}

		// update blog (Un Auth).
		if _, err := blogService.UpdateBlog(adminUsrContext, 1, update); pa.ErrorCode(err) != pa.ENOTFOUND {
			t.Fatal("err != ENOTFOUND")
		}
	})
}

func TestFindBlogs(t *testing.T) {
	t.Run("Ok Find Call (filter - id)", func(t *testing.T) {

	})

	t.Run("Ok Find Call (filter - title)", func(t *testing.T) {

	})

	t.Run("Bad Find Call (Not Found)", func(t *testing.T) {

	})
}

func MustCreateBlog(t *testing.T, db *sqlite.DB, ctx context.Context, blog *pa.Blog) {
	t.Helper()
	if err := sqlite.NewBlogService(db).CreateBlog(ctx, blog); err != nil {
		t.Fatal(err)
	}
}
