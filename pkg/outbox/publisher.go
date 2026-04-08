package outbox

import (
	"context"
	"encoding/json"

	"github.com/Arheon/markus-backend/pkg/outbox/provider/time"
	"github.com/Arheon/markus-backend/pkg/outbox/provider/uuid"
)

type Publisher struct {
	store Store
	time  time.Provider
	uuid  uuid.Provider
}

type PublisherOption func(*Publisher)

func WithTime(time time.Provider) PublisherOption {
	return func(p *Publisher) {
		p.time = time
	}
}

func WithUUID(uuid uuid.Provider) PublisherOption {
	return func(p *Publisher) {
		p.uuid = uuid
	}
}

func NewPublisher(store Store, opts ...PublisherOption) *Publisher {
	p := &Publisher{
		store: store,
		time:  time.NewTimeProvider(),
		uuid:  uuid.NewUUIDProvider(),
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

func Publish[T any](p *Publisher, ctx context.Context, publisher func(c *EventCollector) (T, error)) (T, error) {
	collector := &EventCollector{}

	var result T
	result, err := publisher(collector)
	if err != nil {
		return result, err
	}

	for _, e := range collector.events {
		payload, err := json.Marshal(e)
		if err != nil {
			return result, err
		}

		msg := NewMessage(payload, WithType(e.EventName()))
		if err := p.store.Save(ctx, *msg); err != nil {
			return result, err
		}
	}
	return result, nil
}
