package user

import (
	"context"
	"errors"

	sharedEntity "github.com/Arheon/markus-backend/internal/shared/domain/entity"
	domainEntity "github.com/Arheon/markus-backend/internal/user/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrUserIsExists = errors.New("user is exists")

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateNewUserOrErrorIfExists(ctx context.Context, username string, password string) (*domainEntity.User, error) {
	count, err := gorm.G[domainEntity.User](r.db).Where("username = ?", username).Count(ctx, "username")
	if err != nil {
		return nil, err
	}

	if count > 0 {
		return nil, ErrUserIsExists
	}

	newUuid := uuid.New().String()
	newUser := &domainEntity.User{
		User: sharedEntity.User{
			ID:       newUuid,
			Username: username,
			Password: password,
		},
	}

	err = gorm.G[domainEntity.User](r.db).Create(ctx, newUser)
	if err != nil {
		return nil, err
	}

	return newUser, err
}

func (r *Repository) GetUserByID(ctx context.Context, userID string) (*domainEntity.User, error) {
	return gorm.G[*domainEntity.User](r.db).Where("id = ?", userID).First(ctx)
}
