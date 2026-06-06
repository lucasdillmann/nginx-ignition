package notification

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/notification"
)

type createConfigurationHandler struct {
	commands notification.Commands
}

func (h createConfigurationHandler) handle(ctx *gin.Context) {
	userID, ok := currentUserID(ctx)
	if !ok {
		return
	}

	payload := &configurationRequest{}
	if err := ctx.BindJSON(payload); err != nil {
		panic(err)
	}

	domainModel, err := toConfigurationDomain(userID, uuid.New(), payload)
	if err != nil {
		panic(err)
	}

	if err := h.commands.SaveConfiguration(
		ctx.Request.Context(),
		userID,
		domainModel,
	); err != nil {
		panic(err)
	}

	ctx.JSON(http.StatusCreated, map[string]any{"id": domainModel.ID})
}
