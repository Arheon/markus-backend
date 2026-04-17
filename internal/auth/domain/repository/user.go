package repository

import (
	"context"

	"github.com/Arheon/markus-backend/internal/auth/domain/entity"
	"github.com/google/uuid"
)

//mockery:generate: true
type UserRepository interface {
	GetUserByUsername(ctx context.Context, username string) (*entity.User, error)
	GetUserByID(ctx context.Context, ID uuid.UUID) (*entity.User, error)
}
