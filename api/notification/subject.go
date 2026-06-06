package notification

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/api/common/authorization"
)

func currentUserID(ctx *gin.Context) (uuid.UUID, bool) {
	subject := authorization.CurrentSubject(ctx)
	if subject == nil || subject.User == nil {
		ctx.Status(http.StatusUnauthorized)
		return uuid.Nil, false
	}

	return subject.User.ID, true
}
