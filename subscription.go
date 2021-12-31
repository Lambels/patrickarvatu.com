package pa

import "context"

// Subscription represents a subscription handeled by the event handler on an event / topic.
type Subscription struct {
	// the pk of the subscription.
	ID int

	// the subscribed user.
	UserID int

	// topic to which the user is subscribed, ie: EventTopicNewBlog -> ./event.go.
	Topic string
}

// SubscriptionService represents a service which manages subscriptions in the system.
type SubscriptionService interface {
	// FindSubscriptionByID returns a subscription based on the id.
	// returns ENOTFOUND if the subscription doesent exist.
	FindSubscriptionByID(ctx context.Context, id int) (*Subscription, error)

	// FindSubscriptions returns a range of subscriptions and the length of the range. If filter
	// is specified FindSubscriptions will apply the filter to return set response.
	FindSubscriptions(ctx context.Context, filter SubscriptionFilter) ([]*Subscription, int, error)

	// CreateSubscription creates a subscription on topic. User is passed through context.
	CreateSubscription(ctx context.Context, topic string) (*Subscription, error)

	// DeleteSubscription permanently deletes a subscription. Returns EUNAUTHORIZED if the user owning the subscription
	// isnt the one calling and ENOTFOUND id the subscription doesent exist. User is passed through context.
	DeleteSubscription(ctx context.Context, id int) error
}

// SubscriptionFilter represents a filter used by FindSubscriptions to filter the response.
type SubscriptionFilter struct {
	// fields to filter on.
	ID     *int    `json:"id"`
	UserID *int    `json:"userID"`
	Topic  *string `json:"topic"`

	// restrictions on the result set, used for pagination and set limits.
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}
