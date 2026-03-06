package userinfo

import (
	globalHelpers "github.com/Arheon/markus-backend/internal/shared/infrastructure/helpers"
	getuserinfo "github.com/Arheon/markus-backend/internal/user/domain/query/get_user_info"
	"github.com/Arheon/markus-backend/internal/user/domain/repository"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	logger         *logrus.Logger
	userRepository repository.UserRepository
}

func NewHandler(logger *logrus.Logger, userRepository repository.UserRepository) *Handler {
	return &Handler{
		logger:         logger,
		userRepository: userRepository,
	}
}

func (h *Handler) Handle(ctx *gin.Context) {
	var bindings UserIDBindings
	if err := ctx.ShouldBindUri(&bindings); err != nil {
		globalHelpers.AbortWithBadRequestErrorJSON(ctx)
	}

	query := getuserinfo.NewQuery(ctx, h.userRepository)
	h.logger.Debug("Try to get user infi")
	user, err := query.Handle(bindings.UserID)
	if err != nil {
		globalHelpers.AbortWithError(ctx, err)
		return
	}

	globalHelpers.SuccessJSONV2(ctx, FromUser(user))
}
