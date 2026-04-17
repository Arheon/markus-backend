package user

import (
	"context"

	"github.com/Arheon/markus-backend/internal/auth/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	return gorm.G[*entity.User](r.db).Where("username = ?", username).First(ctx)
}

func (r *Repository) GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	return gorm.G[*entity.User](r.db).Where("id = ?", id.String()).First(ctx)
}
