package host

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/common/valuerange"
	"dillmann.com.br/nginx-ignition/core/nginx"
)

type logsHandler struct {
	commands *nginx.Commands
}

const (
	defaultLineCount = 50
)

var (
	lineCountRange    = valuerange.New(1, 10_000)
	allowedQualifiers = map[string]bool{
		"access": true,
		"error":  true,
	}
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

	qualifier := ctx.Param("qualifier")
	if !allowedQualifiers[qualifier] {
		ctx.Status(http.StatusNotFound)
		return
	}

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	logs, err := h.commands.GetHostLogs(ctx.Request.Context(), id, qualifier, lineCount)
	if err != nil {
		panic(err)
	}

	ctx.JSON(http.StatusOK, logs)
}
