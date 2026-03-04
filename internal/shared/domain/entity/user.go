package entity

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"column:username"`
	Password string `gorm:"column:password"`
}

func (User) TableName() string {
	return "users"
}
