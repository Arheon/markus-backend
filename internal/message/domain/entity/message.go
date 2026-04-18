package entity

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID        uuid.UUID `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime `gorm:"index"`
	Value     string
	MemberID  string
	RoomID    string
}

func (m *Message) TableName() string {
	return "messages"
}
