package health

import (
	"github.com/Arheon/markus-backend/internal/shared/application/handler/health"
	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

func InitController(injector *do.Injector, router *gin.Engine) {
	v1 := router.Group("v1")
	{
		healthGroup := v1.Group("health")

		handler := health.NewHandler(injector)
		healthGroup.GET("/:type", handler.Handle)
	}
}
