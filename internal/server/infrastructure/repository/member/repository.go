package member

import (
	"context"

	"github.com/Arheon/markus-backend/internal/server/domain/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MemberRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *MemberRepository {
	return &MemberRepository{
		db: db,
	}
}

func (r *MemberRepository) GetMemberByUserAndServer(ctx context.Context, serverID string, userID string) (entity.Member, error) {
	return gorm.G[entity.Member](r.db).
		Joins(
			clause.JoinTarget{
				Type:  clause.InnerJoin,
				Table: "server_members",
			},
			func(db gorm.JoinBuilder, joinTable, curTable clause.Table) error {
				db.
					Where("server_members.server_id = servers.id").
					Where("server_members.user_id = ?", userID).
					Where("server_members.server_id = ?", serverID)
				return nil
			},
		).First(ctx)
}
