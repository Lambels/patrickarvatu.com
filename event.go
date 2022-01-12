package pa

import (
	"context"
)

// Event topics.
const (
	// Sub blogs are branched under blogs.
	EventTopicNewSubBlog = "blog:sub_blog:new"

	// Comments are branched under sub blogs.
	EventTopicNewComment = "blog:sub_blog:comment:new"
)

// EventHandler represents a fucntion which is called on each event.
type EventHandler func(ctx context.Context, handler SubscriptionService, event Event) error

// Payload is an iterface to pe used when accepting event payloads, ie: BlogPayload -> ./event.go.
type Payload interface{}

// Event is passed to EventHandler.
type Event struct {
	// The topic of the event, ie: EventTopicNewSubBlog -> ./event.go.
	Topic string

	// The payload of the event, ie: BlogPayload -> ./event.go.
	Payload Payload
}

// SubBlogPayload represents the payload carried by a EventTopicNewSubBlog -> ./event.go.
type SubBlogPayload struct {
	SubBlog *SubBlog `json:"subBlog"`
	BlogID  int      `json:"blogID"`
	Blog    *Blog    `json:"blog"`
}

// CommentPayload represents the payload carried by a EventTopicNewComment -> ./event.go.
type CommentPayload struct {
	Comment   *Comment `json:"comment"`
	SubBlogID int      `json:"subBlogID"`
	SubBlog   *SubBlog `json:"subBlog"`
}

// EventService represents a service which manages auth in the system.
type EventService interface {
	// Push pushes event in the event queue.
	Push(ctx context.Context, event Event) error

	// RegisterHandler registers handler as the handler for topic.
	RegisterHandler(topic string, handler EventHandler)

	// RegisterSubscriptionsHandler registers the subscriptions manager.
	RegisterSubscriptionsHandler(hand SubscriptionService)
}

// NOPEventService is EventService which does nothing.
// Should only be used in tests.
type NOPEventService struct{}

func (n *NOPEventService) Push(ctx context.Context, event Event) error { panic("Not implemented") }

func (n *NOPEventService) RegisterHandler(topic string, handler EventHandler) {}

func (n *NOPEventService) RegisterSubscriptionsHandler(hand SubscriptionService) {}

func NewNOPEventService() EventService { return &NOPEventService{} }
