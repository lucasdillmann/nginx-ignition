package nginx

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"dillmann.com.br/nginx-ignition/core/common/value_range"
	"dillmann.com.br/nginx-ignition/core/nginx"
)

type logsHandler struct {
	commands *nginx.Commands
}

const (
	defaultLineCount = 50
)

var (
	lineCountRange = value_range.New(1, 10_000)
)

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

	logs, err := h.commands.GetMainLogs(ctx.Request.Context(), lineCount)
	if err != nil {
		panic(err)
	}

	ctx.JSON(http.StatusOK, logs)
}
