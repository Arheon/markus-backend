package createnewmessage

type form struct {
	ServerID string `json:"server_id" binding:"required,uuid"`
	RoomID   string `json:"room_id" binding:"required,uuid"`
	Value    string `json:"value" binding:"required"`
}
