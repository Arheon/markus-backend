package createnewmessage

type UriBindings struct {
	RoomID string `uri:"userID" binding:"required,uuid"`
}
