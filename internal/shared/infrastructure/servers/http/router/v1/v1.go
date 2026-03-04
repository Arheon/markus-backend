package v1

import (
	"github.com/Arheon/markus-backend/internal/shared/infrastructure/servers/http/router/v1/health"
	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

func InitV1(injector do.Injector, router *gin.Engine) {
	v1 := router.Group("v1")
	health.InitController(injector, v1)
}
