package asynq_test

import (
	"context"
	"encoding/json"
	"testing"

	pa "github.com/Lambels/patrickarvatu.com"
	"github.com/Lambels/patrickarvatu.com/asynq"
)

func TestEventService(t *testing.T) {
	s := asynq.NewEventService("127.0.0.1:6379")
	defer s.Close()
	go s.Open()

	t.Run("Register Subscription Handler", func(_ *testing.T) {
		s.RegisterSubscriptionsHandler(nil)
	})

	t.Run("Register Handlers", func(t *testing.T) {
		// Wont use ctx and hand for tests
		s.RegisterHandler(pa.EventTopicNewSubBlog, func(_ context.Context, _ pa.SubscriptionService, event pa.Event) error {
			var payload pa.CommentPayload
			if err := json.Unmarshal(event.Payload.([]byte), &payload); err != nil {
				t.Fatal(err)
			}

			t.Logf("Payload: %v", payload.SubBlog.ID)
			return nil
		})
	})

	t.Run("Pushing Event", func(t *testing.T) {
		event := pa.Event{
			Topic: pa.EventTopicNewSubBlog,
			Payload: pa.SubBlogPayload{
				SubBlog: &pa.SubBlog{
					ID: 123,
				},
			},
		}

		if err := s.Push(context.Background(), event); err != nil {
			t.Fatal(err)
		}
	})
}
