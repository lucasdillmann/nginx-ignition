package user

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/api/common/authorization"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func Test_currentHandler(t *testing.T) {
	t.Run("handle", func(t *testing.T) {
		t.Run("returns 200 OK with current user data", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			subject := newUser()
			handler := currentHandler{}
			engine := gin.New()
			engine.Use(func(ginContext *gin.Context) {
				ginContext.Set("ABAC:Subject", &authorization.Subject{User: subject})
				ginContext.Next()
			})
			engine.GET("/api/users/current", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/api/users/current", nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusOK, recorder.Code)
			var response userResponseDTO
			json.Unmarshal(recorder.Body.Bytes(), &response)
			assert.Equal(t, subject.ID, response.ID)
		})

		t.Run("returns 401 Unauthorized when subject is missing", func(t *testing.T) {
			handler := currentHandler{}
			engine := gin.New()
			engine.GET("/api/users/current", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/api/users/current", nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusUnauthorized, recorder.Code)
		})
	})
}
