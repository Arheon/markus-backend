package createnewserver

import (
	createnewserver "github.com/Arheon/markus-backend/internal/server/domain/command/create_new_server"
	"github.com/Arheon/markus-backend/internal/server/domain/repository"
	"github.com/Arheon/markus-backend/internal/shared/infrastructure/helpers"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	repo repository.ServerRepository
}

func NewHandler(repo repository.ServerRepository) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) Handle(ctx *gin.Context) {
	var bindings uriBindings
	if err := ctx.ShouldBindUri(&bindings); err != nil {
		helpers.AbortWithError(ctx, err)
		return
	}

	var form form
	if err := ctx.ShouldBind(&form); err != nil {
		helpers.AbortWithError(ctx, err)
		return
	}

	cmd := createnewserver.NewCommand(h.repo)
	server, err := cmd.Handle(ctx, bindings.UserID, form.ServerName)
	if err != nil {
		helpers.AbortWithError(ctx, err)
		return
	}

	helpers.SuccessJSONV2(ctx, server)
}
