package access_list

import (
	"dillmann.com.br/nginx-ignition/api/common/pagination"
	"dillmann.com.br/nginx-ignition/core/access_list"
	"github.com/gin-gonic/gin"
	"net/http"
)

type listHandler struct {
	commands *access_list.Commands
}

func (h listHandler) handle(ctx *gin.Context) {
	pageSize, pageNumber, searchTerms, err := pagination.ExtractPaginationParameters(ctx)
	if err != nil {
		panic(err)
	}

	page, err := h.commands.List(ctx.Request.Context(), pageSize, pageNumber, searchTerms)
	if err != nil {
		panic(err)
	}

	ctx.JSON(http.StatusOK, pagination.Convert(page, toDto))
}
