package getserverrooms

type Result struct {
	ServerID             string               `json:"server_id"`
	RoomCategories       []ResultRoomCategory `json:"room_categories"`
	RoomsWithoutCategory []ResultRoom         `json:"rooms_without_category"`
}

type ResultRoomCategory struct {
	CategoryID   string       `json:"category_id"`
	CategoryName string       `json:"category_name"`
	Rooms        []ResultRoom `json:"rooms"`
}

type ResultRoom struct {
	RoomID   string `json:"room_id"`
	RoomName string `json:"room_name"`
}
