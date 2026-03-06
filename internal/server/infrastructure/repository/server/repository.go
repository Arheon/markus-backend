package server

import (
	"context"
	"errors"

	domainEntity "github.com/Arheon/markus-backend/internal/server/domain/entity"
	sharedEntity "github.com/Arheon/markus-backend/internal/shared/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrUndefinedUser = errors.New("undefined user")

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetAllServersByUserID(ctx context.Context, userID string) ([]*domainEntity.Server, error) {
	quert := gorm.G[*domainEntity.Server](r.db).Preload("Members", nil).
		Joins(
			clause.JoinTarget{
				Type:  clause.InnerJoin,
				Table: "server_members",
			},
			func(db gorm.JoinBuilder, joinTable, curTable clause.Table) error {
				db.Where("server_members.server_id = servers.id").Where("server_members.user_id = ?", userID)
				return nil
			},
		)

	return quert.Find(ctx)
}

func (r *Repository) CreateNewServer(ctx context.Context, userID string, serverName string) (*domainEntity.Server, error) {
	user, err := gorm.G[domainEntity.User](r.db).Where("id = ?", userID).First(ctx)
	if err != nil {
		return nil, errors.Join(ErrUndefinedUser, err)
	}

	newServer := &domainEntity.Server{
		Server: sharedEntity.Server{
			ID: uuid.New().String(),
		},
		Members: []domainEntity.User{user},
	}

	err = gorm.G[domainEntity.Server](r.db).Create(ctx, newServer)
	if err != nil {
		return nil, err
	}

	return newServer, nil
}
