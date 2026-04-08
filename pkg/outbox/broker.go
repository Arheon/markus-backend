package outbox

import "context"

type Broker interface {
	Publish(ctx context.Context, msg Message) error
}
