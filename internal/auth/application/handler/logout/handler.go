package logout

import (
	"errors"
	"net/http"

	"github.com/Arheon/markus-backend/internal/auth/domain/entity"
	"github.com/Arheon/markus-backend/internal/shared/infrastructure/helpers"
	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
)

var (
	ErrEmptyCookieToken  = errors.New("cookie token is empty")
	ErrEmptyAuthHeader   = errors.New("auth header is empty")
	ErrInvalidAuthHeader = errors.New("auth header is invalid")
	ErrEmptyQueryToken   = errors.New("query token is empty")
	ErrEmptyParamToken   = errors.New("parameter token is empty")

	ErrUndefinedDependency = errors.New("Try to get undefined dependency")
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Handle(ctx *gin.Context, code int, message string) {
	helpers.SuccessJSONV2(ctx, map[string]any{
		"code":    code,
		"message": message,
	})
}

func (h *Handler) LogoutResponse(ctx *gin.Context) {
	claims := jwt.ExtractClaims(ctx)
	user, exists := ctx.Get("id")

	response := map[string]any{
		"code":    http.StatusOK,
		"message": "Successfully logged out",
	}

	if len(claims) > 0 {
		response["logged_out_user"] = claims["id"]
	}

	if exists {
		switch v := user.(type) {
		case string:
			response["user_info"] = v
		case *entity.User:
			response["user_info"] = v.Username
		default:
			response["user_info"] = user
		}
	}

	helpers.SuccessJSONV2(ctx, response)
}
