package sqlite

import (
	"context"

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

func findSubBlogs(ctx context.Context, tx *Tx, filter pa.SubBlogFilter) ([]*pa.SubBlog, int, error) {

}

func createSubBlog(ctx context.Context, tx *Tx, subBlog *pa.SubBlog) error {

}

func updateSubBlog(ctx context.Context, tx *Tx, id int, update pa.SubBlogUpdate) (*pa.SubBlog, error) {

}

func deleteSubBlog(ctx context.Context, tx *Tx, id int) error {

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
