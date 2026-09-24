package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lucasdillmann/nginx-ignition/internal/core/common/configuration"
)

func bodyLimitMiddleware(cfg *configuration.Configuration) (gin.HandlerFunc, error) {
	maxBytes, err := cfg.GetInt("nginx-ignition.server.max-body-bytes")
	if err != nil {
		return nil, err
	}

	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, int64(maxBytes))
		c.Next()
	}, nil
}
