package entity

type RoomCategory struct {
	ID       string `gorm:"primaryKey"`
	Name     string
	Rooms    []Room
	ServerID string
}
