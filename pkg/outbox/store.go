package outbox

import (
	"context"

	"github.com/google/uuid"
)

type Store interface {
	Save(ctx context.Context, msg Message) error
	FetchPending(ctx context.Context, limit int) ([]Message, error)
	MarkProcessed(ctx context.Context, id uuid.UUID) error
	MarkFailed(ctx context.Context, id uuid.UUID, err error) error
}
