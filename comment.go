package pa

import (
	"context"
	"time"
)

type Comment struct {
	ID int `json:"id"`

	SubBlogID int `json:"subBlogID"`
	UserID    int `json:"userID"`

	Content string `json:"content"`

	CreatedAt time.Time `json:"createdAt"`
}

func (c *Comment) Validate() error {
	if c.Content == "" {
		return Errorf(EINVALID, "content is a required field.")
	}

	return nil
}

type CommentService interface {
	FindCommentByID(ctx context.Context, id int) (*Comment, error)

	FindComments(ctx context.Context, filter CommentFilter) ([]*Comment, int, error)

	CreateComment(ctx context.Context, blog *Comment) error

	// under admin control
	UpdateComment(ctx context.Context, update CommentUpdate) (*Comment, error)

	DeleteComment(ctx context.Context, id int) error
}

type CommentFilter struct {
	ID        *int `json:"id"`
	SubBlogID *int `json:"SubBlogID"`
	UserID    *int `json:"userID"`

	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type CommentUpdate struct {
	Content *string `json:"content"`
}
