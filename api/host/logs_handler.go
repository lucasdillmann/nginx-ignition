package host

import (
	"dillmann.com.br/nginx-ignition/core/common/value_range"
	"dillmann.com.br/nginx-ignition/core/nginx"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strconv"
)

type logsHandler struct {
	command *nginx.GetHostLogsCommand
}

const (
	defaultLineCount = 50
)

var (
	lineCountRange    = value_range.New(1, 10_000)
	allowedQualifiers = map[string]bool{
		"access": true,
		"error":  true,
	}
)

func (h logsHandler) handle(context *gin.Context) {
	lineCount := defaultLineCount
	queryValue := context.Query("lines")

	if queryValue != "" {
		var err error
		lineCount, err = strconv.Atoi(queryValue)

		if err != nil || !lineCountRange.Contains(lineCount) {
			context.JSON(http.StatusBadRequest, gin.H{
				"message": fmt.Sprintf(
					"Lines amount should be between %d and %d",
					lineCountRange.Min,
					lineCountRange.Max,
				),
			})
			return
		}
	}

	qualifier := context.Param("qualifier")
	if !allowedQualifiers[qualifier] {
		context.Status(http.StatusNotFound)
		return
	}

	id, err := uuid.Parse(context.Param("id"))
	if err != nil {
		context.Status(http.StatusNotFound)
		return
	}

	logs, err := (*h.command)(id, qualifier, lineCount)
	if err != nil {
		panic(err)
	}

	context.JSON(http.StatusOK, logs)
}
