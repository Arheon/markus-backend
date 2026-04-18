package message

import (
	"context"
	"time"

	"github.com/Arheon/markus-backend/internal/message/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateNewMessage(ctx context.Context, memberID string, roomID string, value string) (*string, error) {
	newUUID := uuid.New()
	err := gorm.G[entity.Message](r.db).Create(ctx, &entity.Message{
		ID:        newUUID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		RoomID:    roomID,
		MemberID:  memberID,
		Value:     value,
	})

	if err != nil {
		return nil, err
	}

	uuidString := newUUID.String()
	return &uuidString, nil
}

func (r *Repository) GetMessagesInRoom(ctx context.Context, roomID string) ([]*entity.Message, error) {
	return gorm.G[*entity.Message](r.db).Where("room_id = ?", roomID).Find(ctx)
}
