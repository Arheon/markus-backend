package getuserservers

import (
	"context"

	domainEntity "github.com/Arheon/markus-backend/internal/server/domain/entity"
	"github.com/Arheon/markus-backend/internal/server/domain/repository"
)

type Query struct {
	ctx  context.Context
	repo repository.ServerRepository
}

func NewQuery(ctx context.Context, repo repository.ServerRepository) *Query {
	return &Query{
		ctx:  ctx,
		repo: repo,
	}
}

func (q *Query) Handle(userID string) ([]*domainEntity.Server, error) {
	return q.repo.GetAllServersByUserID(q.ctx, userID)
}
