package createnewserver

import (
	"context"

	domainEntity "github.com/Arheon/markus-backend/internal/server/domain/entity"
	"github.com/Arheon/markus-backend/internal/server/domain/repository"
)

type Command struct {
	repo repository.ServerRepository
}

func NewCommand(repo repository.ServerRepository) *Command {
	return &Command{
		repo: repo,
	}
}

func (cmd *Command) Handle(ctx context.Context, userID string, serverName string) (*domainEntity.Server, error) {
	return cmd.repo.CreateNewServer(ctx, userID, serverName)
}
