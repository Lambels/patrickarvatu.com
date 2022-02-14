package pa

import (
	"context"
	"time"
)

// Comment represents a comment in the system.
type Comment struct {
	// the pk of the comment.
	ID int `json:"id"`

	// linking fields of the comment.
	SubBlogID int   `json:"subBlogID"`
	UserID    int   `json:"userID"`
	User      *User `json:"user"`

	// content of the comment.
	Content string `json:"content"`

	// timestamp.
	CreatedAt time.Time `json:"createdAt"`
}

// Validate performs basic validation on the comment.
// returns EINVALID if any error is found.
func (c *Comment) Validate() error {
	if c.Content == "" {
		return Errorf(EINVALID, "content is a required field.")
	}
	if c.SubBlogID == 0 {
		return Errorf(EINVALID, "comment must be linked to a sub blog.")
	}
	if c.UserID == 0 {
		return Errorf(EINVALID, "comment must be linked to a user.")
	}

	return nil
}

// CommentService represents a service which manages comments in the system.
type CommentService interface {
	// FindCommentByID returns a comment based on the id.
	// returns ENOTFOUND if the comment doesent exist.
	FindCommentByID(ctx context.Context, id int) (*Comment, error)

	// FindComments returns a range of comments and the length of the range. If filter
	// is specified FindComments will apply the filter to return set response.
	FindComments(ctx context.Context, filter CommentFilter) ([]*Comment, int, error)

	// CreateComment creates a comment.
	CreateComment(ctx context.Context, comment *Comment) error

	// UpdateComment updates a comment based on the update field.
	// returns ENOTFOUND if the comment doesent exist.
	// returns EUNAUTHORIZED if used by anyone other then the adim user.
	UpdateComment(ctx context.Context, id int, update CommentUpdate) (*Comment, error)

	// DeleteComment permanently deletes a comment.
	// returns ENOTFOUND if comment doesent exist.
	// returns EUNAUTHORIZED if used by anyone other then the user owning the comment.
	DeleteComment(ctx context.Context, id int) error
}

// CommentFilter represents a filter used by FindComments to filter the response.
type CommentFilter struct {
	// fields to filter on.
	ID        *int `json:"id"`
	SubBlogID *int `json:"SubBlogID"`
	UserID    *int `json:"userID"`

	// restrictions on the result set, used for pagination and set limits.
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

// CommentUpdate represents an update used by UpdateComment to update a comment.
type CommentUpdate struct {
	// fields which can be updated.
	Content *string `json:"content"`
}
