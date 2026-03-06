package entity

import sharedEntity "github.com/Arheon/markus-backend/internal/shared/domain/entity"

type Server struct {
	sharedEntity.Server `gorm:"embedded"`
	Members             []User `gorm:"many2many:server_members;"`
	RoomCategories      []RoomCategory
	Rooms               []Room
}
