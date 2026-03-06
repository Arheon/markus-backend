package entity

type Room struct {
	ID             string `gorm:"primaryKey"`
	ServerID       string
	RoomCategoryID *string
}
