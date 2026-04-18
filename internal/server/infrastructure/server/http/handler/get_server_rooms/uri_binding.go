package getserverrooms

type uribinding struct {
	ServerID string `uri:"serverID" binding:"required,uuid"`
}
