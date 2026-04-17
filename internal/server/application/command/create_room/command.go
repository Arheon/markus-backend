package createroom

import (
	"context"

	"github.com/Arheon/markus-backend/internal/server/domain/event"
	"github.com/Arheon/markus-backend/internal/server/domain/repository"
	"github.com/Arheon/markus-backend/pkg/outbox"
)

type Command struct {
	repo      repository.RoomRepository
	publisher *outbox.Publisher
}

func NewCommand(repo repository.RoomRepository, publisher *outbox.Publisher) *Command {
	return &Command{
		repo:      repo,
		publisher: publisher,
	}
}

func (c *Command) Handle(ctx context.Context, roomName string, serverID string, roomCategory *string) (*string, error) {
	roomID, err := outbox.Publish(c.publisher, ctx, func(ec *outbox.EventCollector) (*string, error) {
		roomID, err := c.repo.CreateRoom(ctx, roomName, serverID, roomCategory)
		if err != nil {
			return nil, err
		}

		ec.Add(&event.CreateRoomEvent{
			ServerID: serverID,
			RoomID:   *roomID,
		})

		return roomID, nil
	})

	return roomID, err
}
