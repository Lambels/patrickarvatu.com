package sqlite

import (
	"context"

	pa "github.com/Lambels/patrickarvatu.com"
)

// check to see if *BlogService object implements set interface.
var _ pa.BlogService = (*BlogService)(nil)

// BlogService represents a service used to manage blogs.
type BlogService struct {
	db *DB
}

// NewBlogService returns a new instance of BlogService attached to db.
func NewBlogService(db *DB) *BlogService {
	return &BlogService{
		db: db,
	}
}

// FindBlogByID returns a blog based on id.
// returns ENOTFOUND if the blog doesent exist.
func (s *BlogService) FindBlogByID(ctx context.Context, id int) (*pa.Blog, error) {
	tx, err := s.db.BeginTX(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	blog, err := findBlogByID(ctx, tx, id)
	if err != nil {
		return nil, err

	} else if err = attachSubBlogsToBlog(ctx, tx, blog); err != nil { // attach subBlogs to blog.
		return nil, err
	}

	for _, subBlog := range blog.SubBlogs {
		if err := attachCommentsToSubBlog(ctx, tx, subBlog); err != nil { // attach comments to sub blog.
			return nil, err
		}
	}
	return blog, nil
}

// FindBlogs returns a range of blog based on filter.
func (s *BlogService) FindBlogs(ctx context.Context, filter pa.BlogFilter) ([]*pa.Blog, int, error) {
	tx, err := s.db.BeginTX(ctx, nil)
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback()

	blogs, n, err := findBlogs(ctx, tx, filter)
	if err != nil {
		return blogs, n, err
	}

	// loop over each blog.
	for _, blog := range blogs {
		if err := attachSubBlogsToBlog(ctx, tx, blog); err != nil { // attach sub blogs to each blog.
			return blogs, n, err

		} else {
			for _, subBlog := range blog.SubBlogs {
				if err := attachCommentsToSubBlog(ctx, tx, subBlog); err != nil { // attach comments to the sub blogs of each blog.
					return blogs, n, err
				}
			}
		}
	}
	return blogs, n, nil
}

// CreateBlog creates a new blog.
// returns EUNAUTHORIZED if the user trying to create the blog isnt the admin user.
func (s *BlogService) CreateBlog(ctx context.Context, blog *pa.Blog) error {
	tx, err := s.db.BeginTX(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// create blog and dont attach anything as the blog doesent have anything to populate on creation.
	if err := createBlog(ctx, tx, blog); err != nil {
		return err
	}
	return tx.Commit()
}

// UpdateBlog updates blog with id: id.
// returns EUNAUTHORIZED if the user trying to update the blog isnt the admin user.
// returns ENOTFOUND if the blog doesent exist.
func (s *BlogService) UpdateBlog(ctx context.Context, id int, update pa.BlogUpdate) (*pa.Blog, error) {
	tx, err := s.db.BeginTX(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// update and attach sub blogs and comments.
	blog, err := updateBlog(ctx, tx, id, update)
	if err != nil {
		return nil, err
	} else if err := attachSubBlogsToBlog(ctx, tx, blog); err != nil {
		return nil, err
	}

	for _, subBlog := range blog.SubBlogs {
		if err := attachCommentsToSubBlog(ctx, tx, subBlog); err != nil {
			return nil, err
		}
	}

	return blog, tx.Commit()
}

// DeleteBlog permanently deletes the blog specified by id.
// returns EUNAUTHORIZED if the user trying to delete the blog isnt the admin user.
// returns ENOTFOUND if the blog doesent exist.
func (s *BlogService) DeleteBlog(ctx context.Context, id int) error {
	tx, err := s.db.BeginTX(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := deleteBlog(ctx, tx, id); err != nil {
		return err
	}
	return tx.Commit()
}

func findBlogByID(ctx context.Context, tx *Tx, id int) (*pa.Blog, error) {

}

func findBlogs(ctx context.Context, tx *Tx, filter pa.BlogFilter) ([]*pa.Blog, int, error) {

}

func createBlog(ctx context.Context, tx *Tx, blog *pa.Blog) error {

}

func updateBlog(ctx context.Context, tx *Tx, id int, update pa.BlogUpdate) (*pa.Blog, error) {

}

func deleteBlog(ctx context.Context, tx *Tx, id int) error {

}

func attachSubBlogsToBlog(ctx context.Context, tx *Tx, blog *pa.Blog) error {

}
