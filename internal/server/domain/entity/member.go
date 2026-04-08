package entity

type Member struct {
	ID      string `gorm:"primaryKey"`
	UserID  string
	Servers []Server `gorm:"many2many:server_members"`
}
