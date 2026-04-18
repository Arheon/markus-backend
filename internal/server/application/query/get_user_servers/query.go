package getuserservers

import (
	"context"

	"github.com/Arheon/markus-backend/internal/server/domain/repository"
)

type Query struct {
	ctx  context.Context
	repo repository.ServerRepository
}

func NewQuery(ctx context.Context, repo repository.ServerRepository) *Query {
	return &Query{
		ctx:  ctx,
		repo: repo,
	}
}

func (q *Query) Handle(userID string) (*Result, error) {
	servers, err := q.repo.GetAllServersByUserID(q.ctx, userID)
	if err != nil {
		return nil, err
	}

	var serverResults []ResultServer
	for _, server := range servers {
		serverResult := ResultServer{
			ID:      server.ID,
			Name:    server.Name,
			Members: []ResultMember{},
		}

		if len(server.Members) == 0 {
			serverResults = append(serverResults, serverResult)
			continue
		}

		for _, member := range server.Members {
			memberResult := &ResultMember{
				ID:   member.ID,
				Name: member.User.Name,
			}

			serverResult.Members = append(serverResult.Members, *memberResult)
		}
		serverResults = append(serverResults, serverResult)
	}

	return &Result{
		Servers: serverResults,
	}, nil
}
