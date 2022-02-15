package sqlite

import (
	"context"

	pa "github.com/Lambels/patrickarvatu.com"
)

// check to see if *ProjectService object implements set interface.
var _ pa.ProjectService = (*ProjectService)(nil)

// ProjectService represents a service used to manage projects.
type ProjectService struct {
	db *DB
}

// NewProjectService returns a new instance of ProjectService attached to db.
func NewProjectService(db *DB) *ProjectService {
	return &ProjectService{
		db: db,
	}
}

func (s *ProjectService) FindProjectByID(ctx context.Context, id int) (*pa.Project, error) {
	tx, err := s.db.BeginTX(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	proj, err := findProjectByID(ctx, tx, id)
	if err != nil {
		return nil, err

	} else if err = attachTopicsToProject(ctx, tx, proj); err != nil { // attach topics.

		return nil, err
	}

	return proj, nil
}

func (s *ProjectService) FindProjectByName(ctx context.Context, name string) (*pa.Project, error) {
	tx, err := s.db.BeginTX(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	proj, err := findProjectByName(ctx, tx, name)
	if err != nil {
		return nil, err

	} else if err = attachTopicsToProject(ctx, tx, proj); err != nil { // attach topics.

		return nil, err
	}

	return proj, nil
}

func (s *ProjectService) FindProjects(ctx context.Context, filter pa.ProjectFilter) ([]*pa.Project, int, error) {
	tx, err := s.db.BeginTX(ctx, nil)
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback()

	projects, n, err := findProjects(ctx, tx, filter)
	if err != nil {
		return projects, n, err
	}

	// loop over each project.
	for _, project := range projects {
		if err := attachTopicsToProject(ctx, tx, project); err != nil { // attach topics.
			return projects, n, err
		}
	}

	return projects, n, err
}

func (s *ProjectService) CreateOrUpdateProject(ctx context.Context, project *pa.Project) error {
	tx, err := s.db.BeginTX(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	proj, err := findProjectByName(ctx, tx, project.Name)
	switch pa.ErrorCode(err) {
	case pa.ENOTFOUND: // (project doesent exist)
		if err := createProject(ctx, tx, project); err != nil {
			return err
		}

	case "": // nil error ErrorCode is an empty string. (project exists)
		if err := updateProject(ctx, tx, proj.ID, project); err != nil {
			return err
		}

	default: // error code pa.EINTERNAL, could be replaced by `case pa.EINTERNAL:` yet this code is more idiomatic.
		return err
	}

	return tx.Commit()
}

func (s *ProjectService) DeleteProject(ctx context.Context, name string) error {
	tx, err := s.db.BeginTX(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := deleteProject(ctx, tx, name); err != nil {
		return err
	}

	return tx.Commit()
}

// projects interface functions ---------------------------------------------------------

func findProjectByID(ctx context.Context, tx *Tx, id int) (*pa.Project, error) {
	filter := pa.ProjectFilter{
		ID: &id,
	}

	projects, _, err := findProjects(ctx, tx, filter)
	if err != nil {
		return nil, err
	} else if len(projects) == 0 {
		return nil, pa.Errorf(pa.ENOTFOUND, "project not found.")
	}

	return projects[0], nil
}

// findProjectByName is a helper function to interface with findProjects and returns ENOTFOUND
// if project doesent exist.
func findProjectByName(ctx context.Context, tx *Tx, name string) (*pa.Project, error) {
	filter := pa.ProjectFilter{
		Name: &name,
	}

	projects, _, err := findProjects(ctx, tx, filter)
	if err != nil {
		return nil, err
	} else if len(projects) == 0 {
		return nil, pa.Errorf(pa.ENOTFOUND, "project not found.")
	}

	return projects[0], nil
}

func findProjects(ctx context.Context, tx *Tx, filter pa.ProjectFilter) (_ []*pa.Project, n int, err error) {
	return nil, 0, nil
}

// TODO: handle topics
func createProject(ctx context.Context, tx *Tx, project *pa.Project) error {
	return nil
}

func updateProject(ctx context.Context, tx *Tx, id int, project *pa.Project) error {
	return nil
}

func deleteProject(ctx context.Context, tx *Tx, name string) error {
	return nil
}

// topics: many 2 many interface functions ----------------------------------------------
// to not be used directly.

func createNewTopic(ctx context.Context, tx *Tx, content string) error {
	return nil
}

func findTopicByContent(ctx context.Context, tx *Tx, content string) (*pa.Topic, error) {
	return nil, nil
}

func createNewTopicLink(ctx context.Context, tx *Tx, topicLink *pa.TopicLink) error {
	return nil
}

func findTopicsByProjectID(ctx context.Context, tx *Tx, id int) ([]*pa.Topic, error) {
	return nil, nil
}

func attachTopicsToProject(ctx context.Context, tx *Tx, project *pa.Project) error {
	return nil
}
