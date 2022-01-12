package sqlite

import (
	"context"
	"fmt"
	"strings"

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
			return subs, n, err
		}

	case pa.EventTopicNewComment:
		// users subscribed to a sub blog recieve a notification when a new comment is added on
		// the sub blog hence we look for sub blog subscriptions when looking for users subscribed to
		// pa.EventTopicNewComment.
		subs, n, err = findSubBlogSubscriptions(ctx, tx, filter)
		if err != nil {
			return subs, n, err
		}

	default:
		fmt.Errorf("findSubscriptions: unidentified topic %v", *filter.Topic)
	}

	return subs, n, nil
}

// findBlogSubscriptions finds blog subscriptions pointed to by the filter.
func findBlogSubscriptions(ctx context.Context, tx *Tx, filter pa.SubscriptionFilter) (_ []*pa.Subscription, n int, err error) {
	// build where and args statement method.
	// not vulnerable to sql injection attack.
	where, args := []string{"1 = 1"}, []interface{}{}

	if v := filter.ID; v != nil {
		where = append(where, "id = ?")
		args = append(args, *v)
	}
	if v := filter.UserID; v != nil {
		where = append(where, "user_id = ?")
		args = append(args, *v)
	}
	if v := filter.Payload; v != nil {
		where = append(where, "blog_id = ?")
		args = append(args, v.(pa.SubBlogPayload).BlogID)
	}

	rows, err := tx.QueryContext(ctx, `
		SELECT
			id,
			user_id,
			blog_id,
			COUNT(*) OVER()
		FROM blog_subscriptions
		WHERE`+strings.Join(where, " AND ")+`
		ORDER BY id ASC
		`+FormatLimitOffset(filter.Limit, filter.Offset)+`
	`,
		args...,
	)

	if err != nil {
		return nil, n, err
	}

	// deserialize rows.
	subscriptions := []*pa.Subscription{}
	for rows.Next() {
		var subscription *pa.Subscription
		var payload pa.SubBlogPayload

		// set topic similar to all subscriptions.
		subscription.Topic = pa.EventTopicNewSubBlog

		if err := rows.Scan(
			&subscription.ID,
			&subscription.UserID,
			&payload.BlogID,
			&n,
		); err != nil {
			return nil, 0, err
		}

		// attach blog to payload through blog id.
		if payload.Blog, err = findBlogByID(ctx, tx, payload.BlogID); err != nil {
			return nil, 0, err
		}

		// attach payload to subscription.
		subscription.Payload = payload

		subscriptions = append(subscriptions, subscription)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return subscriptions, n, nil
}

// findSubBlogSubscriptions finds sub blog subscriptions pointed to by the filter.
func findSubBlogSubscriptions(ctx context.Context, tx *Tx, filter pa.SubscriptionFilter) (_ []*pa.Subscription, n int, err error) {
	// build where and args statement method.
	// not vulnerable to sql injection attack.
	where, args := []string{"1 = 1"}, []interface{}{}

	if v := filter.ID; v != nil {
		where = append(where, "id = ?")
		args = append(args, *v)
	}
	if v := filter.UserID; v != nil {
		where = append(where, "user_id = ?")
		args = append(args, *v)
	}
	if v := filter.Payload; v != nil {
		where = append(where, "sub_blog_id = ?")
		args = append(args, v.(pa.CommentPayload).SubBlogID)
	}

	rows, err := tx.QueryContext(ctx, `
		SELECT
			id,
			user_id,
			sub_blog_id,
			COUNT(*) OVER()
		FROM sub_blog_subscriptions
		WHERE`+strings.Join(where, " AND ")+`
		ORDER BY id ASC
		`+FormatLimitOffset(filter.Limit, filter.Offset)+`
	`,
		args...,
	)

	if err != nil {
		return nil, n, err
	}

	// deserialize rows.
	subscriptions := []*pa.Subscription{}
	for rows.Next() {
		var subscription *pa.Subscription
		var payload pa.CommentPayload

		// set topic similar to all subscriptions.
		subscription.Topic = pa.EventTopicNewComment

		if err := rows.Scan(
			&subscription.ID,
			&subscription.UserID,
			&payload.SubBlogID,
			&n,
		); err != nil {
			return nil, 0, err
		}

		// attach blog to payload through blog id.
		if payload.SubBlog, err = findSubBlogByID(ctx, tx, payload.SubBlogID); err != nil {
			return nil, 0, err
		}

		// attach payload to subscription.
		subscription.Payload = payload

		subscriptions = append(subscriptions, subscription)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return subscriptions, n, nil

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
	sub.UserID = pa.UserIDFromContext(ctx)

	result, err := tx.ExecContext(ctx, `
		INSERT INTO blog_subscriptions (
			user_id,
			blog_id,
		)
		VALUES(?, ?)
	`,
		sub.UserID,
		sub.Payload.(pa.SubBlogPayload).BlogID,
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	// set id from database to subscription obj.
	sub.ID = int(id)
	return nil
}

// createSubBlogSubscription creates a sub blog subscription.
func createSubBlogSubscription(ctx context.Context, tx *Tx, sub *pa.Subscription) error {
	sub.UserID = pa.UserIDFromContext(ctx)

	result, err := tx.ExecContext(ctx, `
		INSERT INTO sub_blog_subscriptions (
			user_id,
			sub_blog_id,
		)
		VALUES(?, ?)
	`,
		sub.UserID,
		sub.Payload.(pa.CommentPayload).SubBlogID,
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	// set id from database to subscription obj.
	sub.ID = int(id)
	return nil
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
	sub, _, err := findBlogSubscriptions(ctx, tx, filter)
	if err != nil {
		return err

	} else if len(sub) == 0 {
		return pa.Errorf(pa.ENOTFOUND, "subscription does not exist.")
	}

	if pa.UserIDFromContext(ctx) != sub[0].UserID {
		return pa.Errorf(pa.EUNAUTHORIZED, "user is unathorized.")
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM blog_subscriptions WHERE id = ?`, filter.ID); err != nil {
		return err
	}
	return nil
}

// deleteSubBlogSubscription deletes the SubBlogSubscription pointed to by the filter.
func deleteSubBlogSubscription(ctx context.Context, tx *Tx, filter pa.SubscriptionFilter) error {
	sub, _, err := findSubBlogSubscriptions(ctx, tx, filter)
	if err != nil {
		return err

	} else if len(sub) == 0 {
		return pa.Errorf(pa.ENOTFOUND, "subscription does not exist.")
	}

	if pa.UserIDFromContext(ctx) != sub[0].UserID {
		return pa.Errorf(pa.EUNAUTHORIZED, "user is unathorized.")
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM sub_blog_subscriptions WHERE id = ?`, filter.ID); err != nil {
		return err
	}
	return nil
}

// TODO: revise code
