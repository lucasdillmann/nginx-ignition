package server

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/lucasdillmann/nginx-ignition/internal/api/common/apierror"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/configuration"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func Test_bodyLimitMiddleware(t *testing.T) {
	setup := func(maxBodyBytes int) *gin.Engine {
		cfg := configuration.NewWithOverrides(map[string]string{
			"nginx-ignition.server.max-body-bytes": strconv.Itoa(maxBodyBytes),
		})

		middleware, err := bodyLimitMiddleware(cfg)
		if err != nil {
			panic(err)
		}

		engine := gin.New()
		engine.Use(gin.CustomRecoveryWithWriter(nil, apierror.Handler))
		engine.Use(middleware)

		engine.POST("/test", func(ctx *gin.Context) {
			if _, err := io.ReadAll(ctx.Request.Body); err != nil {
				panic(err)
			}

			ctx.Status(http.StatusOK)
		})

		return engine
	}

	t.Run("allows request within body limit", func(t *testing.T) {
		engine := setup(1024)
		body := bytes.Repeat([]byte("x"), 100)

		recorder := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/test", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		engine.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
	})

	t.Run("rejects request exceeding body limit", func(t *testing.T) {
		engine := setup(1024)
		body := bytes.Repeat([]byte("x"), 2048)

		recorder := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/test", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		engine.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusRequestEntityTooLarge, recorder.Code)
	})

	t.Run("handles exactly at limit", func(t *testing.T) {
		engine := setup(1024)
		body := bytes.Repeat([]byte("x"), 1024)

		recorder := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/test", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		engine.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
	})

	t.Run("handles empty body", func(t *testing.T) {
		engine := setup(1024)

		recorder := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/test", nil)
		engine.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
	})
}
