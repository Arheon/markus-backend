package createnewserver

import (
	"context"

	domainEntity "github.com/Arheon/markus-backend/internal/server/domain/entity"
	"github.com/Arheon/markus-backend/internal/server/domain/event"
	"github.com/Arheon/markus-backend/internal/server/domain/repository"
	"github.com/Arheon/markus-backend/pkg/outbox"
)

type Command struct {
	repo      repository.ServerRepository
	publisher *outbox.Publisher
}

func NewCommand(
	repo repository.ServerRepository,
	publisher *outbox.Publisher,
) *Command {
	return &Command{
		repo:      repo,
		publisher: publisher,
	}
}

func (cmd *Command) Handle(ctx context.Context, userID string, serverName string) (*domainEntity.Server, error) {
	return outbox.Publish(cmd.publisher, ctx, func(c *outbox.EventCollector) (*domainEntity.Server, error) {
		server, err := cmd.repo.CreateNewServer(ctx, userID, serverName)
		if err != nil {
			return nil, err
		}

		c.Add(&event.CreateServerEvent{
			ServerID: server.ID,
			Name:     serverName,
			MemberIDs: []string{
				userID,
			},
		})

		return server, nil
	})
}
