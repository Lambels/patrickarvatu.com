package pa

import "context"

// Project represents a github api repo response simplified.
// this is consumed by our frontend to avoid getting github rate-limited.
type Project struct {
	// pk in our system.
	ID int `json:"id"`

	// github api fields.
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Topics      []string `json:"topics"`
	HtmlURL     string   `json:"html_url"`
}

// TopicLink represents a link between a topic and a project.
type TopicLink struct {
	ProjectID int
	TopicID   int
}

// Topic represents a topic from github.
type Topic struct {
	ID      int
	Content string
}

// Validate performs basic validation on the project.
// returns EINVALID if any error is found.
func (p *Project) Validate() error {
	if p.Name == "" {
		return Errorf(EINVALID, "name is a required field.")
	}
	if p.HtmlURL == "" {
		return Errorf(EINVALID, "htmlUrl is a required field.")
	}
	return nil
}

// ProjectService represents a service which manages projects in the system.
type ProjectService interface {
	// FindProjectByID returns a project based on the id.
	// returns ENOTFOUND if the project doesent exist.
	FindProjectByID(ctx context.Context, id int) (*Project, error)

	// FindProjectByName returns a project based on the name.
	// returns ENOTFOUND if the project doesent exist.
	// (helper function to also return ENOTFOUND when filtering on name)
	FindProjectByName(ctx context.Context, name string) (*Project, error)

	// FindProjects returns a range of preojects and the length of the range. If filter
	// is specified FindProjects will apply the filter to return set response.
	FindProjects(ctx context.Context, filter ProjectFilter) ([]*Project, int, error)

	// CreateOrUpdateProject checks for existing id field on project or any duplicate name, if any
	// found the project field will be used to update the pointed to project.
	// returns EUNAUTHORIZED if not used by admin user or internally.
	CreateOrUpdateProject(ctx context.Context, project *Project) error

	// DeleteProject permanently deletes a project based on name.
	// returns ENOTFOUND if project doesent exist.
	// returns EUNAUTHORIZED if not used by admin user or internally.
	DeleteProject(ctx context.Context, name string) error
}

// ProjectFilter represents a filter used by FindProjects to filter the response.
type ProjectFilter struct {
	// fields to filter on.
	ID   *int    `json:"id"`
	Name *string `json:"name"`

	// restrictions on the result set, used for pagination and set limits.
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}
