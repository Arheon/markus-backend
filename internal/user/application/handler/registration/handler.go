package registration

import (
	globalHelpers "github.com/Arheon/markus-backend/internal/shared/infrastructure/helpers"
	createuser "github.com/Arheon/markus-backend/internal/user/domain/command/create_user"
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
		logger,
		userRepository,
	}
}

func (h *Handler) Handle(ctx *gin.Context) {
	var form registrationForm
	if err := ctx.ShouldBind(&form); err != nil {
		globalHelpers.AbortWithBadRequestErrorJSON(ctx)
		return
	}

	command := createuser.NewCommand(ctx, h.userRepository)
	h.logger.Debug("Try to create new user")

	user, err := command.Handle(form.Username, form.Password)
	if err != nil {
		globalHelpers.AbortWithError(ctx, err)
		return
	}

	globalHelpers.SuccessJSONV2(ctx, map[string]any{
		"message": "User create successfuly",
		"user":    user,
	})
}
