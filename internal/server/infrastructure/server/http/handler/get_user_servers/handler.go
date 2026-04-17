package getuserservers

import (
	"github.com/Arheon/markus-backend/internal/auth/domain/entity"
	getuserservers "github.com/Arheon/markus-backend/internal/server/application/query/get_user_servers"
	"github.com/Arheon/markus-backend/internal/server/domain/repository"
	"github.com/Arheon/markus-backend/internal/shared/infrastructure/helpers"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	repo repository.ServerRepository
}

func NewHandler(repo repository.ServerRepository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) Handle(ctx *gin.Context) {
	identiti, exists := ctx.Get("id")
	if !exists {
		helpers.AbortWithUnauthorizedErrorJSON(ctx)
		return
	}

	user := identiti.(*entity.User)

	query := getuserservers.NewQuery(ctx, h.repo)
	serverList, err := query.Handle(user.ID)
	if err != nil {
		helpers.AbortWithError(ctx, err)
		return
	}

	helpers.SuccessJSONV2(ctx, serverList)
}
