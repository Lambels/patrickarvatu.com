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

	t.Run("Register Subscription Handler", func(t *testing.T) {
		s.RegisterSubscriptionsHandler(nil)
	})

	t.Run("Register Handlers", func(t *testing.T) {
		// Wont use ctx and hand for tests
		s.RegisterHandler(pa.EventTopicNewBlog, func(ctx context.Context, hand pa.SubscriptionService, event pa.Event) error {
			payload := new(pa.BlogPayload)
			if err := json.Unmarshal(event.Payload.([]byte), payload); err != nil {
				t.Fatal(err)
			}

			t.Logf("Payload: %v", payload.Blog.ID)
			return nil
		})
	})

	t.Run("Pushing Event", func(t *testing.T) {
		event := pa.Event{
			Topic: pa.EventTopicNewBlog,
			Payload: pa.BlogPayload{
				Blog: &pa.Blog{
					ID: 123,
				},
			},
		}

		if err := s.Push(context.Background(), event); err != nil {
			t.Fatal(err)
		}
	})
}
