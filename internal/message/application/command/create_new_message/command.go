package createnewmessage

import (
	"context"

	"github.com/Arheon/markus-backend/internal/message/domain/event"
	"github.com/Arheon/markus-backend/internal/message/domain/repository"
	"github.com/Arheon/markus-backend/pkg/outbox"
)

type Command struct {
	repo      repository.MessageRepository
	publisher *outbox.Publisher
}

func NewCommand(
	repo repository.MessageRepository,
	publisher *outbox.Publisher,
) *Command {
	return &Command{
		repo:      repo,
		publisher: publisher,
	}
}

func (c *Command) Handle(ctx context.Context, memberID string, roomID string, value string) (*Result, error) {
	messageID, err := outbox.Publish[string](c.publisher, ctx, func(ec *outbox.EventCollector) (*string, error) {
		messageID, err := c.repo.CreateNewMessage(ctx, memberID, roomID, value)
		if err != nil {
			return nil, err
		}

		ec.Add(&event.CreateMessageEvent{
			MemberID: memberID,
			RoomID:   roomID,
			Value:    value,
		})

		return messageID, nil
	})

	if err != nil {
		return nil, err
	}

	return &Result{MessageID: *messageID}, nil
}
