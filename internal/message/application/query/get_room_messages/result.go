package getroommessages

type Result struct {
	RoomID   string          `json:"room_id"`
	Messages []ResultMessage `json:"messages"`
}

type ResultMessage struct {
	MemberID string `json:"member_id"`
	Value    string `json:"value"`
}
