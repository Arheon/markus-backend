package router

import (
	v1 "github.com/Arheon/markus-backend/internal/shared/infrastructure/servers/http/router/v1"
	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

func NewRouter(c *do.Injector, router *gin.Engine) *gin.Engine {
	if router == nil {
		router = gin.New()
	}

	router.Use(gin.Recovery())

	v1.InitV1(c, router)
	return router
}
