package healthcheck

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type readinessHandler struct{}

func (h readinessHandler) handle(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"ready":   true,
		"message": "I'm alive (but the cake is a lie 👀🍰)",
	})
}
