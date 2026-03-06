package repository

import (
	"context"

	domainEntity "github.com/Arheon/markus-backend/internal/server/domain/entity"
)

type ServerRepository interface {
	GetAllServersByUserID(ctx context.Context, userID string) ([]*domainEntity.Server, error)
	CreateNewServer(ctx context.Context, userID string, serverName string) (*domainEntity.Server, error)
}
