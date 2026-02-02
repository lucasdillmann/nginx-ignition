package logline

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"dillmann.com.br/nginx-ignition/core/nginx"
)

func ExtractSearchParams(ctx *gin.Context) *nginx.LogSearch {
	searchQuery := ctx.Query("search")
	if searchQuery == "" {
		return nil
	}

	output := &nginx.LogSearch{
		Query: strings.TrimSpace(searchQuery),
	}

	surroundingLines := ctx.Query("surroundingLines")
	output.SurroundingLines, _ = strconv.Atoi(surroundingLines)

	if output.SurroundingLines < 0 {
		output.SurroundingLines = 0
	} else if output.SurroundingLines > 10 {
		output.SurroundingLines = 10
	}

	return output
}
