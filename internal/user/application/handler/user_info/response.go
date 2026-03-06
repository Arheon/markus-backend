package userinfo

import "github.com/Arheon/markus-backend/internal/user/domain/entity"

type Response struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func FromUser(user *entity.User) *Response {
	return &Response{
		ID:   user.ID,
		Name: user.Name,
	}
}
