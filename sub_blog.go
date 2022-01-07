package pa

import (
	"context"
	"time"
)

type SubBlog struct {
	ID int `json:"id"`

	BlogID int `json:"blogID"`

	Title    string     `json:"title"`
	Content  string     `json:"body"`
	Comments []*Comment `json:"comments"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (s *SubBlog) Validate() error {
	if s.Title == "" {
		return Errorf(EINVALID, "title is a required field.")
	}
	if s.Content == "" {
		return Errorf(EINVALID, "content is a required field.")
	}

	return nil
}

// SubBlogService represents a service which manages auth in the system.
type SubBlogService interface {
	FindSubBlogByID(ctx context.Context, id int) (*SubBlog, error)

	FindSubBlogs(ctx context.Context, filter SubBlogFilter) ([]*SubBlog, int, error)

	CreateSubBlog(ctx context.Context, subBlog *SubBlog) error

	UpdateSubBlog(ctx context.Context, id int, update SubBlogUpdate) (*SubBlog, error)

	DeleteSubBlog(ctx context.Context, id int) error
}

type SubBlogFilter struct {
	ID     *int    `json:"id"`
	Title  *string `json:"title"`
	BlogID *int    `json:"blogID"`

	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type SubBlogUpdate struct {
	Title   *string `json:"title"`
	Content *string `json:"content"`
}
