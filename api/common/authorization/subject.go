package authorization

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"dillmann.com.br/nginx-ignition/core/user"
)

type Subject struct {
	TokenID string
	User    *user.User
	claims  *jwt.MapClaims
}

func CurrentSubject(ctx *gin.Context) *Subject {
	subject, _ := ctx.Get(RequestSubject)
	if subject == nil {
		return nil
	}

	return subject.(*Subject)
}
