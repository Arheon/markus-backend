package createroom

type form struct {
	ServerID       string  `json:"server_id" binding:"required,uuid"`
	RoomName       string  `json:"room_name" binding:"required"`
	RoomCategoryID *string `json:"room_category" binding:"omitempty,uuid"`
}
