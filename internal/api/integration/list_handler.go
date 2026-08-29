package integration

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"nginx-ignition/internal/api/common/pagination"
	"nginx-ignition/internal/core/integration"
)

type listHandler struct {
	commands integration.Commands
}

func (h listHandler) handle(ctx *gin.Context) {
	pageSize, pageNumber, searchTerms, err := pagination.ExtractPaginationParameters(ctx)
	if err != nil {
		panic(err)
	}

	enabledOnly := ctx.Query("enabledOnly") == "true"

	page, err := h.commands.List(
		ctx.Request.Context(),
		pageSize,
		pageNumber,
		searchTerms,
		enabledOnly,
	)
	if err != nil {
		panic(err)
	}

	ctx.JSON(http.StatusOK, pagination.Convert(page, toDTO))
}
