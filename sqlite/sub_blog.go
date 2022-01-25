package sqlite

import (
	"context"
	"strings"

	pa "github.com/Lambels/patrickarvatu.com"
)

// check to see if *SubBlogService object implements set interface.
var _ pa.SubBlogService = (*SubBlogService)(nil)

// SubBlogService represents a service used to manage sub blogs.
type SubBlogService struct {
	db *DB
}

// NewSubBlogService returns a new instance of SubBlogService attached to db.
func NewSubBlogService(db *DB) *SubBlogService {
	return &SubBlogService{
		db: db,
	}
}

// FindSubBlogByID returns a sub blog based on id.
// returns ENOTFOUND if the sub blog doesent exist.
func (s *SubBlogService) FindSubBlogByID(ctx context.Context, id int) (*pa.SubBlog, error) {
	tx, err := s.db.BeginTX(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	subBlog, err := findSubBlogByID(ctx, tx, id)
	if err != nil {
		return nil, err

	} else if err := attachCommentsToSubBlog(ctx, tx, subBlog); err != nil {
		return nil, err
	}
	return subBlog, nil
}

// FindSubBlogs returns a range of sub blog based on filter.
func (s *SubBlogService) FindSubBlogs(ctx context.Context, filter pa.SubBlogFilter) ([]*pa.SubBlog, int, error) {
	tx, err := s.db.BeginTX(ctx, nil)
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback()

	subBlogs, n, err := findSubBlogs(ctx, tx, filter)
	if err != nil {
		return subBlogs, n, err
	}

	for _, subBlog := range subBlogs {
		if err := attachCommentsToSubBlog(ctx, tx, subBlog); err != nil {
			return subBlogs, n, err
		}
	}
	return subBlogs, n, nil
}

// CreateSubBlog creates a new sub blog.
// returns EUNAUTHORIZED if the user trying to create the sub blog isnt the admin user.
func (s *SubBlogService) CreateSubBlog(ctx context.Context, subBlog *pa.SubBlog) error {
	tx, err := s.db.BeginTX(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := createSubBlog(ctx, tx, subBlog); err != nil {
		return err
	} // nothing to attach as the only field to be populated is comments which cant exist before sub blog.
	return tx.Commit()
}

// UpdateSubBlog updates sub blog with id: id.
// returns EUNAUTHORIZED if the user trying to update the sub blog isnt the admin user.
// returns ENOTFOUND if the sub blog doesent exist.
func (s *SubBlogService) UpdateSubBlog(ctx context.Context, id int, update pa.SubBlogUpdate) (*pa.SubBlog, error) {
	tx, err := s.db.BeginTX(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	subBlog, err := updateSubBlog(ctx, tx, id, update)
	if err != nil {
		return nil, err
	} else if err := attachCommentsToSubBlog(ctx, tx, subBlog); err != nil {
		return nil, err
	}
	return subBlog, tx.Commit()
}

// DeleteSubBlog permanently deletes the sub blog specified by id.
// returns EUNAUTHORIZED if the user trying to delete the sub blog isnt the admin user.
// returns ENOTFOUND if the sub blog doesent exist.
func (s *SubBlogService) DeleteSubBlog(ctx context.Context, id int) error {
	tx, err := s.db.BeginTX(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := deleteSubBlog(ctx, tx, id); err != nil {
		return err
	}
	return tx.Commit()
}

func findSubBlogByID(ctx context.Context, tx *Tx, id int) (*pa.SubBlog, error) {
	filter := pa.SubBlogFilter{
		ID: &id,
	}
	subBlogs, _, err := findSubBlogs(ctx, tx, filter)
	if err != nil {
		return nil, err
	} else if len(subBlogs) == 0 {
		return nil, pa.Errorf(pa.ENOTFOUND, "sub blog not found.")
	}
	return subBlogs[0], nil
}

func findSubBlogs(ctx context.Context, tx *Tx, filter pa.SubBlogFilter) (_ []*pa.SubBlog, n int, err error) {
	// build where and args statement method.
	// not vulnerable to sql injection attack.
	where, args := []string{"1 = 1"}, []interface{}{}

	if v := filter.ID; v != nil {
		where = append(where, "id = ?")
		args = append(args, *v)
	}
	if v := filter.Title; v != nil {
		where = append(where, "title = ?")
		args = append(args, *v)
	}
	if v := filter.BlogID; v != nil {
		where = append(where, "blog_id = ?")
		args = append(args, *v)
	}

	rows, err := tx.QueryContext(ctx, `
		SELECT
			id,
			title,
			blog_id,
			content,
			created_at,
			updated_at,
			COUNT(*) OVER()
		FROM sub_blogs
		WHERE `+strings.Join(where, " AND ")+`
		ORDER BY id ASC
		`+FormatLimitOffset(filter.Limit, filter.Offset)+`
	`,
		args...,
	)

	if err != nil {
		return nil, n, err
	}
	defer rows.Close()

	// deserialize rows.
	subBlogs := []*pa.SubBlog{}
	for rows.Next() {
		var subBlog pa.SubBlog

		if err := rows.Scan(
			&subBlog.ID,
			&subBlog.Title,
			&subBlog.BlogID,
			&subBlog.Content,
			(*NullTime)(&subBlog.CreatedAt),
			(*NullTime)(&subBlog.UpdatedAt),
			&n,
		); err != nil {
			return nil, 0, err
		}

		if err := attachCommentsToSubBlog(ctx, tx, &subBlog); err != nil {
			return nil, 0, err
		}

		subBlogs = append(subBlogs, &subBlog)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return subBlogs, n, nil
}

func createSubBlog(ctx context.Context, tx *Tx, subBlog *pa.SubBlog) error {
	if !pa.IsAdminContext(ctx) {
		return pa.Errorf(pa.EUNAUTHORIZED, "user isnt admin.")
	}

	subBlog.CreatedAt = tx.now
	subBlog.UpdatedAt = subBlog.CreatedAt

	if err := subBlog.Validate(); err != nil {
		return err
	}

	result, err := tx.ExecContext(ctx, `
		INSERT INTO sub_blogs (
			blog_id,
			title,
			content,
			created_at,
			updated_at
		)
		VALUES(?, ?, ?, ?, ?)
	`,
		subBlog.BlogID,
		subBlog.Title,
		subBlog.Content,
		(*NullTime)(&subBlog.CreatedAt),
		(*NullTime)(&subBlog.UpdatedAt),
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	// set id from database to blog obj.
	subBlog.ID = int(id)
	return nil
}

func updateSubBlog(ctx context.Context, tx *Tx, id int, update pa.SubBlogUpdate) (*pa.SubBlog, error) {
	if !pa.IsAdminContext(ctx) {
		return nil, pa.Errorf(pa.EUNAUTHORIZED, "user isnt admin.")
	}

	subBlog, err := findSubBlogByID(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	if v := update.Content; v != nil {
		subBlog.Content = *v
	}
	if v := update.Title; v != nil {
		subBlog.Title = *v
	}

	if err := subBlog.Validate(); err != nil {
		return nil, err
	}

	subBlog.UpdatedAt = tx.now

	if _, err := tx.ExecContext(ctx, `
		UPDATE sub_blogs
		SET content		= ?,
			title 		= ?,
			updated_at 	= ?
		WHERE id = ?	
	`,
		subBlog.Content,
		subBlog.Title,
		(*NullTime)(&subBlog.UpdatedAt),
		id,
	); err != nil {
		return nil, err
	}

	return subBlog, nil
}

func deleteSubBlog(ctx context.Context, tx *Tx, id int) error {
	if !pa.IsAdminContext(ctx) {
		return pa.Errorf(pa.EUNAUTHORIZED, "user isnt admin.")
	}

	if _, err := findSubBlogByID(ctx, tx, id); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM sub_blogs WHERE id = ?`, id); err != nil {
		return err
	}
	return nil
}

func attachCommentsToSubBlog(ctx context.Context, tx *Tx, subBlog *pa.SubBlog) error {
	filter := pa.CommentFilter{
		SubBlogID: &subBlog.ID,
	}

	comments, _, err := findComments(ctx, tx, filter)
	if err != nil {
		return err
	} // we dont care if there are no comments.

	// append found comments under sub blog: subBlog.
	subBlog.Comments = append(subBlog.Comments, comments...)
	return nil
}
