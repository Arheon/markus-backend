package getroommessages

import (
	getroommessages "github.com/Arheon/markus-backend/internal/message/application/query/get_room_messages"
	"github.com/Arheon/markus-backend/internal/message/domain/repository"
	"github.com/Arheon/markus-backend/internal/shared/infrastructure/helpers"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	logger *logrus.Logger
	repo   repository.MessageRepository
}

func NewHandler(
	logger *logrus.Logger,
	repo repository.MessageRepository,
) *Handler {
	return &Handler{
		logger: logger,
		repo:   repo,
	}
}

func (h *Handler) Handle(ctx *gin.Context) {
	var f form
	if err := ctx.ShouldBindJSON(&f); err != nil {
		h.logger.Error("Unexpected body binding error ", err.Error())
		helpers.AbortWithBadRequestErrorJSON(ctx)
		return
	}

	query := getroommessages.NewQuery(h.repo)
	result, err := query.Handle(ctx, f.RoomID)
	if err != nil {
		helpers.AbortWithError(ctx, err)
		return
	}

	helpers.SuccessJSONV2(ctx, result)
}
