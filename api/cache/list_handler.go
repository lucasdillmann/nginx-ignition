package cache

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"dillmann.com.br/nginx-ignition/api/common/pagination"
	"dillmann.com.br/nginx-ignition/core/cache"
)

type listHandler struct {
	commands *cache.Commands
}

func (h listHandler) handle(ctx *gin.Context) {
	pageSize, pageNumber, searchTerms, err := pagination.ExtractPaginationParameters(ctx)
	if err != nil {
		panic(err)
	}

	page, err := h.commands.List(ctx, pageSize, pageNumber, searchTerms)
	if err != nil {
		panic(err)
	}

	response := pagination.Convert(page, toResponseDto)
	ctx.JSON(http.StatusOK, response)
}
