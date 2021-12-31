package pa

import (
	"context"
)

// event topics
const (
	EventTopicNewBlog = "blog:new"

	// Sub blogs are branched under blogs
	EventTopicNewSubBlog = "blog:sub_blog:new"

	// Comments are branched under sub blogs
	EventTopicNewComment = "blog:sub_blog:comment:new"
)

type EventHandler func(ctx context.Context, handler SubscriptionService, event Event) error

type Payload interface{}

type Event struct {
	Topic string

	Payload Payload
}

type BlogPayload struct {
	Blog *Blog
}

type SubBlogPayload struct {
	SubBlog *SubBlog
}

type CommentPayload struct {
	Comment *Comment
}

// EventService represents a service which manages auth in the system.
type EventService interface {
	Push(ctx context.Context, event Event) error

	RegisterHandler(topic string, handler EventHandler)

	RegisterSubscriptionsHandler(hand SubscriptionService)
}

type NOPEventService struct{}

func (n *NOPEventService) Push(ctx context.Context, event Event) error { panic("Not implemented") }

func (n *NOPEventService) RegisterHandler(topic string, handler EventHandler) {}

func (n *NOPEventService) RegisterSubscriptionsHandler(hand SubscriptionService) {}

func NewNOPEventService() EventService { return &NOPEventService{} }
