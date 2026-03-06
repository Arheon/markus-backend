package entity

import sharedEntity "github.com/Arheon/markus-backend/internal/shared/domain/entity"

type User struct {
	sharedEntity.User `gorm:"embedded"`
	Messages          []Message
	Servers           []Server `gorm:"many2many:server_members;"`
}
