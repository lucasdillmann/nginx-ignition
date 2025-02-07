package integration

import (
	"dillmann.com.br/nginx-ignition/api/common/pagination"
	"dillmann.com.br/nginx-ignition/core/integration"
	"github.com/gin-gonic/gin"
	"net/http"
)

type listOptionsHandler struct {
	command *integration.ListOptionsCommand
}

func (h listOptionsHandler) handle(context *gin.Context) {
	pageSize, pageNumber, searchTerms, err := pagination.ExtractPaginationParameters(context)
	if err != nil {
		panic(err)
	}

	id := context.Param("id")
	if id == "" {
		context.Status(http.StatusNotFound)
		return
	}

	page, err := (*h.command)(id, pageSize, pageNumber, searchTerms)
	if err != nil {
		panic(err)
	}

	if page == nil {
		context.Status(http.StatusNotFound)
		return
	}

	context.JSON(http.StatusOK, pagination.Convert(page, toOptionDto))
}
