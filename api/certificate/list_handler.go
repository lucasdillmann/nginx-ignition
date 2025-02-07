package certificate

import (
	"dillmann.com.br/nginx-ignition/api/common/pagination"
	"dillmann.com.br/nginx-ignition/core/certificate"
	"github.com/gin-gonic/gin"
	"net/http"
)

type listHandler struct {
	command *certificate.ListCommand
}

func (h listHandler) handle(context *gin.Context) {
	pageSize, pageNumber, searchTerms, err := pagination.ExtractPaginationParameters(context)
	if err != nil {
		panic(err)
	}

	page, err := (*h.command)(pageSize, pageNumber, searchTerms)
	if err != nil {
		panic(err)
	}

	context.JSON(http.StatusOK, pagination.Convert(page, toCertificateResponse))
}
