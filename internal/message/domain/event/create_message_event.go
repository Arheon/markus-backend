package event

type CreateMessageEvent struct {
	MemberID string `json:"member_id"`
	RoomID   string `json:"room_id"`
	Value    string `json:"value"`
}

func (e *CreateMessageEvent) EventName() string {
	return "event.message.create_message"
}
