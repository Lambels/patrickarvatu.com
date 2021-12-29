package asynq

import (
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
			Addr: redisDSN,
		},
		asynq.Config{
			Concurrency: 10,
		},
	)

	s.worker = srv

	client := asynq.NewClient(
		asynq.RedisClientOpt{
			Addr: redisDSN,
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

func (e *EventService) Push(topic string, payload pa.Payload) error {
	if jsonPayload, err := json.Marshal(payload); err == nil {
		_, err := e.client.Enqueue(asynq.NewTask(topic, jsonPayload))
		return err
	} else {
		return err
	}
}

// Register all handlers before opening the server
func (e *EventService) RegisterHandler(topic string, handler asynq.HandlerFunc) {
	e.mux.HandleFunc(topic, handler)
}
