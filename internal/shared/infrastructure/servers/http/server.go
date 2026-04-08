package http

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/Arheon/markus-backend/configs"
	"github.com/Arheon/markus-backend/internal/shared/infrastructure/servers/http/router"
	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
	"github.com/sirupsen/logrus"
)

func Init(env string, injector do.Injector) {
	logger := do.MustInvokeAs[*logrus.Logger](injector)
	cfg := do.MustInvokeAs[*config.Config](injector)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer func() {
		injector.RootScope().ShutdownOnSignals()
		stop()
	}()

	if env == "PROD" {
		gin.SetMode(gin.ReleaseMode)
	}

	r, err := router.NewRouter(injector, nil)
	if err != nil {
		logger.Fatal(err)
	}

	srv := &http.Server{
		Addr:    cfg.Server.Port,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()

	stop()
	logger.Debug("shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown: ", err)
	}

	logger.Debug("Server exiting")

}

