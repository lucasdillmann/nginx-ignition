package pagination

import (
	"dillmann.com.br/nginx-ignition/api/common/api_error"
	"dillmann.com.br/nginx-ignition/core/common/value_range"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

var (
	pageSizeRange = value_range.New(1, 1000)
)

func ExtractPaginationParameters(context *gin.Context) (int, int, *string, error) {
	pageSize := context.DefaultQuery("pageSize", "25")
	pageNumber := context.DefaultQuery("pageNumber", "0")
	searchTerms := context.Query("searchTerms")

	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		return 0, 0, nil, api_error.New(
			http.StatusBadRequest,
			"Page size must be an integer",
		)
	}

	if !pageSizeRange.Contains(pageSizeInt) {
		return 0, 0, nil, api_error.New(
			http.StatusBadRequest,
			"Page size must be between "+strconv.Itoa(pageSizeRange.Min)+" and "+strconv.Itoa(pageSizeRange.Max),
		)
	}

	pageNumberInt, err := strconv.Atoi(pageNumber)
	if err != nil {
		return 0, 0, nil, api_error.New(
			http.StatusBadRequest,
			"Page number must be an integer",
		)
	}

	if pageNumberInt < 0 {
		return 0, 0, nil, api_error.New(
			http.StatusBadRequest,
			"Page number must be greater than or equal to 0",
		)
	}

	searchTermsPtr := &searchTerms
	if strings.TrimSpace(searchTerms) == "" {
		searchTermsPtr = nil
	}

	return pageSizeInt, pageNumberInt, searchTermsPtr, nil
}
