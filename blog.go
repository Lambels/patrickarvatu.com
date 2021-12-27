package pa

import (
	"context"
	"time"
)

type Blog struct {
	ID int `json:"id"`

	Title       string     `json:"title"`
	Description string     `json:"description"`
	SubBlogs    []*SubBlog `json:"subBlogs"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (b *Blog) Validate() error {
	if b.Title == "" {
		return Errorf(EINVALID, "title is a required field.")
	}
	return nil
}

type BlogService interface {
	FindBlogByID(ctx context.Context, id int) (*Blog, error)

	FindBlogs(ctx context.Context, filter BlogFilter) ([]*Blog, int, error)

	CreateBlog(ctx context.Context, blog *Blog) error

	UpdateBlog(ctx context.Context, update BlogUpdate) (*Blog, error)

	DeleteBlog(ctx context.Context, id int) error
}

type BlogFilter struct {
	ID    *int    `json:"id"`
	Title *string `json:"title"`

	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type BlogUpdate struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}
