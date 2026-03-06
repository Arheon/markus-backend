package internal

import (
	"log"
	"time"

	config "github.com/Arheon/markus-backend/configs"
	serverDomain "github.com/Arheon/markus-backend/internal/server/domain/entity"
	serverRepository "github.com/Arheon/markus-backend/internal/server/domain/repository"
	serverServerRepository "github.com/Arheon/markus-backend/internal/server/infrastructure/repository/server"
	sharedDomainEntity "github.com/Arheon/markus-backend/internal/shared/domain/entity"
	"github.com/Arheon/markus-backend/internal/shared/infrastructure/servers/http"
	userDomainEntity "github.com/Arheon/markus-backend/internal/user/domain/entity"
	userRepository "github.com/Arheon/markus-backend/internal/user/domain/repository"
	userUserRepository "github.com/Arheon/markus-backend/internal/user/infrastructure/repository/user"
	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/samber/do/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/Arheon/markus-backend/internal/auth/application/handler/identity"
	"github.com/Arheon/markus-backend/internal/auth/application/handler/login"
	"github.com/Arheon/markus-backend/internal/auth/application/handler/logout"
	authDomainEntity "github.com/Arheon/markus-backend/internal/auth/domain/entity"
	authRepository "github.com/Arheon/markus-backend/internal/auth/domain/repository"
	jwtMiddlewareHelpers "github.com/Arheon/markus-backend/internal/auth/infrastructure/middleware/jwt"
	authUserRepository "github.com/Arheon/markus-backend/internal/auth/infrastructure/repository/user"
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

		db.AutoMigrate(
			&authDomainEntity.User{},
			&sharedDomainEntity.User{},
			&userDomainEntity.User{},
			&serverDomain.User{},
			&serverDomain.Message{},
			&serverDomain.Room{},
			&serverDomain.RoomCategory{},
			&serverDomain.Server{},
		)

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
		logger, err := do.InvokeAs[*logrus.Logger](i)
		if err != nil {
			return nil, err
		}

		config, err := do.InvokeAs[*config.Config](i)
		if err != nil {
			return nil, err
		}

		userRepo, err := do.InvokeAs[authRepository.UserRepository](i)
		if err != nil {
			return nil, err
		}

		identityHandler := identity.Handler{}
		loginHandler := login.NewHandler(logger, userRepo)
		logoutHandler := logout.NewHandler()

		identityKey := config.Auth.IdentityKey
		if identityKey == "" {
			identityKey = "id"
		}

		authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
			Realm:       "test zone",
			Key:         []byte("secret key"),
			Timeout:     time.Hour,
			MaxRefresh:  time.Hour,
			IdentityKey: identityKey,
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

	do.Provide(injector, func(i do.Injector) (authRepository.UserRepository, error) {
		db, err := do.InvokeAs[*gorm.DB](i)
		if err != nil {
			return nil, err
		}

		return authUserRepository.NewRepository(db), nil
	})

	do.Provide(injector, func(i do.Injector) (userRepository.UserRepository, error) {
		db, err := do.InvokeAs[*gorm.DB](i)
		if err != nil {
			return nil, err
		}

		return userUserRepository.NewRepository(db), nil
	})

	do.Provide(injector, func(i do.Injector) (serverRepository.ServerRepository, error) {
		db, err := do.InvokeAs[*gorm.DB](i)
		if err != nil {
			return nil, err
		}

		return serverServerRepository.NewRepository(db), nil
	})

	http.Init(env, injector)

	return nil
}
