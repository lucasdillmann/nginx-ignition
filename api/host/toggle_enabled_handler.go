package host

import (
	"dillmann.com.br/nginx-ignition/core/host"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type toggleEnabledHandler struct {
	getCommand  *host.GetCommand
	saveCommand *host.SaveCommand
}

func (h toggleEnabledHandler) handle(context *gin.Context) {
	id, err := uuid.Parse(context.Param("id"))
	if err != nil {
		context.Status(http.StatusNotFound)
		return
	}

	data, err := (*h.getCommand)(id)
	if err != nil {
		panic(err)
	}

	if data == nil {
		context.Status(http.StatusNotFound)
		return
	}

	data.Enabled = !data.Enabled
	err = (*h.saveCommand)(data)

	if err != nil {
		panic(err)
	}

	context.Status(http.StatusNoContent)
}
