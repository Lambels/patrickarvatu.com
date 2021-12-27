package pa

// event topics
const (
	EventTopicNewBlog    = "blog:new"
	EventTopicNewSubBlog = "blog:new:sub_blog"
	EventTopicNewComment = "blog:new:comment"
)

type EventHandler func(event Event) Error

type Payload interface{}

type Event struct {
	Topic string

	Payload []byte
}

type NewBlogPayload struct {
	Blog *Blog
}

type NewSubBlog struct {
	SubBlog *SubBlog
}

type EventService interface {
	Push(topic string, payload Payload)

	Subscribe(topic string, handler EventHandler)
}
