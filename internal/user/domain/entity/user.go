package entity

import "github.com/Arheon/markus-backend/internal/shared/domain/entity"

type User struct {
	entity.User `gorm:"embedded"`
	Name        string ``
}
