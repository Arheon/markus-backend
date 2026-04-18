package repository

import (
	"context"

	"github.com/Arheon/markus-backend/internal/message/domain/entity"
)

type MessageRepository interface {
	CreateNewMessage(ctx context.Context, memberID string, roomID string, value string) (*string, error)
	GetMessagesInRoom(ctx context.Context, roomID string) ([]*entity.Message, error)
}
