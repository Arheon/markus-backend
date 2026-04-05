package createuser

import (
	"context"
	"errors"

	"github.com/Arheon/markus-backend/internal/user/domain/repository"
	"github.com/Arheon/markus-backend/internal/user/domain/service/auth"

	globalHelpers "github.com/Arheon/markus-backend/internal/shared/infrastructure/helpers"
	"github.com/microcosm-cc/bluemonday"
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

func (c *Command) Handle(username string, password string) (string, error) {
	var sanitiser = bluemonday.StrictPolicy()
	username = sanitiser.Sanitize(username)
	password = sanitiser.Sanitize(password)

	if err := auth.ValidateUsername(username); err != nil {
		return "", err
	}

	if err := auth.ValidatePasswordStrangth(password); err != nil {
		return "", err
	}

	hashPassword, err := globalHelpers.HashPassword(password)
	user, err := c.userRepo.CreateNewUserOrErrorIfExists(c.ctx, username, hashPassword)
	if err != nil {
		return "", err
	}

	return user.ID, nil
}
