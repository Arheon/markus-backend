package getroommessages

type form struct {
	RoomID string `json:"room_id" binding:"required,uuid"`
}
