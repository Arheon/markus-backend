package createnewmessage

import (
	"github.com/Arheon/markus-backend/internal/auth/domain/entity"
	createnewmessage "github.com/Arheon/markus-backend/internal/message/application/command/create_new_message"
	"github.com/Arheon/markus-backend/internal/message/domain/repository"
	"github.com/Arheon/markus-backend/internal/shared/infrastructure/helpers"
	"github.com/Arheon/markus-backend/pkg/outbox"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	serverRepo "github.com/Arheon/markus-backend/internal/server/domain/repository"
)

type Handler struct {
	logger     *logrus.Logger
	repo       repository.MessageRepository
	memberRepo serverRepo.MemberRepository
	outbox     *outbox.Publisher
}

func NewHandler(
	publisher *outbox.Publisher,
	logger *logrus.Logger,
	repo repository.MessageRepository,
	memberRepo serverRepo.MemberRepository,
) *Handler {
	return &Handler{
		memberRepo: memberRepo,
		logger:     logger,
		outbox:     publisher,
		repo:       repo,
	}
}

func (h *Handler) Handle(ctx *gin.Context) {
	identity, exists := ctx.Get("id")
	if !exists {
		helpers.AbortWithAccessDeniedErrorJSON(ctx)
		return
	}

	user := identity.(*entity.User)

	var f form
	if err := ctx.ShouldBindJSON(&f); err != nil {
		h.logger.Error("Unexpected body bindig error ", err.Error())
		helpers.AbortWithBadRequestErrorJSON(ctx)
		return
	}

	cmd := createnewmessage.NewCommand(h.memberRepo, h.repo, h.outbox)
	result, err := cmd.Handle(ctx, user.ID, f.ServerID, f.RoomID, f.Value)
	if err != nil {
		h.logger.Error("Unexpected error ", err.Error())
		helpers.AbortWithError(ctx, err)
		return
	}

	helpers.SuccessJSONV2(ctx, result)
}
