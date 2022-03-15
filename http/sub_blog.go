package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	pa "github.com/Lambels/patrickarvatu.com"
	"github.com/go-chi/chi/v5"
)

// registerBlogRoutes registers the sub blog routes under r.
func (s *Server) registerSubBlogRoutes(r chi.Router) {
	r.Get("/", s.handleGetSubBlogs)
	r.Get("/{subBlogID}", s.handleGetSubBlog)
	r.Get("/{subBlogID}/comments", s.handleGetComments)

	r.Route("/", func(r chi.Router) {
		r.Use(s.adminAuthMiddleware)

		r.Post("/", s.handleCreateSubBlog)

		r.Patch("/{subBlogID}", s.handleUpdateSubBlog)

		r.Delete("/{subBlogID}", s.handleDeleteSubBlog)
	})
}

// handleGetSubBlogs handels GET '/sub-blogs/', '/blogs/{blogID}/sub-blogs'
// looks for blogID and over writes it with anything passed in the body.
func (s *Server) handleGetSubBlogs(w http.ResponseWriter, r *http.Request) {
	var filter pa.SubBlogFilter

	param := chi.URLParam(r, "blogID")
	if param != "" { // we have a blog id param
		id, err := strconv.Atoi(param)
		if err != nil {
			SendError(w, r, pa.Errorf(pa.EINVALID, "invalid id format"))
			return
		}

		filter.BlogID = &id
	}

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

	// fetch data from database.
	subBlogs, n, err := s.SubBlogService.FindSubBlogs(r.Context(), filter)
	if err != nil {
		SendError(w, r, err)
		return
	}

	SendJSON(w, getSubBlogsResponse{
		N:        n,
		SubBlogs: subBlogs,
	})
}

// handleGetSubBlog handels GET '/sub-blogs/{subBlogID}'
// looks for sub blog with id: subBlogID.
func (s *Server) handleGetSubBlog(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "subBlogID"))
	if err != nil {
		SendError(w, r, pa.Errorf(pa.EINVALID, "invalid id format"))
		return
	}

	// fetch sub blog from database.
	subBlog, err := s.SubBlogService.FindSubBlogByID(r.Context(), id)
	if err != nil {
		SendError(w, r, err)
		return
	}

	// send response.
	SendJSON(w, subBlog)
}

// handleCreateSubBlog handels POST '/sub-blogs/'
// creates a sub blog with the request body and pushes a pa.EventTopicNewSubBlog -> ./event.go.
func (s *Server) handleCreateSubBlog(w http.ResponseWriter, r *http.Request) {
	var subBlog pa.SubBlog

	switch header := r.Header.Get("Content-Type"); {
	case header == "application/json":
		// decode body.
		if err := json.NewDecoder(r.Body).Decode(&subBlog); err != nil {
			SendError(w, r, pa.Errorf(pa.EINVALID, "invalid JSON body"))
			return
		}

	case strings.HasPrefix(header, "multipart/form-data"):
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		blogID, err := strconv.Atoi(r.FormValue("blogID"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		subBlog.BlogID = blogID
		subBlog.Title = r.FormValue("title")

		f, _, err := r.FormFile("content")
		if err != nil {
			SendError(w, r, err)
			return
		}

		buf, err := ioutil.ReadAll(f)
		if err != nil {
			SendError(w, r, err)
			return
		}

		subBlog.Content = string(buf)

	default:
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	// create sub blog.
	if err := s.SubBlogService.CreateSubBlog(r.Context(), &subBlog); err != nil {
		SendError(w, r, err)
		return
	}

	// push event.
	if err := s.publishNewEvent(r.Context(), pa.Event{
		Topic: pa.EventTopicNewSubBlog,
		Payload: pa.SubBlogPayload{
			BlogID: subBlog.BlogID, // attach only blog id to payload for easy redirect.
		},
	}); err != nil {
		SendError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// handleUpdateSubBlog handels PATCH '/sub-blogs/{subBlogID}'
// updates a sub blog based on request body and subBlogID.
func (s *Server) handleUpdateSubBlog(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "subBlogID"))
	if err != nil {
		SendError(w, r, pa.Errorf(pa.EINVALID, "invalid id format"))
		return
	}

	// decode body.
	var update pa.SubBlogUpdate
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		SendError(w, r, pa.Errorf(pa.EINVALID, "invalid JSON body"))
		return
	}

	// update sub blog.
	subBlog, err := s.SubBlogService.UpdateSubBlog(r.Context(), id, update)
	if err != nil {
		SendError(w, r, err)
		return
	}

	// send response.
	SendJSON(w, subBlog)
}

// handleDeleteSubBlog handels DELETE '/sub-blogs/{subBlogID}'
// permanently deletes the sub blog pointed to by subBlogID.
func (s *Server) handleDeleteSubBlog(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "subBlogID"))
	if err != nil {
		SendError(w, r, pa.Errorf(pa.EINVALID, "invalid id format"))
		return
	}

	// delete sub blog.
	if err := s.SubBlogService.DeleteSubBlog(r.Context(), id); err != nil {
		SendError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
