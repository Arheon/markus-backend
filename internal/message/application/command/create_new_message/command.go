package createnewmessage

import (
	"context"

	"github.com/Arheon/markus-backend/internal/message/domain/event"
	"github.com/Arheon/markus-backend/internal/message/domain/repository"
	serverRepositry "github.com/Arheon/markus-backend/internal/server/domain/repository"
	"github.com/Arheon/markus-backend/pkg/outbox"
)

type Command struct {
	repo       repository.MessageRepository
	memberRepo serverRepositry.MemberRepository
	publisher  *outbox.Publisher
}

func NewCommand(
	memberRepo serverRepositry.MemberRepository,
	repo repository.MessageRepository,
	publisher *outbox.Publisher,
) *Command {
	return &Command{
		memberRepo: memberRepo,
		repo:       repo,
		publisher:  publisher,
	}
}

func (c *Command) Handle(ctx context.Context, userID string, serverID string, roomID string, value string) (*Result, error) {
	messageID, err := outbox.Publish(c.publisher, ctx, func(ec *outbox.EventCollector) (*string, error) {
		member, err := c.memberRepo.GetMemberByUserAndServer(ctx, serverID, userID)
		if err != nil {
			return nil, err
		}

		messageID, err := c.repo.CreateNewMessage(ctx, member.ID, roomID, value)
		if err != nil {
			return nil, err
		}

		ec.Add(&event.CreateMessageEvent{
			MemberID: userID,
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
