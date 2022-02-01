package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	pa "github.com/Lambels/patrickarvatu.com"
	"github.com/go-chi/chi/v5"
)

// registerSubscriptionRoutes registers the subscription routes under r.
func (s *Server) registerSubscriptionRoutes(r chi.Router) {
	r.Get("/", s.handleGetSubscriptions)
	r.Get("/{topic}/{subscriptionID}", s.handleGetSubscription)

	r.Post("/", s.handleCreateSubscription)

	r.Delete("/{topic}/{subscriptionID}", s.handleDeleteSubscription)
}

// handleGetSubscriptions handels GET '/subscriptions/'
// retrieves subscriptions based on request body.
func (s *Server) handleGetSubscriptions(w http.ResponseWriter, r *http.Request) {
	var filter pa.SubscriptionFilter

	// get filter params from:
	switch w.Header().Get("Content-Type") {
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

	// overwrite filter with caller uID.
	temp := pa.UserIDFromContext(r.Context())
	filter.UserID = &temp

	resp := new(getSubscriptionsResponse)
	// we have no specific topic -> get both
	if *filter.Topic == "" {
		// subscription filter with new sub blog topic.
		v := pa.EventTopicNewSubBlog
		filter.Topic = &v
		blogSubscriptions, nBlogSubs, err := s.SubscriptionService.FindSubscriptions(r.Context(), filter)
		if err != nil {
			SendError(w, r, err)
			return
		}

		// subscription filter with new comment topic.
		v = pa.EventTopicNewComment
		filter.Topic = &v
		subBlogSubscriptions, nSubBlogSubs, err := s.SubscriptionService.FindSubscriptions(r.Context(), filter)
		if err != nil {
			SendError(w, r, err)
			return
		}

		resp.serializeIn(blogSubscriptions...)
		resp.serializeIn(subBlogSubscriptions...)
		resp.N = nBlogSubs + nSubBlogSubs
	} else { // we have a topic so we send it to find subscription to handle it.
		subscriptions, n, err := s.SubscriptionService.FindSubscriptions(r.Context(), filter)
		if err != nil {
			SendError(w, r, err)
			return
		}

		resp.serializeIn(subscriptions...)
		resp.N = n
	}

	// send response.
	SendJSON(w, resp)
}

// handleGetSubscription handels GET '/subscriptions/{topic}/{subscriptionID}'
// retrieves subscription with id: subscriptionID and topic: topic.
func (s *Server) handleGetSubscription(w http.ResponseWriter, r *http.Request) {
	topic := chi.URLParam(r, "topic")
	if topic == "" {
		SendError(w, r, pa.Errorf(pa.EINVALID, "empty topic"))
		return
	}

	subscriptionID, err := strconv.Atoi(chi.URLParam(r, "subscriptionID"))
	if err != nil {
		SendError(w, r, pa.Errorf(pa.EINVALID, "invalid id format"))
		return
	}

	// fetch subscription from database.
	subscription, err := s.SubscriptionService.FindSubscriptionByID(r.Context(), subscriptionID, topic)
	if err != nil {
		SendError(w, r, err)
		return
	}

	// send response.
	SendJSON(w, subscription)
}

// handleCreateSubscription handels POST '/subscriptions/'
// creates a subscription with the request body.
func (s *Server) handleCreateSubscription(w http.ResponseWriter, r *http.Request) {
	var subscription pa.Subscription

	// decode body.
	if err := json.NewDecoder(r.Body).Decode(&subscription); err != nil {
		SendError(w, r, err)
		return
	}

	// create subscription.
	if err := s.SubscriptionService.CreateSubscription(r.Context(), &subscription); err != nil {
		SendError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// handleDeleteSubscription handels DELETE '/subscriptions/{topic}/{subscriptionID}'
// deletes subscription with id: subscriptionID and topic: topic.
func (s *Server) handleDeleteSubscription(w http.ResponseWriter, r *http.Request) {
	topic := chi.URLParam(r, "topic")
	if topic == "" {
		SendError(w, r, pa.Errorf(pa.EINVALID, "empty topic"))
		return
	}

	subscriptionID, err := strconv.Atoi(chi.URLParam(r, "subscriptionID"))
	if err != nil {
		SendError(w, r, pa.Errorf(pa.EINVALID, "invalid id format"))
		return
	}

	// delete subscription
	if err := s.SubscriptionService.DeleteSubscription(r.Context(), subscriptionID, topic); err != nil {
		SendError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
