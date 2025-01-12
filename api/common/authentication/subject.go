package authentication

import (
	"dillmann.com.br/nginx-ignition/core/user"
	"github.com/golang-jwt/jwt/v5"
)

type Subject struct {
	TokenID string
	User    *user.User
	claims  *jwt.MapClaims
}
