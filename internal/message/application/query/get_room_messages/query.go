package getroommessages

import (
	"context"

	"github.com/Arheon/markus-backend/internal/message/domain/repository"
)

type Query struct {
	repo repository.MessageRepository
}

func NewQuery(
	repo repository.MessageRepository,
) *Query {
	return &Query{
		repo: repo,
	}
}

func (q *Query) Handle(ctx context.Context, roomID string) (*Result, error) {
	messages, err := q.repo.GetMessagesInRoom(ctx, roomID)
	if err != nil {
		return nil, err
	}

	var messageResults []ResultMessage
	for _, message := range messages {
		messageResult := ResultMessage{
			MemberID: message.MemberID,
			Value:    message.Value,
		}

		messageResults = append(messageResults, messageResult)
	}

	return &Result{RoomID: roomID, Messages: messageResults}, nil
}
