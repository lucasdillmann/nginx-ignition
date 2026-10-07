package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/lucasdillmann/nginx-ignition/internal/api/core/authorization"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/user"
)

type apiTokenDeleteHandler struct {
	commands user.Commands
}

func (h apiTokenDeleteHandler) handle(ctx *gin.Context) {
	currentSubject := authorization.CurrentSubject(ctx)
	if currentSubject == nil || currentSubject.User == nil {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	if currentSubject.Kind == authorization.APIKind {
		ctx.Status(http.StatusBadRequest)
		return
	}

	tokenID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		panic(err)
	}

	if err := h.commands.DeleteAPIToken(
		ctx.Request.Context(),
		currentSubject.User.ID,
		tokenID,
	); err != nil {
		panic(err)
	}

	ctx.Status(http.StatusNoContent)
}
