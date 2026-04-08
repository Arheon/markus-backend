package outbox

import (
	"context"

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

func (p *Publisher) Publish(ctx context.Context, msg Message) error {
	msg.ID = p.uuid.NewUUID()
	msg.CreatedAt = p.time.Now().UTC()

	return p.store.Save(ctx, msg)
}
