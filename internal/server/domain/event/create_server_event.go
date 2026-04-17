package event

type CreateServerEvent struct {
	ServerID  string
	Name      string
	MemberIDs []string
}

func (e *CreateServerEvent) EventName() string {
	return "event.server.create_server"
}
