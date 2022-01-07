package sqlite

import (
	"context"

	pa "github.com/Lambels/patrickarvatu.com"
)

// check to see if *CommentService object implements set interface.
var _ pa.CommentService = (*CommentService)(nil)

// CommentService represents a service used to manage comments.
type CommentService struct {
	db *DB
}

// NewCommentService returns a new instance of CommentService attached to db.
func NewCommentService(db *DB) *CommentService {
	return &CommentService{
		db: db,
	}
}

// FindCommentByID returns a comment based on id.
// returns ENOTFOUND if the comment doesent exist.
func (s *CommentService) FindCommentByID(ctx context.Context, id int) (*pa.Comment, error) {
	tx, err := s.db.BeginTX(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	comment, err := findCommentByID(ctx, tx, id)
	if err != nil {
		return nil, err

	} else if err := attachUserToComment(ctx, tx, comment); err != nil {
		return nil, err
	}

	return comment, nil
}

// FindComments returns a range of comment based on filter.
func (s *CommentService) FindComments(ctx context.Context, filter pa.CommentFilter) ([]*pa.Comment, int, error) {
	tx, err := s.db.BeginTX(ctx, nil)
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback()

	comments, n, err := findComments(ctx, tx, filter)
	if err != nil {
		return comments, n, err
	}

	for _, comment := range comments {
		if err := attachUserToComment(ctx, tx, comment); err != nil {
			return comments, n, err
		}
	}
	return comments, n, nil
}

// CreateComment creates a new comment.
func (s *CommentService) CreateComment(ctx context.Context, comment *pa.Comment) error {
	tx, err := s.db.BeginTX(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := createComment(ctx, tx, comment); err != nil {
		return err

	} else if err := attachUserToComment(ctx, tx, comment); err != nil {
		return err
	}
	return tx.Commit()
}

// UpdateComment updates comment with id: id.
// returns EUNAUTHORIZED if the user isnt trying to update isnt the admin user.
// returns ENOTFOUND if the comment doesent exist.
func (s *CommentService) UpdateComment(ctx context.Context, id int, update pa.CommentUpdate) (*pa.Comment, error) {
	tx, err := s.db.BeginTX(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	comment, err := updateComment(ctx, tx, id, update)
	if err != nil {
		return nil, err

	} else if err := attachUserToComment(ctx, tx, comment); err != nil {
		return nil, err
	}
	return comment, tx.Commit()
}

// DeleteComment permanently deletes the comment specified by id.
// returns EUNAUTHORIZED if the user isnt trying to delete his own comment.
// returns ENOTFOUND if the comment doesent exist.
func (s *CommentService) DeleteComment(ctx context.Context, id int) error {
	tx, err := s.db.BeginTX(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := deleteComment(ctx, tx, id); err != nil {
		return err
	}
	return tx.Commit()
}

func findCommentByID(ctx context.Context, tx *Tx, id int) (*pa.Comment, error) {
	filter := pa.CommentFilter{
		ID: &id,
	}
	comments, _, err := findComments(ctx, tx, filter)
	if err != nil {
		return nil, err

	} else if len(comments) == 0 {
		return nil, pa.Errorf(pa.ENOTFOUND, "comment not found")
	}
	return comments[0], nil
}

func findComments(ctx context.Context, tx *Tx, filter pa.CommentFilter) ([]*pa.Comment, int, error) {

}

func attachUserToComment(ctx context.Context, tx *Tx, comment *pa.Comment) error {

}

func createComment(ctx context.Context, tx *Tx, comment *pa.Comment) error {

}

func updateComment(ctx context.Context, tx *Tx, id int, update pa.CommentUpdate) (*pa.Comment, error) {

}

func deleteComment(ctx context.Context, tx *Tx, id int) error {

}
