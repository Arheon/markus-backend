package jwt

import (
	"github.com/Arheon/markus-backend/internal/shared/domain/entity"
	"github.com/golang-jwt/jwt/v5"
)

var identityKey = "id"

func PayloadFunc(data any) jwt.MapClaims {
	if v, ok := data.(*entity.User); ok {
		return jwt.MapClaims{
			identityKey: v.Username,
		}
	}
	return jwt.MapClaims{}
}
