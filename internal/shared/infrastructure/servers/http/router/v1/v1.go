package v1

import (
	"github.com/Arheon/markus-backend/internal/shared/infrastructure/servers/http/router/v1/auth"
	"github.com/Arheon/markus-backend/internal/shared/infrastructure/servers/http/router/v1/health"
	"github.com/Arheon/markus-backend/internal/shared/infrastructure/servers/http/router/v1/user"
	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

func InitV1(injector do.Injector, router *gin.Engine) error {
	v1 := router.Group("v1")
	if err := health.InitController(injector, v1); err != nil {
		return err
	}

	if err := auth.InitContrioller(injector, v1); err != nil {
		return err
	}

	if err := user.InitContrioller(injector, v1); err != nil {
		return err
	}

	return nil
}
