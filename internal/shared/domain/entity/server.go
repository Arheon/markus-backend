package entity

type Server struct {
	ID   string `gorm:"primaryKey"`
	Name string
}

func (Server) TableName() string {
	return "servers"
}
