package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/lucasdillmann/nginx-ignition/internal/api/core/authorization"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/user"
)

type logoutHandler struct {
	commands user.Commands
}

func (h logoutHandler) handle(ctx *gin.Context) {
	subject := authorization.CurrentSubject(ctx)

	sessionID, err := uuid.Parse(subject.TokenID)
	if err != nil {
		panic(err)
	}

	if err := h.commands.StopSession(ctx.Request.Context(), sessionID); err != nil {
		panic(err)
	}

	ctx.Status(http.StatusNoContent)
}
