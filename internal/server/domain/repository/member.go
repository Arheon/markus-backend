package repository

import (
	"context"

	"github.com/Arheon/markus-backend/internal/server/domain/entity"
)

type MemberRepository interface {
	GetMemberByUserAndServer(ctx context.Context, serverID string, userID string) (entity.Member, error)
}
