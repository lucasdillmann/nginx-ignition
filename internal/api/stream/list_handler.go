package stream

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lucasdillmann/nginx-ignition/internal/api/common/pagination"
	"github.com/lucasdillmann/nginx-ignition/internal/core/stream"
)

type listHandler struct {
	commands stream.Commands
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
