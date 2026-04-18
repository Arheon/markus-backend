package getserverrooms

import (
	getserverrooms "github.com/Arheon/markus-backend/internal/server/application/query/get_server_rooms"
	"github.com/Arheon/markus-backend/internal/server/domain/repository"
	"github.com/Arheon/markus-backend/internal/shared/infrastructure/helpers"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	logger   *logrus.Logger
	roomRepo repository.RoomRepository
}

func NewHandler(
	logger *logrus.Logger,
	roomRepo repository.RoomRepository,
) *Handler {
	return &Handler{
		logger:   logger,
		roomRepo: roomRepo,
	}
}

func (h *Handler) Handle(ctx *gin.Context) {
	var binding uribinding
	if err := ctx.ShouldBindUri(&binding); err != nil {
		h.logger.Error("Unexpected uri binding error ", err.Error())
		helpers.AbortWithBadRequestErrorJSON(ctx)
		return
	}

	query := getserverrooms.NewQuery(h.roomRepo)
	result, err := query.Handle(ctx, binding.ServerID)
	if err != nil {
		h.logger.Error("Unexpected query error ", err.Error())
		helpers.AbortWithError(ctx, err)
		return
	}

	helpers.SuccessJSONV2(ctx, result)
}
