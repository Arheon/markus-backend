package dependencyinjection

import (
	"time"

	config "github.com/Arheon/markus-backend/configs"
	"github.com/Arheon/markus-backend/internal/auth/domain/repository"
	jwtMiddleware "github.com/Arheon/markus-backend/internal/auth/infrastructure/middleware/jwt"
	"github.com/Arheon/markus-backend/internal/auth/infrastructure/repository/user"
	"github.com/Arheon/markus-backend/internal/auth/infrastructure/server/http/handler/identity"
	"github.com/Arheon/markus-backend/internal/auth/infrastructure/server/http/handler/login"
	"github.com/Arheon/markus-backend/internal/auth/infrastructure/server/http/handler/logout"
	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/samber/do/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func InitModule(injector do.Injector) {
	do.Provide(injector, func(i do.Injector) (repository.UserRepository, error) {
		db, err := do.InvokeAs[*gorm.DB](i)
		if err != nil {
			return nil, err
		}

		return user.NewRepository(db), nil
	})

	do.Provide(injector, func(i do.Injector) (*jwt.GinJWTMiddleware, error) {
		logger, err := do.InvokeAs[*logrus.Logger](i)
		if err != nil {
			return nil, err
		}

		config, err := do.InvokeAs[*config.Config](i)
		if err != nil {
			return nil, err
		}

		userRepo, err := do.InvokeAs[repository.UserRepository](i)

		if err != nil {
			return nil, err
		}

		identityHandler := identity.NewHandler(userRepo, logger)
		loginHandler := login.NewHandler(logger, userRepo)
		logoutHandler := logout.NewHandler()

		identityKey := config.Auth.IdentityKey
		if identityKey == "" {
			identityKey = "id"
		}

		authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
			Realm:       config.Env,
			Key:         []byte(config.Secret),
			Timeout:     time.Hour,
			MaxRefresh:  time.Hour,
			IdentityKey: identityKey,
			PayloadFunc: jwtMiddleware.PayloadFunc,

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
}
