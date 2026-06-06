package notification

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/notification"
)

type putConfigurationHandler struct {
	commands notification.Commands
}

func (h putConfigurationHandler) handle(ctx *gin.Context) {
	userID, ok := currentUserID(ctx)
	if !ok {
		return
	}

	payload := &configurationRequest{}
	if err := ctx.BindJSON(payload); err != nil {
		panic(err)
	}

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	domainModel, err := toConfigurationDomain(userID, id, payload)
	if err != nil {
		panic(err)
	}

	if err := h.commands.SaveConfiguration(
		ctx.Request.Context(),
		userID,
		domainModel,
	); err != nil {
		if errors.Is(err, notification.ErrConfigurationNotFound) {
			ctx.Status(http.StatusNotFound)
			return
		}
		panic(err)
	}

	ctx.Status(http.StatusNoContent)
}
