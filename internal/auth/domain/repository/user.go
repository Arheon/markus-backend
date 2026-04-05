package repository

import (
	"context"

	"github.com/Arheon/markus-backend/internal/auth/domain/entity"
)

//mockery:generate: true
type UserRepository interface {
	GetUserByUsername(ctx context.Context, username string) (*entity.User, error)
}
