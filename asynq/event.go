package asynq

import (
	"context"
	"encoding/json"

	pa "github.com/Lambels/patrickarvatu.com"
	"github.com/hibiken/asynq"
)

var _ pa.EventService = (*EventService)(nil)

type EventService struct {
	worker *asynq.Server
	mux    *asynq.ServeMux
	client *asynq.Client
}

func NewEventService(redisDSN string) *EventService {
	s := &EventService{}

	mux := asynq.NewServeMux()

	s.mux = mux

	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     redisDSN,
			Password: "",
		},
		asynq.Config{
			Concurrency: 10,
		},
	)

	s.worker = srv

	client := asynq.NewClient(
		asynq.RedisClientOpt{
			Addr:     redisDSN,
			Password: "",
		},
	)

	s.client = client

	return s
}

func (e *EventService) Open() error {
	return e.worker.Run(e.mux)
}

func (e *EventService) Close() error {
	e.worker.Shutdown()
	return e.client.Close()
}

func (e *EventService) Push(ctx context.Context, event pa.Event) error {
	if jsonPayload, err := json.Marshal(event.Payload); err == nil {
		_, err := e.client.EnqueueContext(ctx, asynq.NewTask(event.Topic, jsonPayload))
		return err
	} else {
		return err
	}
}

// Register all handlers before opening the server
func (e *EventService) RegisterHandler(topic string, handler pa.EventHandler) {
	e.mux.HandleFunc(topic, func(ctx context.Context, t *asynq.Task) error {
		return handler(
			ctx,
			pa.Event{
				Topic:   topic,
				Payload: t.Payload(),
			},
		)
	})
}
