package entity

import "time"

type Message struct {
	ID        string `gorm:"primaryKey"`
	Value     string
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    string
	RoomID    string
}
