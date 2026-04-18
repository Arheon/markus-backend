package entity

type Room struct {
	ID             string `gorm:"primaryKey"`
	RoomName       string
	RoomCategory   *RoomCategory
	ServerID       string
	RoomCategoryID *string
}
