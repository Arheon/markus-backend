package repository

import (
	"context"

	domainEntity "github.com/Arheon/markus-backend/internal/server/domain/entity"
)

//mockery:generate: true
type ServerRepository interface {
	GetServerByID(ctx context.Context, serverID string) (*domainEntity.Server, error)
	GetAllServersByUserID(ctx context.Context, userID string) ([]*domainEntity.Server, error)
	CreateNewServer(ctx context.Context, userID string, serverName string) (*domainEntity.Server, error)
}
