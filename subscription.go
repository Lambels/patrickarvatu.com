package pa

import "context"

// Subscription represents a subscription handeled by the event handler on an event / topic.
type Subscription struct {
	// the pk of the subscription.
	ID int

	// the subscribed user.
	UserID int

	// topic to which the user is subscribed, ie: EventTopicNewBlog -> ./event.go.
	// each topic will map to a table in the database, logic handeled in service.
	Topic string

	// Payload is used to provide unique information about each Topic. The type of payload can be identified
	// by the topic, for example a Topic of type EventTopicNewBlog will come with a payload of type
	// BlogPayload -> ./event.go
	Payload Payload
}

// SubscriptionService represents a service which manages subscriptions in the system.
type SubscriptionService interface {
	// FindSubscriptionByID returns a subscription based on the id and topic.
	// returns ENOTFOUND if the subscription doesent exist.
	FindSubscriptionByID(ctx context.Context, id int, topic string) (*Subscription, error)

	// FindSubscriptions returns a range of subscriptions and the length of the range. If filter
	// is specified FindSubscriptions will apply the filter to return set response.
	FindSubscriptions(ctx context.Context, filter SubscriptionFilter) ([]*Subscription, int, error)

	// CreateSubscription creates a subscription on topic. User is passed through context.
	CreateSubscription(ctx context.Context, subscription *Subscription) error

	// DeleteSubscription permanently deletes a subscription. Returns EUNAUTHORIZED if the user owning the subscription
	// isnt the one calling and ENOTFOUND id the subscription doesent exist. User is passed through context.
	DeleteSubscription(ctx context.Context, id int, topic string) error
}

// SubscriptionFilter represents a filter used by FindSubscriptions to filter the response.
type SubscriptionFilter struct {
	// fields to filter on.
	ID      *int    `json:"id"`
	UserID  *int    `json:"userID"`
	Topic   *string `json:"topic"`
	Payload Payload `json:"payload"`

	// restrictions on the result set, used for pagination and set limits.
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}
