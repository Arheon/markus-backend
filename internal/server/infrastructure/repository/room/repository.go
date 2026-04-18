package room

import (
	"context"

	"github.com/Arheon/markus-backend/internal/server/domain/entity"
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

func (r *Repository) GetRoomByID(ctx context.Context, roomID string) (*entity.Room, error) {
	return gorm.G[*entity.Room](r.db).Preload("RoomCategory", nil).Where("id = ?", roomID).First(ctx)
}

func (r *Repository) CreateRoom(ctx context.Context, roomName string, serverID string, categoryID *string) (*string, error) {
	newRoomUUID := uuid.NewString()

	if err := gorm.G[entity.Room](r.db).Create(ctx, &entity.Room{
		ID:             newRoomUUID,
		RoomName:       roomName,
		ServerID:       serverID,
		RoomCategoryID: categoryID,
	}); err != nil {
		return nil, err
	}

	return &newRoomUUID, nil
}

func (r *Repository) GetServerRooms(ctx context.Context, serverID string) ([]entity.Room, error) {
	return gorm.G[entity.Room](r.db).Preload("RoomCategory", nil).Where("server_id = ?", serverID).Find(ctx)
}
