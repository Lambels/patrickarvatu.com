package sqlite_test

import (
	"context"
	"testing"

	pa "github.com/Lambels/patrickarvatu.com"
	"github.com/Lambels/patrickarvatu.com/sqlite"
)

// TODO: write test

func MustCreateSubBlog(t *testing.T, db *sqlite.DB, ctx context.Context, subBlog *pa.SubBlog) {
	t.Helper()
	if err := sqlite.NewSubBlogService(db).CreateSubBlog(ctx, subBlog); err != nil {
		t.Fatal(err)
	}
}
