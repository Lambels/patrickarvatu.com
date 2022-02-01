package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	pa "github.com/Lambels/patrickarvatu.com"
	"github.com/go-chi/chi/v5"
)

// registerCommentRoutes registers the comment routes under r.
func (s *Server) registerCommentRoutes(r chi.Router) {
	r.Get("/", s.handleGetComments)
	r.Get("/{commentID}", s.handleGetComment)

	r.Post("/", s.handleCreateComment)

	r.Patch("/{commentID}", s.handleUpdateComment)

	r.Delete("/{commentID}", s.handleDeleteComment)
}

// handleGetSubComments handels GET '/comments/', '/sub-blogs/{subBlogID}/comments'
// looks for subBlogID and over writes it with anything passed in the body.
func (s *Server) handleGetComments(w http.ResponseWriter, r *http.Request) {
	var filter pa.CommentFilter

	param := chi.URLParam(r, "subBlogID")
	if param != "" { // we have a sub blog id param
		id, err := strconv.Atoi(param)
		if err != nil {
			SendError(w, r, pa.Errorf(pa.EINVALID, "invalid id format"))
			return
		}

		filter.SubBlogID = &id
	}

	// get filter params from:
	switch r.Header.Get("Content-Type") {
	case "application/json":
		if err := json.NewDecoder(r.Body).Decode(&filter); err != nil {
			SendError(w, r, err)
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
	comments, n, err := s.CommentService.FindComments(r.Context(), filter)
	if err != nil {
		SendError(w, r, err)
		return
	}

	SendJSON(w, getCommentsResponse{
		N:        n,
		Comments: comments,
	})
}

// handleGetComment handels GET '/comments/{commentID}'
// looks for comment with id: commentID.
func (s *Server) handleGetComment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "commentID"))
	if err != nil {
		SendError(w, r, pa.Errorf(pa.EINVALID, "invalid id format"))
		return
	}

	// fetch comment from database.
	comment, err := s.CommentService.FindCommentByID(r.Context(), id)
	if err != nil {
		SendError(w, r, err)
		return
	}

	// send response.
	SendJSON(w, comment)
}

// handleCreateComment handels POST '/comments/'
// creates a comment with the request body, pushes a pa.EventTopicNewComment -> ./event.go and
// creates a subscription on the sub blog on which the comment exists.
func (s *Server) handleCreateComment(w http.ResponseWriter, r *http.Request) {
	var comment pa.Comment
	// decode body.
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		SendError(w, r, err)
		return
	}

	// create comment.
	if err := s.CommentService.CreateComment(r.Context(), &comment); err != nil {
		SendError(w, r, err)
		return
	}

	// push event.
	if err := s.EventService.Push(r.Context(), pa.Event{
		Topic: pa.EventTopicNewComment,
		Payload: pa.CommentPayload{
			SubBlogID: comment.SubBlogID, // attach only sub blog id to payload for easy redirect.
		},
	}); err != nil {
		SendError(w, r, err)
		return
	}

	// create subscription.
	if err := s.SubscriptionService.CreateSubscription(r.Context(), &pa.Subscription{
		Topic: pa.EventTopicNewComment, // user id gets allocated by create subscription so we save an allocation, create subscription
		Payload: pa.CommentPayload{ // is only accesable by auth users.
			SubBlogID: comment.SubBlogID,
		},
	}); err != nil {
		SendError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// handleUpdateComment handels PATCH '/comments/{commentID}'
// updates a comment based on request body and commentID.
func (s *Server) handleUpdateComment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "commentID"))
	if err != nil {
		SendError(w, r, pa.Errorf(pa.EINVALID, "invalid id format"))
		return
	}

	// decode body.
	var update pa.CommentUpdate
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		SendError(w, r, pa.Errorf(pa.EINVALID, "invalid JSON body"))
		return
	}

	// update comment.
	comment, err := s.CommentService.UpdateComment(r.Context(), id, update)
	if err != nil {
		SendError(w, r, err)
		return
	}

	// send response.
	SendJSON(w, comment)
}

// TODO: revise code

// handleDeleteComment handels DELETE '/comments/{commentID}'
// permanently deletes the comment pointed to by commentID and deletes any existing subscription on the sub blog on which
// the comment lives if there are no more left.
func (s *Server) handleDeleteComment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "commentID"))
	if err != nil {
		SendError(w, r, pa.Errorf(pa.EINVALID, "invalid id format"))
		return
	}

	// fetch comment before delete.
	comment, err := s.CommentService.FindCommentByID(r.Context(), id)
	if err != nil {
		SendError(w, r, err)
		return
	}

	// delete comment.
	if err = s.CommentService.DeleteComment(r.Context(), id); err != nil {
		SendError(w, r, err)
		return
	}

	uID := pa.UserIDFromContext(r.Context())

	// fetch comments after delete.
	comments, _, err := s.CommentService.FindComments(r.Context(), pa.CommentFilter{
		UserID:    &uID,
		SubBlogID: &comment.SubBlogID,
	})
	if err != nil {
		SendError(w, r, err)
		return

	} else if len(comments) == 0 { // we have no more comments on sub blog.
		// delete subscription if exists.
		if subs, _, err := s.SubscriptionService.FindSubscriptions(r.Context(), pa.SubscriptionFilter{ // check for subscription on sub blog.
			UserID: &uID,
			Payload: pa.CommentPayload{
				SubBlogID: comment.SubBlogID,
			},
		}); err != nil {
			SendError(w, r, err)
			return
		} else if len(subs) == 1 { // 1 user can only have 1 subscription on 1 comment.
			// delete subscription.
			if err := s.SubscriptionService.DeleteSubscription(r.Context(), subs[0].ID, subs[0].Topic); err != nil {
				SendError(w, r, err)
				return
			}
		}
	}
	w.WriteHeader(http.StatusNoContent)
}
