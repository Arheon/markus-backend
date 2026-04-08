package entity

type Server struct {
	ID             string `gorm:"primaryKey"`
	Name           string
	Members        []Member `gorm:"many2many:server_members;"`
	RoomCategories []RoomCategory
	Rooms          []Room
}
