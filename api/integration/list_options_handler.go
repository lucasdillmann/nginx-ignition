package integration

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/api/common/pagination"
	"dillmann.com.br/nginx-ignition/core/integration"
)

type listOptionsHandler struct {
	commands *integration.Commands
}

func (h listOptionsHandler) handle(ctx *gin.Context) {
	pageSize, pageNumber, searchTerms, err := pagination.ExtractPaginationParameters(ctx)
	if err != nil {
		panic(err)
	}

	id := ctx.Param("id")
	if id == "" {
		ctx.Status(http.StatusNotFound)
		return
	}

	uuidValue, err := uuid.Parse(id)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	tcpOnly := ctx.Query("tcpOnly") == "true"

	page, err := h.commands.ListOptions(ctx.Request.Context(), uuidValue, pageSize, pageNumber, searchTerms, tcpOnly)
	if err != nil {
		panic(err)
	}

	if page == nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	ctx.JSON(http.StatusOK, pagination.Convert(page, toOptionDto))
}
