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

}

func (s *ProjectService) FindProjectByName(ctx context.Context, name string) (*pa.Project, error) {

}

func (s *ProjectService) FindProjects(ctx context.Context, filter pa.ProjectFilter) ([]*pa.Project, int, error) {

}

func (s *ProjectService) CreateOrUpdateProject(ctx context.Context, project *pa.Project) error {

}

func (s *ProjectService) DeleteProject(ctx context.Context, title string) error {

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

}

func createProject(ctx context.Context, tx *Tx, project *pa.Project) error {

}

func updateProject(ctx context.Context, tx *Tx, id int, project *pa.Project) error {

}

func deleteProject(ctx context.Context, tx *Tx, name string) error {

}

// topics: many 2 many interface functions ----------------------------------------------

func createNewTopicDescription(ctx context.Context, tx *Tx, content string) error {

}

func createNewTopicLink(ctx context.Context, topicDescID, projectID int) error {

}

func findTopicsByProjectID(ctx context.Context, tx *Tx, id int) []string {

}

func attachTopicsToProject(ctx context.Context, tx *Tx, project *pa.Project) error {

}
