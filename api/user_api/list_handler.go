package user_api

import (
	"dillmann.com.br/nginx-ignition/api/common/pagination"
	"dillmann.com.br/nginx-ignition/core/user"
	"github.com/gin-gonic/gin"
	"net/http"
)

type listHandler struct {
	command *user.ListCommand
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

	context.JSON(http.StatusOK, pagination.Convert(page, toDto))
}
