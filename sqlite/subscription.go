package sqlite

import (
	"context"
	"fmt"

	pa "github.com/Lambels/patrickarvatu.com"
)

// check to see if *SubscriptionService object implements set interface.
var _ pa.SubscriptionService = (*SubscriptionService)(nil)

// SubscriptionService represents a serivce used to manage subscriptions.
type SubscriptionService struct {
	db *DB
}

// NewSubscriptionService returns a new instance of SubscriptionService attached to db.
func NewSubscriptionService(db *DB) *SubscriptionService {
	return &SubscriptionService{
		db: db,
	}
}

// findXXXSubscription -------------------------------------------------------------

// FindSubscriptionByID returns a subscription based on the id and topic.
// returns ENOTFOUND if the subscription doesent exist.
func (s *SubscriptionService) FindSubscriptionByID(ctx context.Context, id int, topic string) (*pa.Subscription, error) {
	tx, err := s.db.BeginTX(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	return findSubscriptionByID(ctx, tx, id, topic)
}

// FindSubscriptions returns a range of subscriptions based on filter. Only returns subscriptions
// which the user owns.
func (s *SubscriptionService) FindSubscriptions(ctx context.Context, filter pa.SubscriptionFilter) ([]*pa.Subscription, int, error) {
	tx, err := s.db.BeginTX(ctx, nil)
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback()

	return findSubscriptions(ctx, tx, filter)
}

// findSubscriptionByID is a helper function creating a filter with only the id and topic
// to pass further down the line to findSubscriptions.
func findSubscriptionByID(ctx context.Context, tx *Tx, id int, topic string) (sub *pa.Subscription, err error) {
	filter := pa.SubscriptionFilter{
		ID:    &id,
		Topic: &topic,
	}

	subs, _, err := findSubscriptions(ctx, tx, filter)
	if err != nil {
		return nil, err
	} else if len(subs) == 0 {
		return nil, &pa.Error{
			Code:    pa.ENOTFOUND,
			Message: "subscription not found.",
		}
	}

	return subs[0], nil
}

// findSubscriptions runs the appropiate find function mapped to the topic.
func findSubscriptions(ctx context.Context, tx *Tx, filter pa.SubscriptionFilter) (subs []*pa.Subscription, n int, err error) {
	switch *filter.Topic {
	case pa.EventTopicNewSubBlog:
		// users subscribed to a blog recieve a notification when a sub blog is created
		// hence we search for blog subscriptions when looking for users subscribed to
		// pa.EventTopicNewSubBlog.
		subs, n, err = findBlogSubscriptions(ctx, tx, filter)
		if err != nil {
			return nil, 0, err
		}

	case pa.EventTopicNewComment:
		// users subscribed to a sub blog recieve a notification when a new comment is added on
		// the sub blog hence we look for sub blog subscriptions when looking for users subscribed to
		// pa.EventTopicNewComment.
		subs, n, err = findSubBlogSubscriptions(ctx, tx, filter)
		if err != nil {
			return nil, 0, err
		}

	default:
		fmt.Errorf("findSubscriptions: unidentified topic %v", *filter.Topic)
	}

	return subs, n, nil
}

// findBlogSubscriptions finds blog subscriptions pointed to by the filter.
func findBlogSubscriptions(ctx context.Context, tx *Tx, filter pa.SubscriptionFilter) (_ []*pa.Subscription, n int, err error) {

}

// findSubBlogSubscriptions finds sub blog subscriptions pointed to by the filter.
func findSubBlogSubscriptions(ctx context.Context, tx *Tx, filter pa.SubscriptionFilter) (_ []*pa.Subscription, n int, err error) {

}

// createXXXSubscription -------------------------------------------------------------

// CreateSubscription creates a new subscription and links it to the user.
func (s *SubscriptionService) CreateSubscription(ctx context.Context, subscription *pa.Subscription) error {
	tx, err := s.db.BeginTX(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := createSubscription(ctx, tx, subscription); err != nil {
		return err
	}
	return tx.Commit()
}

// ceateSubscription runs the appropiate create function mapped to the topic.
func createSubscription(ctx context.Context, tx *Tx, subscription *pa.Subscription) (err error) {
	switch subscription.Topic {
	case pa.EventTopicNewSubBlog:
		err = createBlogSubscription(ctx, tx, subscription)

	case pa.EventTopicNewComment:
		err = createSubBlogSubscription(ctx, tx, subscription)

	default:
		err = fmt.Errorf("findSubscriptions: unidentified topic %v", subscription.Topic)
	}

	return err
}

// createBlogSubscription creates a blog subscription.
func createBlogSubscription(ctx context.Context, tx *Tx, sub *pa.Subscription) error {

}

// createSubBlogSubscription creates a sub blog subscription.
func createSubBlogSubscription(ctx context.Context, tx *Tx, sub *pa.Subscription) error {

}

// deleteXXXSubscription -------------------------------------------------------------

// DeleteSubscription permanently deletes the subscription only if the owner of the subscription
// is the user himself.
func (s *SubscriptionService) DeleteSubscription(ctx context.Context, id int, topic string) error {
	tx, err := s.db.BeginTX(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	filter := pa.SubscriptionFilter{
		ID:    &id,
		Topic: &topic,
	}

	if err := deleteSubscription(ctx, tx, filter); err != nil {
		return err
	}
	return tx.Commit()
}

// deleteSubscrption runs the appropiate delete function mapped to the topic.
func deleteSubscription(ctx context.Context, tx *Tx, filter pa.SubscriptionFilter) (err error) {
	switch *filter.Topic {
	case pa.EventTopicNewSubBlog:
		err = deleteBlogSubscription(ctx, tx, filter)

	case pa.EventTopicNewComment:
		err = deleteSubBlogSubscription(ctx, tx, filter)

	default:
		err = fmt.Errorf("findSubscriptions: unidentified topic %v", *filter.Topic)
	}

	return err
}

// deleteBlogSubscription deletes the BlogSubscription pointed to by the filter.
func deleteBlogSubscription(ctx context.Context, tx *Tx, filter pa.SubscriptionFilter) error {

}

// deleteSubBlogSubscription deletes the SubBlogSubscription pointed to by the filter.
func deleteSubBlogSubscription(ctx context.Context, tx *Tx, filter pa.SubscriptionFilter) error {

}
