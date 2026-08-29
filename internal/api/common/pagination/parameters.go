package pagination

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/lucasdillmann/nginx-ignition/internal/api/common/apierror"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/i18n"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/valuerange"
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
			i18n.M(ctx.Request.Context(), i18n.K.ApiCommonPaginationMustBeAnInteger).
				V("type", "size"),
		)
	}

	if !pageSizeRange.Contains(pageSize) {
		return 0, 0, nil, apierror.New(
			http.StatusBadRequest,
			i18n.M(ctx.Request.Context(), i18n.K.ApiCommonPaginationMustBeBetweenRange).
				V("type", "size").
				V("min", pageSizeRange.Min).
				V("max", pageSizeRange.Max),
		)
	}

	pageNumber, err = strconv.Atoi(pageNumberStr)
	if err != nil {
		return 0, 0, nil, apierror.New(
			http.StatusBadRequest,
			i18n.M(ctx.Request.Context(), i18n.K.ApiCommonPaginationMustBeAnInteger).
				V("type", "number"),
		)
	}

	if pageNumber < 0 {
		return 0, 0, nil, apierror.New(
			http.StatusBadRequest,
			i18n.M(ctx.Request.Context(), i18n.K.ApiCommonPaginationCantBeNegative).
				V("type", "number"),
		)
	}

	searchTermsPtr := &searchTermsStr
	if strings.TrimSpace(searchTermsStr) == "" {
		searchTermsPtr = nil
	}

	return pageSize, pageNumber, searchTermsPtr, nil
}
