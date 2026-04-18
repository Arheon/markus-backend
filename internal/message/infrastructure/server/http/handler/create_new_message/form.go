package createnewmessage

type form struct {
	MemberID string `json:"member_id" binding:"required,uuid"`
	RoomID   string `json:"room_id" binding:"required,uuid"`
	Value    string `json:"value" binding:"required"`
}
