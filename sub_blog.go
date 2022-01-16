package pa

import (
	"context"
	"time"
)

// SubBlog represents an sub blog object in the system.
// SubBlog has no reason to store any user ID as the admin user is the only one
// who can interact with SubBlogService.CreateSubBlog().
type SubBlog struct {
	// the pk of the sub blog.
	ID int `json:"id"`

	// the id of the blog under which the sub blog is.
	BlogID int `json:"blogID"`

	// the descriptive fields of the sub blog.
	Title    string     `json:"title"`
	Content  string     `json:"body"`
	Comments []*Comment `json:"comments"`

	// timestamps.
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Validate performs basic validation on the sub blog.
// returns EINVALID if any error is found.
func (s *SubBlog) Validate() error {
	if s.Title == "" {
		return Errorf(EINVALID, "title is a required field.")
	}
	if s.Content == "" {
		return Errorf(EINVALID, "content is a required field.")
	}
	if s.BlogID == 0 {
		return Errorf(EINVALID, "sub blog must be linked to blog.")
	}

	return nil
}

// SubBlogService represents a service which manages auth in the system.
type SubBlogService interface {
	// FindSubBlogByID returns a sub blog based on the id.
	// returns ENOTFOUND if the sub blog doesent exist.
	FindSubBlogByID(ctx context.Context, id int) (*SubBlog, error)

	// FindSubBlogs returns a range of sub blogs and the length of the range. If filter
	// is specified FindSubBlogs will apply the filter to return set response.
	FindSubBlogs(ctx context.Context, filter SubBlogFilter) ([]*SubBlog, int, error)

	// CreateSubBlog creates a sub blog.
	// returns EUNAUTHORIZED if used by anyone other then the adim user.
	CreateSubBlog(ctx context.Context, subBlog *SubBlog) error

	// UpdateSubBlog updates a sub blog based on the update field.
	// returns ENOTFOUND if sub blog doesent exist.
	// returns EUNAUTHORIZED if used by anyone other then the adim user.
	UpdateSubBlog(ctx context.Context, id int, update SubBlogUpdate) (*SubBlog, error)

	// DeleteSubBlog permanently deletes a sub blog.
	// returns ENOTFOUND if sub blog doesent exist.
	// returns EUNAUTHORIZED if used by anyone other then the adim user.
	DeleteSubBlog(ctx context.Context, id int) error
}

// SubBlogFilter represents a filter used by FindSubBlogs to filter the response.
type SubBlogFilter struct {
	// fields to filter on.
	ID     *int    `json:"id"`
	Title  *string `json:"title"`
	BlogID *int    `json:"blogID"`

	// restrictions on the result set, used for pagination and set limits.
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

// SubBlogUpdate represents an update used by UpdateSubBlog to update a sub blog.
type SubBlogUpdate struct {
	// fields which can be updated.
	Title   *string `json:"title"`
	Content *string `json:"content"`
}
