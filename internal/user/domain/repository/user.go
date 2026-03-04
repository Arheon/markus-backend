package repository

import (
	"context"

	"github.com/Arheon/markus-backend/internal/user/domain/entity"
)

type UserRepository interface {
	CreateNewUserOrErrorIfExists(context context.Context, username string, password string) error
	GetUserByID(context context.Context, userID string) (*entity.User, error)
}
