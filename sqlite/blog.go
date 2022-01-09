package sqlite

import (
	"context"
	"strings"

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
	filter := pa.BlogFilter{
		ID: &id,
	}

	blogs, _, err := findBlogs(ctx, tx, filter)
	if err != nil {
		return nil, err
	} else if len(blogs) == 0 {
		return nil, pa.Errorf(pa.ENOTFOUND, "blog not found.")
	}

	return blogs[0], nil
}

func findBlogs(ctx context.Context, tx *Tx, filter pa.BlogFilter) (_ []*pa.Blog, n int, err error) {
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

	rows, err := tx.QueryContext(ctx, `
		SELECT
			id,
			title,
			description,
			created_at,
			updated_at,
			COUNT(*) OVER()
		FROM blogs
		WHERE`+strings.Join(where, " AND ")+`
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
	blogs := []*pa.Blog{}
	for rows.Next() {
		var blog *pa.Blog

		if err := rows.Scan(
			&blog.ID,
			&blog.Title,
			&blog.Description,
			(*NullTime)(&blog.CreatedAt),
			(*NullTime)(&blog.UpdatedAt),
			&n,
		); err != nil {
			return nil, 0, err
		}

		blogs = append(blogs, blog)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return blogs, n, nil
}

func createBlog(ctx context.Context, tx *Tx, blog *pa.Blog) error {
	if !pa.IsAdminContext(ctx) {
		return pa.Errorf(pa.EUNAUTHORIZED, "user isnt admin.")
	}

	blog.CreatedAt = tx.now
	blog.UpdatedAt = blog.CreatedAt

	if err := blog.Validate(); err != nil {
		return err
	}

	result, err := tx.ExecContext(ctx, `
		INSERT INTO blogs (
			title,
			description,
			created_at,
			updated_at,
		)
		VALUES(?, ?, ?, ?)
	`,
		blog.Title,
		blog.Description,
		(*NullTime)(&blog.CreatedAt),
		(*NullTime)(&blog.UpdatedAt),
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	// set id from database to blog obj.
	blog.ID = int(id)
	return nil
}

func updateBlog(ctx context.Context, tx *Tx, id int, update pa.BlogUpdate) (*pa.Blog, error) {
	if !pa.IsAdminContext(ctx) {
		return nil, pa.Errorf(pa.EUNAUTHORIZED, "user isnt admin.")
	}

	blog, err := findBlogByID(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	if v := update.Title; v != nil {
		blog.Title = *v
	}
	if v := update.Description; v != nil {
		blog.Description = *v
	}

	if err := blog.Validate(); err != nil {
		return blog, err
	}

	blog.UpdatedAt = tx.now

	if _, err := tx.ExecContext(ctx, `
		UPDATE blogs,
		SET title 		= ?,
			description = ?,
			updated_at	= ?,
		WHERE id = ?
	`,
		blog.Title,
		blog.Description,
		blog.UpdatedAt,
		id,
	); err != nil {
		return nil, err
	}

	return blog, nil
}

func deleteBlog(ctx context.Context, tx *Tx, id int) error {
	if !pa.IsAdminContext(ctx) {
		return pa.Errorf(pa.EUNAUTHORIZED, "user isnt admin.")
	}

	if _, err := findBlogByID(ctx, tx, id); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM blogs WHERE id = ?`, id); err != nil {
		return err
	}
	return nil
}

func attachSubBlogsToBlog(ctx context.Context, tx *Tx, blog *pa.Blog) error {
	filter := pa.SubBlogFilter{
		BlogID: &blog.ID,
	}
	subBlogs, _, err := findSubBlogs(ctx, tx, filter)
	if err != nil {
		return err
	}

	blog.SubBlogs = append(blog.SubBlogs, subBlogs...)
	return nil
}
