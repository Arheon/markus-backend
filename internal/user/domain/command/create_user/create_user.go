package createuser

import (
	"context"
	"errors"

	"github.com/Arheon/markus-backend/internal/user/domain/entity"
	"github.com/Arheon/markus-backend/internal/user/domain/repository"
)

var ErrCantCreateUser = errors.New("can't create user")

type Command struct {
	ctx      context.Context
	userRepo repository.UserRepository
}

func NewCommand(ctx context.Context, userRepo repository.UserRepository) *Command {
	return &Command{
		ctx,
		userRepo,
	}
}

func (c *Command) Handle(username string, password string) (*entity.User, error) {
	user, err := c.userRepo.CreateNewUserOrErrorIfExists(c.ctx, username, password)
	if err != nil {
		return nil, err
	}

	return user, nil
}
