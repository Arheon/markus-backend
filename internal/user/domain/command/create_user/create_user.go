package createuser

import (
	"context"
	"errors"

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

func (c *Command) Handle(username string, password string) error {
	if err := c.userRepo.CreateNewUserOrErrorIfExists(c.ctx, username, password); err != nil {
		return errors.Join(err, ErrCantCreateUser)
	}

	return nil
}
