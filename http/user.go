package http

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"strconv"

	pa "github.com/Lambels/patrickarvatu.com"
	"github.com/go-chi/chi/v5"
)

// registerUserRoutes registers the user routes under r.
func (s *Server) registerUserRoutes(r chi.Router) {
	r.Get("/{userID}", s.handleGetUser)
	r.Get("/{userID}/profile", s.handleUserProfile)
	r.Patch("/{userID}/refresh-api-key", s.handleRefreshApiKey)
	r.Delete("/{userID}", s.handleDeleteUser)
}

// handleGetUser handels GET '/users/{userID}'.
// sends a getOtherUser or getMyUser response depending if the user is requesting himself or not.
func (s *Server) handleGetUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		SendError(w, r, pa.Errorf(pa.EINVALID, "invalid id format"))
		return
	}

	// user is trying to see his own details which is equal to GET `oauth/user/me`.
	if pa.UserIDFromContext(r.Context()) == id {
		s.handleMe(w, r) // usually bad practice but there is no possible error outcome from handleMe
		return
	}

	// fetch user from database.
	user, err := s.UserService.FindUserByID(r.Context(), id)
	if err != nil {
		SendError(w, r, err)
		return
	}

	// send limited response.
	SendJSON(w, getOtherUserResponse{
		Username:  user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.String(),
	})
}

// handleUserProfile handels GET '/users/{userID}/profile'.
// sends a getOtherUserProfile or getMyProfile response depending if the user is requesting his own
// profile or not.
func (s *Server) handleUserProfile(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		SendError(w, r, pa.Errorf(pa.EINVALID, "invalid id format"))
		return
	}

	// fetch user obj from database.
	user, err := s.UserService.FindUserByID(r.Context(), id)
	if err != nil {
		SendError(w, r, err)
		return
	}

	// fetch comments from database.
	comments, nComments, err := s.CommentService.FindComments(r.Context(), pa.CommentFilter{UserID: &id})
	if err != nil {
		SendError(w, r, err)
		return
	}

	// user is trying to see his own profile.
	if pa.UserIDFromContext(r.Context()) == id {
		// subscription filter with user id.
		subFilter := pa.SubscriptionFilter{
			UserID: &id,
		}

		// subscription filter with new sub blog topic.
		v := pa.EventTopicNewSubBlog
		subFilter.Topic = &v
		blogSubscriptions, nBlogSubs, err := s.SubscriptionService.FindSubscriptions(r.Context(), subFilter)
		if err != nil {
			SendError(w, r, err)
			return
		}

		// subscription filter with new comment topic.
		v = pa.EventTopicNewComment
		subFilter.Topic = &v
		subBlogSubscriptions, nSubBlogSubs, err := s.SubscriptionService.FindSubscriptions(r.Context(), subFilter)
		if err != nil {
			SendError(w, r, err)
			return
		}

		// build response.
		var response getMyProfileResponse
		response.User = user

		response.Comments = getCommentsResponse{
			N:        nComments,
			Comments: comments,
		}

		subsResponse := new(getSubscriptionsResponse)
		subsResponse.serializeIn(blogSubscriptions...)    // serialize in blog subscriptions.
		subsResponse.serializeIn(subBlogSubscriptions...) // serialize in sub blog subscriptions.
		subsResponse.N = nBlogSubs + nSubBlogSubs
		response.Subscriptions = *subsResponse

		SendJSON(w, response)
		return
	}

	// send limited response.
	SendJSON(w, getOtherUserProfileResponse{
		User: getOtherUserResponse{
			Username:  user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.String(),
		},
		Comments: getCommentsResponse{
			N:        nComments,
			Comments: comments,
		},
	})
}

// handleRefreshApiKey handels PATCH '/users/{userID}/refresh-api-key'.
// refreshes with a new api key and send a getMyUser response if the user tries to refresh his own api key.
func (s *Server) handleRefreshApiKey(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		SendError(w, r, pa.Errorf(pa.EINVALID, "invalid id format"))
		return
	}

	// generate a random api-key.
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		SendError(w, r, err)
		return
	}
	// encode rand bytes.
	apiKey := base64.StdEncoding.EncodeToString(buf)

	// update user with new api key.
	user, err := s.UserService.UpdateUser(r.Context(), id, pa.UserUpdate{ApiKey: &apiKey})
	if err != nil {
		SendError(w, r, err)
		return
	}

	// send response with updated user.
	SendJSON(w, user)
}

// handleDeleteUser handels DELETE '/users/{userID}'.
// permanently deletes the user pointed to by userID and clears the session.
func (s *Server) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		SendError(w, r, pa.Errorf(pa.EINVALID, "invalid id format"))
		return
	}

	// delete user from database and clear session.
	if err := s.UserService.DeleteUser(r.Context(), id); err != nil {
		SendError(w, r, err)
		return
	} else if err := s.setSession(w, pa.Session{}); err != nil {
		SendError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
