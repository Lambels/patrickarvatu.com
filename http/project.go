package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	pa "github.com/Lambels/patrickarvatu.com"
	"github.com/go-chi/chi/v5"
)

// registerProjectRoutes registers the project routes under r.
func (s *Server) registerProjectRoutes(r chi.Router) {
	r.Get("/", s.handleGetProjects)
	r.Get("/{projectIDOrName}", s.handleGetProject)

	r.Route("/", func(r chi.Router) {
		r.Use(s.adminAuthMiddleware)

		r.Post("/", s.handleCreateOrUpdateProject)

		r.Delete("/{projectName}", s.handleDeleteProject)
	})
}

// handleGetProjects handels GET '/projects/'
// retrieves projects based on request body.
func (s *Server) handleGetProjects(w http.ResponseWriter, r *http.Request) {
	var filter pa.ProjectFilter

	// get filter params from:
	switch r.Header.Get("Content-Type") {
	case "application/json":
		if err := json.NewDecoder(r.Body).Decode(&filter); err != nil {
			SendError(w, r, pa.Errorf(pa.EINVALID, "invalid JSON body"))
			return
		}

	default:
		var offset int
		var err error
		if v := r.URL.Query().Get("offset"); v != "" {
			if offset, err = strconv.Atoi(v); err != nil {
				SendError(w, r, pa.Errorf(pa.EINVALID, "invalid offset format"))
				return
			}
		}

		filter.Offset = offset
		filter.Limit = 20
	}

	// fetch projects from database.
	projects, n, err := s.ProjectService.FindProjects(r.Context(), filter)
	if err != nil {
		SendError(w, r, err)
		return
	}

	SendJSON(w, getProjectsResponse{
		N:        n,
		Projects: projects,
	})
}

// handleGetProject handels GET '/projects/{projectIDOrName}'
// checks if projectIDOrName is an integer, if not returns project with name: projectIDOrName
// else returns project with id: projectIDOrName.
func (s *Server) handleGetProject(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "projectIDOrName"))
	if err != nil { // we have a project name.
		// fetch project from database.
		project, err := s.ProjectService.FindProjectByName(r.Context(), chi.URLParam(r, "projectIDOrName"))
		if err != nil {
			SendError(w, r, err)
			return
		}

		// send response.
		SendJSON(w, project)
		return
	}

	// we have a integer.
	// fetch project from database.
	project, err := s.ProjectService.FindProjectByID(r.Context(), id)
	if err != nil {
		SendError(w, r, err)
		return
	}

	// send response.
	SendJSON(w, project)
}

// handleCreateOrDeleteProject handels POST '/projects/'
// creates or updates a blog with the request body.
func (s *Server) handleCreateOrUpdateProject(w http.ResponseWriter, r *http.Request) {
	var project pa.Project

	// decode body.
	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
		SendError(w, r, err)
		return
	}

	// create or update project.
	if err := s.ProjectService.CreateOrUpdateProject(r.Context(), &project); err != nil {
		SendError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// handleDeleteProject handels DELETE '/projects/{projectName}'
// deletes project with name: projectName.
func (s *Server) handleDeleteProject(w http.ResponseWriter, r *http.Request) {
	// delete project.
	if err := s.ProjectService.DeleteProject(r.Context(), chi.URLParam(r, "projectName")); err != nil {
		SendError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
