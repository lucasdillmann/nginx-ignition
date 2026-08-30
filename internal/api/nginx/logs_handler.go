package nginx

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/lucasdillmann/nginx-ignition/internal/api/common/logline"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/valuerange"
	"github.com/lucasdillmann/nginx-ignition/internal/core/nginx"
)

type logsHandler struct {
	commands nginx.Commands
}

const (
	defaultLineCount = 50
)

var lineCountRange = valuerange.New(1, 99_999)

func (h logsHandler) handle(ctx *gin.Context) {
	lineCount := defaultLineCount
	queryValue := ctx.Query("lines")

	if queryValue != "" {
		var err error
		lineCount, err = strconv.Atoi(queryValue)

		if err != nil || !lineCountRange.Contains(lineCount) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": fmt.Sprintf(
					"Lines amount should be between %d and %d",
					lineCountRange.Min,
					lineCountRange.Max,
				),
			})
			return
		}
	}

	search := logline.ExtractSearchParams(ctx)

	logs, err := h.commands.GetMainLogs(ctx.Request.Context(), lineCount, search)
	if err != nil {
		panic(err)
	}

	payload := logline.ToResponseDTOs(logs)
	ctx.JSON(http.StatusOK, payload)
}
