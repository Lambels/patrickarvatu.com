package http

import (
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"strconv"

	pa "github.com/Lambels/patrickarvatu.com"
	"github.com/go-chi/chi/v5"
)

// registerBlogRoutes registers the blog routes under r.
func (s *Server) registerBlogRoutes(r chi.Router) {
	fs := http.FileServer(http.Dir(s.conf.FileStructure.BlogImagesDir))
	r.Handle("/images/", http.StripPrefix("/images", fs))

	r.Get("/", s.handleGetBlogs)
	r.Get("/{blogID}", s.handleGetBlog)
	r.Get("/{blogID}/sub-blogs", s.handleGetSubBlogs)

	r.Route("/", func(r chi.Router) {
		r.Use(s.adminAuthMiddleware)

		r.Post("/", s.handleCreateBlog)

		r.Patch("/{blogID}", s.handleUpdateBlog)

		r.Put("/{blogID}/image", s.handleAttachBlogImage)

		r.Delete("/{blogID}", s.handleDeleteBlog)
		r.Delete("/{blogID}/image", s.handleDeleteBlogImage)
	})
}

// handleGetBlogs handels GET '/blogs/'
// retrieves blogs based on request body.
func (s *Server) handleGetBlogs(w http.ResponseWriter, r *http.Request) {
	var filter pa.BlogFilter

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

	// fetch blogs from database.
	blogs, n, err := s.BlogService.FindBlogs(r.Context(), filter)
	if err != nil {
		SendError(w, r, err)
		return
	}

	SendJSON(w, getBlogsResponse{
		N:     n,
		Blogs: blogs,
	})
}

// handleGetBlog handels GET '/blogs/{blogID}'
// retrieves blog with id: blogID.
func (s *Server) handleGetBlog(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "blogID"))
	if err != nil {
		SendError(w, r, pa.Errorf(pa.EINVALID, "invalid id format"))
		return
	}

	// fetch blog from database.
	blog, err := s.BlogService.FindBlogByID(r.Context(), id)
	if err != nil {
		SendError(w, r, err)
		return
	}

	// send response.
	SendJSON(w, blog)
}

// handleCreateBlog handels POST '/blogs/'
// creates a blog with the request body.
func (s *Server) handleCreateBlog(w http.ResponseWriter, r *http.Request) {
	var blog pa.Blog

	// decode body.
	if err := json.NewDecoder(r.Body).Decode(&blog); err != nil {
		SendError(w, r, err)
		return
	}

	// create blog.
	if err := s.BlogService.CreateBlog(r.Context(), &blog); err != nil {
		SendError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// handleUpdateBlog handels PATCH '/blogs/'
// creates a blog with the request body.
func (s *Server) handleUpdateBlog(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "blogID"))
	if err != nil {
		SendError(w, r, pa.Errorf(pa.EINVALID, "invalid id format"))
		return
	}

	// decode body.
	var update pa.BlogUpdate
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		SendError(w, r, err)
		return
	}

	// update blog.
	blog, err := s.BlogService.UpdateBlog(r.Context(), id, update)
	if err != nil {
		SendError(w, r, err)
		return
	}

	// send response.
	SendJSON(w, blog)
}

// handleDeleteBlog handels DELETE '/blogs/{blogID}'
// deletes blog with id: blogID.
func (s *Server) handleDeleteBlog(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "blogID"))
	if err != nil {
		SendError(w, r, pa.Errorf(pa.EINVALID, "invalid id format"))
		return
	}

	// delete blog.
	if err := s.BlogService.DeleteBlog(r.Context(), id); err != nil {
		SendError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleAttachBlogImage(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "blogID"))
	if err != nil {
		SendError(w, r, pa.Errorf(pa.EINVALID, "invalid id format"))
		return
	}

	if _, err := s.BlogService.FindBlogByID(r.Context(), id); err != nil { // cant attach image to un unexisting blog.
		SendError(w, r, err)
		return
	} else if ext, err := mime.ExtensionsByType(r.Header.Get("Content-Type")); err != nil {
		SendError(w, r, err)
		return
	} else if len(ext) == 0 {
		SendError(w, r, pa.Errorf(pa.EINVALID, "invalid mime type"))
		return
	} else if err := s.BlogsFileSystem.CreateFile(r.Context(), "/"+fmt.Sprint(id)+ext[0], r.Body); err != nil {
		SendError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (s *Server) handleDeleteBlogImage(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "blogID")

	if ext, err := mime.ExtensionsByType(r.Header.Get("Content-Type")); err != nil {
		SendError(w, r, err)
		return
	} else if len(ext) == 0 {
		SendError(w, r, pa.Errorf(pa.EINVALID, "invalid mime type"))
		return
	} else if err := s.BlogsFileSystem.DeleteFile(r.Context(), "/"+string(id)+ext[0]); err != nil {
		SendError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
