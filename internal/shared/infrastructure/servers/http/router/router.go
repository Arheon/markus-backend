package router

import (
	v1 "github.com/Arheon/markus-backend/internal/shared/infrastructure/servers/http/router/v1"
	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

func NewRouter(c do.Injector, router *gin.Engine) (*gin.Engine, error) {
	if router == nil {
		router = gin.New()
	}

	router.Use(gin.Recovery())

	if err := v1.InitV1(c, router); err != nil {
		return nil, err
	}
	return router, nil
}
