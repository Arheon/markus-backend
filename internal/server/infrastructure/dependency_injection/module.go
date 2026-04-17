package dependencyinjection

import (
	"github.com/Arheon/markus-backend/internal/server/domain/repository"
	"github.com/Arheon/markus-backend/internal/server/infrastructure/repository/room"
	"github.com/Arheon/markus-backend/internal/server/infrastructure/repository/server"
	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

func InitModule(injector do.Injector) {
	do.Provide(injector, func(i do.Injector) (repository.ServerRepository, error) {
		db, err := do.InvokeAs[*gorm.DB](i)
		if err != nil {
			return nil, err
		}

		return server.NewRepository(db), nil
	})

	do.Provide(injector, func(i do.Injector) (repository.RoomRepository, error) {
		db, err := do.InvokeAs[*gorm.DB](i)
		if err != nil {
			return nil, err
		}

		return room.NewRepository(db), nil
	})
}
