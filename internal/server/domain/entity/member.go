package entity

import "github.com/Arheon/markus-backend/internal/user/domain/entity"

type Member struct {
	ID      string `gorm:"primaryKey"`
	UserID  string
	Servers []Server `gorm:"many2many:server_members"`
	User    entity.User
}
