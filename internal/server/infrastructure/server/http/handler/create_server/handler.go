package createnewserver

import (
	"github.com/Arheon/markus-backend/internal/auth/domain/entity"
	createnewserver "github.com/Arheon/markus-backend/internal/server/application/command/create_new_server"
	"github.com/Arheon/markus-backend/internal/server/domain/repository"
	"github.com/Arheon/markus-backend/internal/shared/infrastructure/helpers"
	"github.com/Arheon/markus-backend/pkg/outbox"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	repo      repository.ServerRepository
	publisher *outbox.Publisher
}

func NewHandler(
	repo repository.ServerRepository,
	publisher *outbox.Publisher,
) *Handler {
	return &Handler{
		repo:      repo,
		publisher: publisher,
	}
}

func (h *Handler) Handle(ctx *gin.Context) {
	identity, exists := ctx.Get("id")
	if !exists {
		helpers.AbortWithUnauthorizedErrorJSON(ctx)
		return
	}

	user := identity.(*entity.User)

	var form form
	if err := ctx.ShouldBind(&form); err != nil {
		helpers.AbortWithError(ctx, err)
		return
	}

	cmd := createnewserver.NewCommand(h.repo, h.publisher)
	server, err := cmd.Handle(ctx, user.ID, form.ServerName)
	if err != nil {
		helpers.AbortWithError(ctx, err)
		return
	}

	helpers.SuccessJSONV2(ctx, server)
}
