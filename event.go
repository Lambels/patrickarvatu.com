package pa

import "github.com/hibiken/asynq"

// event topics
const (
	EventTopicNewBlog = "blog:new"

	// Sub blogs are branched under blogs
	EventTopicNewSubBlog = "blog:sub_blog:new"

	// Comments are branched under sub blogs
	EventTopicNewComment = "blog:sub_blog:comment:new"
)

type Payload interface{}

type NewBlogPayload struct {
	Blog *Blog
}

type NewSubBlog struct {
	SubBlog *SubBlog
}

type EventService interface {
	Push(topic string, payload Payload) error

	RegisterHandler(topic string, handler asynq.HandlerFunc)
}

type NOPEventService struct{}

func (n *NOPEventService) Push(topic string, payload Payload) error

func (n *NOPEventService) RegisterHandler(topic string, handler asynq.HandlerFunc)

func NewNOPEventService() EventService { return &NOPEventService{} }
