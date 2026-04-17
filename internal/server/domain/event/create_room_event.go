package event

type CreateRoomEvent struct {
	ServerID string `json:"server_name"`
	RoomID   string `json:"room_id"`
}

func (e *CreateRoomEvent) EventName() string {
	return "event.server.create_room"
}
