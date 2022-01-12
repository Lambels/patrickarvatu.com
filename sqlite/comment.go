package sqlite

import (
	"context"
	"strings"

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

func findComments(ctx context.Context, tx *Tx, filter pa.CommentFilter) (_ []*pa.Comment, n int, err error) {
	// build where and args statement method.
	// not vulnerable to sql injection attack.
	where, args := []string{"1 = 1"}, []interface{}{}

	if v := filter.ID; v != nil {
		where = append(where, "id = ?")
		args = append(args, *v)
	}
	if v := filter.SubBlogID; v != nil {
		where = append(where, "sub_blog_id = ?")
		args = append(args, *v)
	}
	if v := filter.UserID; v != nil {
		where = append(where, "user_id = ?")
		args = append(args, *v)
	}

	rows, err := tx.QueryContext(ctx, `
		SELECT
			id,
			sub_blog_id,
			user_id,
			content,
			created_at,
			COUNT(*) OVER()
		FROM comments
		WHERE`+strings.Join(where, " AND ")+`
		ORDER BY id ASC
		`+FormatLimitOffset(filter.Limit, filter.Offset)+`
	`,
		args...,
	)

	if err != nil {
		return nil, n, err
	}

	// deserialize rows.
	comments := []*pa.Comment{}
	for rows.Next() {
		var comment *pa.Comment

		if err := rows.Scan(
			&comment.ID,
			&comment.SubBlogID,
			&comment.UserID,
			&comment.Content,
			(*NullTime)(&comment.CreatedAt),
			&n,
		); err != nil {
			return nil, 0, err
		}

		comments = append(comments, comment)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return comments, n, nil
}

func createComment(ctx context.Context, tx *Tx, comment *pa.Comment) error {
	comment.UserID = pa.UserIDFromContext(ctx)
	comment.CreatedAt = tx.now

	if err := comment.Validate(); err != nil {
		return err
	}

	result, err := tx.ExecContext(ctx, `
		INSERT INTO comments (
			sub_blog_id,
			user_id,
			content,
			created_at,
		)
		VALUES(?, ?, ?, ?)
	`,
		comment.SubBlogID,
		comment.UserID,
		comment.Content,
		(*NullTime)(&comment.CreatedAt),
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	// set id from database to comment obj.
	comment.ID = int(id)
	return nil
}

func updateComment(ctx context.Context, tx *Tx, id int, update pa.CommentUpdate) (*pa.Comment, error) {
	if !pa.IsAdminContext(ctx) {
		return nil, pa.Errorf(pa.EUNAUTHORIZED, "user isnt admin.")
	}

	comment, err := findCommentByID(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	if v := update.Content; v != nil {
		comment.Content = *v
	}

	if err := comment.Validate(); err != nil {
		return comment, err
	}

	if _, err := tx.ExecContext(ctx, `
		UPDATE comments,
		SET content = ?,
		WHERE id = ?
	`,
		comment.Content,
		id,
	); err != nil {
		return nil, err
	}

	return comment, nil
}

func deleteComment(ctx context.Context, tx *Tx, id int) error {
	comment, err := findCommentByID(ctx, tx, id)
	if err != nil {
		return err
	}

	if comment.UserID != pa.UserIDFromContext(ctx) {
		return pa.Errorf(pa.EUNAUTHORIZED, "user cant delete comment")
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM comments WHERE id = ?`, id); err != nil {
		return err
	}
	return nil
}

func attachUserToComment(ctx context.Context, tx *Tx, comment *pa.Comment) error {
	user, err := findUserByID(ctx, tx, comment.UserID)
	if err != nil {
		return err
	}

	comment.User = user
	return nil
}
