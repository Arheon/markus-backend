package gorm

import (
	"context"
	"time"

	"github.com/Arheon/markus-backend/pkg/outbox"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
	return &Store{db: db}
}

func (s *Store) Save(ctx context.Context, msg outbox.Message) error {
	return s.db.WithContext(ctx).Create(&msg).Error
}

func (s *Store) FetchPending(ctx context.Context, limit int) ([]outbox.Message, error) {
	msgs, err := gorm.G[outbox.Message](s.db).Where("processed_at IS NULL").Limit(limit).Find(ctx)
	return msgs, err
}

func (s *Store) MarkProcessed(ctx context.Context, id uuid.UUID) error {
	return s.db.WithContext(ctx).
		Model(&outbox.Message{}).
		Where("id = ?", id).
		Update("processed_at", time.Now()).Error
}

func (s *Store) MarkFailed(ctx context.Context, id uuid.UUID, err error) error {

	return s.db.WithContext(ctx).
		Model(&outbox.Message{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"attempts": gorm.Expr("attempts + 1"),
			"error":    err.Error(),
		}).Error
}
