package pa

import (
	"context"
	"time"
)

// Blog represents an blog object in the system.
// Blog has no reason to store any user ID as the admin user is the only one
// who can interact with BlogService.CreateBlog().
type Blog struct {
	// the pk of the blog.
	ID int `json:"id"`

	// the descriptive fields of the blog.
	Title       string     `json:"title"`
	Description string     `json:"description"`
	SubBlogs    []*SubBlog `json:"subBlogs"` // the list of sub blogs contained by the blog.

	// timestamps.
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Validate performs basic validation on the blog.
// returns EINVALID if any error is found.
func (b *Blog) Validate() error {
	if b.Title == "" {
		return Errorf(EINVALID, "title is a required field.")
	}
	return nil
}

// BlogService represents a service which manages auth in the system.
type BlogService interface {
	// FindBlogByID returns a blog based on the id.
	// returns ENOTFOUND if the blog doesent exist.
	FindBlogByID(ctx context.Context, id int) (*Blog, error)

	// FindBlogs returns a range of blogs and the length of the range. If filter
	// is specified FindBlogs will apply the filter to return set response.
	FindBlogs(ctx context.Context, filter BlogFilter) ([]*Blog, int, error)

	// CreateBlog creates a blog.
	// returns EUNAUTHORIZED if used by anyone other then the adim user.
	CreateBlog(ctx context.Context, blog *Blog) error

	// UpdateBlog updates a blog based on the update field.
	// returns ENOTFOUND if blog doesent exist.
	// returns EUNAUTHORIZED if used by anyone other then the adim user.
	UpdateBlog(ctx context.Context, id int, update BlogUpdate) (*Blog, error)

	// DeleteBlog permanently deletes a blog.
	// returns ENOTFOUND if blog doesent exist.
	// returns EUNAUTHORIZED if used by anyone other then the adim user.
	DeleteBlog(ctx context.Context, id int) error
}

// BlogFilter represents a filter used by FindBlogs to filter the response.
type BlogFilter struct {
	// fields to filter on.
	ID    *int    `json:"id"`
	Title *string `json:"title"`

	// restrictions on the result set, used for pagination and set limits.
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

// BlogUpdate represents an update used by UpdateBlog to update a blog.
type BlogUpdate struct {
	// fields which can be updated.
	Title       *string `json:"title"`
	Description *string `json:"description"`
}
