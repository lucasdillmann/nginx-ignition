package pagination

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"dillmann.com.br/nginx-ignition/api/common/apierror"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
	"dillmann.com.br/nginx-ignition/core/common/valuerange"
)

var pageSizeRange = valuerange.New(1, 1000)

func ExtractPaginationParameters(ctx *gin.Context) (
	pageSize int,
	pageNumber int,
	searchTerms *string,
	err error,
) {
	pageSizeStr := ctx.DefaultQuery("pageSize", "25")
	pageNumberStr := ctx.DefaultQuery("pageNumber", "0")
	searchTermsStr := ctx.Query("searchTerms")

	pageSize, err = strconv.Atoi(pageSizeStr)
	if err != nil {
		return 0, 0, nil, apierror.New(
			http.StatusBadRequest,
			i18n.M(ctx.Request.Context(), i18n.K.PaginationErrorMustBeAnInteger).V("type", "size"),
		)
	}

	if !pageSizeRange.Contains(pageSize) {
		return 0, 0, nil, apierror.New(
			http.StatusBadRequest,
			i18n.M(ctx.Request.Context(), i18n.K.PaginationErrorMustBeBetweenRange).
				V("type", "size").
				V("min", pageSizeRange.Min).
				V("max", pageSizeRange.Max),
		)
	}

	pageNumber, err = strconv.Atoi(pageNumberStr)
	if err != nil {
		return 0, 0, nil, apierror.New(
			http.StatusBadRequest,
			i18n.M(ctx.Request.Context(), i18n.K.PaginationErrorMustBeAnInteger).
				V("type", "number"),
		)
	}

	if pageNumber < 0 {
		return 0, 0, nil, apierror.New(
			http.StatusBadRequest,
			i18n.M(ctx.Request.Context(), i18n.K.PaginationErrorCantBeNegative).
				V("type", "number"),
		)
	}

	searchTermsPtr := &searchTermsStr
	if strings.TrimSpace(searchTermsStr) == "" {
		searchTermsPtr = nil
	}

	return pageSize, pageNumber, searchTermsPtr, nil
}
