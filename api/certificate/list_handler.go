package certificate

import (
	"net/http"

	"dillmann.com.br/nginx-ignition/api/common/pagination"
	"dillmann.com.br/nginx-ignition/core/certificate"
	"github.com/gin-gonic/gin"
)

type listHandler struct {
	commands *certificate.Commands
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

	ctx.JSON(http.StatusOK, pagination.Convert(page, toCertificateResponse))
}
