package createroom

type form struct {
	RoomName       string  `json:"room_name" binding:"required"`
	RoomCategoryID *string `json:"room_category_id" binding:"omitempty,uuid"`
}
