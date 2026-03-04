package getuserinfo

import (
	"context"

	"github.com/Arheon/markus-backend/internal/user/domain/entity"
	"github.com/Arheon/markus-backend/internal/user/domain/repository"
)

type Query struct {
	ctx            context.Context
	userRepository repository.UserRepository
}

func NewQuery(ctx context.Context, userRepository repository.UserRepository) *Query {
	return &Query{
		userRepository: userRepository,
	}
}

func (q *Query) Handle(id string) (*entity.User, error) {
	return q.userRepository.GetUserByID(q.ctx, id)
}
