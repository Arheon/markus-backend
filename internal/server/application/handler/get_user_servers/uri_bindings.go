package getuserservers

type uriBindings struct {
	UserID string `uri:"userID" binding:"required,uuid"`
}
