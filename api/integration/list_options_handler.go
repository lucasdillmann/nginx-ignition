package integration

import (
	"net/http"

	"github.com/gin-gonic/gin"

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

	tcpOnly := false
	if ctx.Query("tcpOnly") == "true" {
		tcpOnly = true
	}

	page, err := h.commands.ListOptions(ctx.Request.Context(), id, pageSize, pageNumber, searchTerms, tcpOnly)
	if err != nil {
		panic(err)
	}

	if page == nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	ctx.JSON(http.StatusOK, pagination.Convert(page, toOptionDto))
}
