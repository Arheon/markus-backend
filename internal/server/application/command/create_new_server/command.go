package createnewserver

import (
	"context"
	"encoding/json"

	domainEntity "github.com/Arheon/markus-backend/internal/server/domain/entity"
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
	server, err := cmd.repo.CreateNewServer(ctx, userID, serverName)
	if err != nil {
		return nil, err
	}

	payload, err := json.Marshal(server)
	if err != nil {
		return nil, err
	}

	message := outbox.NewMessage(payload)
	message.Type = "ServerCreate"

	cmd.publisher.Publish(ctx, *message)

	return server, nil
}
