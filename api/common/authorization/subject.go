package authorization

import (
	"dillmann.com.br/nginx-ignition/core/user"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Subject struct {
	TokenID string
	User    *user.User
	claims  *jwt.MapClaims
}

func CurrentSubject(ctx *gin.Context) *Subject {
	subject, _ := ctx.Get(RequestSubject)
	return subject.(*Subject)
}
