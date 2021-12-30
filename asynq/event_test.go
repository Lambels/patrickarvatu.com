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

	t.Run("Register Handlers", func(t *testing.T) {
		t.Log("Registering handler")
		s.RegisterHandler(pa.EventTopicNewBlog, func(ctx context.Context, event pa.Event) error {
			payload := new(pa.BlogPayload)
			if err := json.Unmarshal(event.Payload.([]byte), payload); err != nil {
				t.Fatal(err)
			}

			t.Logf("Payload: %v", payload.Blog.ID)
			return nil
		})
	})

	t.Run("Pushing Event", func(t *testing.T) {
		ctx := context.Background()
		userCtx := pa.NewContextWithUser(ctx, &pa.User{ID: 123})

		event := pa.Event{
			Topic: pa.EventTopicNewBlog,
			Payload: pa.BlogPayload{
				Blog: &pa.Blog{
					ID: 123,
				},
			},
		}

		t.Log("Pushing event")
		if err := s.Push(userCtx, event); err != nil {
			t.Fatal(err)
		}
	})
}
