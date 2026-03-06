package jwt

import (
	"github.com/Arheon/markus-backend/internal/auth/domain/entity"
	"github.com/golang-jwt/jwt/v5"
)

func PayloadFunc(data any) jwt.MapClaims {
	if v, ok := data.(*entity.User); ok {
		return jwt.MapClaims{
			"id": v.Username,
		}
	}
	return jwt.MapClaims{}
}
