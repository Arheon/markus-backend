package createroom

import (
	createroom "github.com/Arheon/markus-backend/internal/server/application/command/create_room"
	"github.com/Arheon/markus-backend/internal/server/domain/repository"
	"github.com/Arheon/markus-backend/internal/shared/infrastructure/helpers"
	"github.com/Arheon/markus-backend/pkg/outbox"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	repo      repository.RoomRepository
	publisher *outbox.Publisher
	logger    *logrus.Logger
}

func NewHandler(
	repo repository.RoomRepository,
	publisher *outbox.Publisher,
	logger *logrus.Logger,
) *Handler {
	return &Handler{
		repo:      repo,
		publisher: publisher,
		logger:    logger,
	}
}

func (h *Handler) Handle(ctx *gin.Context) {
	var binding uribinding
	if err := ctx.ShouldBindUri(&binding); err != nil {
		h.logger.Error("Unexpected uri binding error ", err.Error())
		helpers.AbortWithBadRequestErrorJSON(ctx)
		return
	}

	var form form
	if err := ctx.ShouldBind(&form); err != nil {
		h.logger.Error("Undefined gin context error ", err)
		helpers.AbortWithBadRequestErrorJSON(ctx)
		return
	}

	cmd := createroom.NewCommand(h.repo, h.publisher)
	roomID, err := cmd.Handle(ctx, form.RoomName, binding.ServerID, form.RoomCategoryID)
	if err != nil {
		helpers.AbortWithError(ctx, err)
	}

	helpers.SuccessJSONV2(ctx, map[string]string{
		"room_id": *roomID,
	})
}
