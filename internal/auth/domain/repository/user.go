package repository

import "github.com/Arheon/markus-backend/internal/auth/domain/entity"

type UserRepository interface {
	GetUserByUsername(username string) (*entity.User, error)
}
