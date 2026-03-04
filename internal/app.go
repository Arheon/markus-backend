package internal

import (
	"log"
	"time"

	"github.com/Arheon/markus-backend/configs"
	"github.com/Arheon/markus-backend/internal/shared/infrastructure/servers/http"
	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/samber/do/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/Arheon/markus-backend/internal/auth/application/handler/identity"
	"github.com/Arheon/markus-backend/internal/auth/application/handler/login"
	"github.com/Arheon/markus-backend/internal/auth/application/handler/logout"
	authRepository "github.com/Arheon/markus-backend/internal/auth/domain/repository"
	jwtMiddlewareHelpers "github.com/Arheon/markus-backend/internal/auth/infrastructure/middleware/jwt"
)

func InitApp(env string) error {
	logger := logrus.New()

	injector := do.New()
	cfg, err := config.NewConfig(env)
	if err != nil {
		logger.Error("Undefined config error")
		return err
	}

	dbConf := cfg.Database

	do.Provide(injector, func(i do.Injector) (*gorm.DB, error) {
		db, err := gorm.Open(postgres.Open(dbConf.DSN), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
		})

		if nil != err {
			log.Fatalln("Connection to database failed", err)
		} else {
			log.Println("Connected to database successfully", nil)
		}

		return db, err
	})

	do.Provide(injector, func(i do.Injector) (*logrus.Logger, error) {
		return logger, nil
	})

	do.Provide(injector, func(i do.Injector) (*config.Config, error) {
		return cfg, nil
	})

	do.Provide(injector, func(i do.Injector) (*jwt.GinJWTMiddleware, error) {
		config, err := do.InvokeAs[*config.Config](i)
		if err != nil {
			return nil, err
		}

		userRepo, err := do.InvokeAs[*authRepository.UserRepository](i)
		if err != nil {
			return nil, err
		}

		identityHandler := identity.Handler{}
		loginHandler := login.NewHandler(*userRepo)
		logoutHandler := logout.NewHandler()

		authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
			Realm:       "test zone",
			Key:         []byte("secret key"),
			Timeout:     time.Hour,
			MaxRefresh:  time.Hour,
			IdentityKey: config.Auth.IdentityKey,
			PayloadFunc: jwtMiddlewareHelpers.PayloadFunc,

			IdentityHandler: identityHandler.Handle,
			Authenticator:   loginHandler.HandleAuthenticate,
			Authorizer:      loginHandler.HandleAuthorize,
			Unauthorized:    logoutHandler.Handle,
			LogoutResponse:  logoutHandler.LogoutResponse,
			TokenLookup:     "header: Authorization, query: token, cookie: jwt",
			TokenHeadName:   "Bearer",
			TimeFunc:        time.Now,
		})

		return authMiddleware, nil
	})

	http.Init(env, injector)

	return nil
}
