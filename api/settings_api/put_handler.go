package settings_api

import (
	"dillmann.com.br/nginx-ignition/core/settings"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type putHandler struct {
	command *settings.SaveCommand
}

func (h putHandler) handle(context *gin.Context) {
	var payload SettingsDto
	if err := json.NewDecoder(context.Request.Body).Decode(&payload); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.New().Struct(payload); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	domain := toDomain(&payload)
	if err := (*h.command)(domain); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.Status(http.StatusNoContent)
}
