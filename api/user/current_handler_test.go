package user

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/api/common/authorization"
	"dillmann.com.br/nginx-ignition/core/user"
)

func Test_CurrentHandler_Handle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("returns 200 OK with current user data", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := uuid.New()
		mockUser := &user.User{ID: id, Name: "Current User"}

		handler := currentHandler{}
		r := gin.New()
		r.Use(func(c *gin.Context) {
			c.Set("ABAC:Subject", &authorization.Subject{User: mockUser})
			c.Next()
		})
		r.GET("/api/users/current", handler.handle)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/users/current", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp userResponseDTO
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, id, resp.ID)
	})

	t.Run("returns 401 Unauthorized when subject is missing", func(t *testing.T) {
		handler := currentHandler{}
		r := gin.New()
		r.GET("/api/users/current", handler.handle)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/users/current", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
