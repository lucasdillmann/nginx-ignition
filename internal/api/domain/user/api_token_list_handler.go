package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lucasdillmann/nginx-ignition/internal/api/core/authorization"
	"github.com/lucasdillmann/nginx-ignition/internal/api/core/pagination"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/user"
)

type apiTokenListHandler struct {
	commands user.Commands
}

func (h apiTokenListHandler) handle(ctx *gin.Context) {
	currentSubject := authorization.CurrentSubject(ctx)
	if currentSubject == nil || currentSubject.User == nil {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	pageSize, pageNumber, searchTerms, err := pagination.ExtractPaginationParameters(ctx)
	if err != nil {
		panic(err)
	}

	page, err := h.commands.ListAPITokens(
		ctx.Request.Context(),
		currentSubject.User.ID,
		pageSize,
		pageNumber,
		searchTerms,
	)
	if err != nil {
		panic(err)
	}

	ctx.JSON(http.StatusOK, pagination.Convert(page, toAPITokenDTO))
}
