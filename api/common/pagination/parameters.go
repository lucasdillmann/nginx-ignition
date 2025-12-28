package pagination

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"dillmann.com.br/nginx-ignition/api/common/apierror"
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
			"Page size must be an integer",
		)
	}

	if !pageSizeRange.Contains(pageSize) {
		return 0, 0, nil, apierror.New(
			http.StatusBadRequest,
			fmt.Sprintf("Page size must be between %d and %d", pageSizeRange.Min, pageSizeRange.Max),
		)
	}

	pageNumber, err = strconv.Atoi(pageNumberStr)
	if err != nil {
		return 0, 0, nil, apierror.New(
			http.StatusBadRequest,
			"Page number must be an integer",
		)
	}

	if pageNumber < 0 {
		return 0, 0, nil, apierror.New(
			http.StatusBadRequest,
			"Page number must be greater than or equal to 0",
		)
	}

	searchTermsPtr := &searchTermsStr
	if strings.TrimSpace(searchTermsStr) == "" {
		searchTermsPtr = nil
	}

	return pageSize, pageNumber, searchTermsPtr, nil
}
