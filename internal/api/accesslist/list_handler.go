package accesslist

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"nginx-ignition/internal/api/common/pagination"
	"nginx-ignition/internal/core/accesslist"
)

type listHandler struct {
	commands accesslist.Commands
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

	ctx.JSON(http.StatusOK, pagination.Convert(page, toDTO))
}
