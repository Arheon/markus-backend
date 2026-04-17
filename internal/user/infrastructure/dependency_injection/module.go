package dependencyinjection

import (
	"github.com/Arheon/markus-backend/internal/user/domain/repository"
	"github.com/Arheon/markus-backend/internal/user/infrastructure/repository/user"
	"github.com/samber/do/v2"
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
}
