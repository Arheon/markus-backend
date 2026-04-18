package repository

import (
	"context"

	"github.com/Arheon/markus-backend/internal/server/domain/entity"
)

type RoomRepository interface {
	CreateRoom(ctx context.Context, roomName string, serverID string, categoryID *string) (*string, error)
	GetServerRooms(ctx context.Context, serverID string) ([]entity.Room, error)
	GetRoomByID(ctx context.Context, roomID string) (*entity.Room, error)
}
