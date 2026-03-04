package userinfo

type UserIDBindings struct {
	UserID string `uri:"userID" binding:"required,uuid"`
}
