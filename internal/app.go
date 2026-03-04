package internal

import (
	"log"

	"github.com/Arheon/markus-backend/config"
	"github.com/samber/do/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

	return nil
}
