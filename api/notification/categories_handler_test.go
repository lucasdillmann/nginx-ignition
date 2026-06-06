package notification

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
	"dillmann.com.br/nginx-ignition/core/notification"
	"dillmann.com.br/nginx-ignition/core/user"
)

func Test_categoriesHandler(t *testing.T) {
	t.Run("handle", func(t *testing.T) {
		t.Run("returns 200 OK with categories on success", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			userID := uuid.New()
			categories := []notification.CategoryInfo{
				{ID: notification.CategoryCertificateRenewed},
			}
			commands := notification.NewMockedCommands(controller)
			commands.EXPECT().
				ListCategories(gomock.Any()).
				Return(categories, nil)

			handler := categoriesHandler{commands: commands}
			engine := gin.New()
			engine.Use(func(ginContext *gin.Context) {
				ginContext.Set("ABAC:Subject", &authorization.Subject{User: &user.User{ID: userID}})
				ginContext.Next()
			})
			engine.GET("/categories", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, "/categories", nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusOK, recorder.Code)
			var response []categoryResponse
			json.Unmarshal(recorder.Body.Bytes(), &response)
			assert.Len(t, response, 1)
			assert.Equal(t, string(notification.CategoryCertificateRenewed), response[0].ID)
		})

		t.Run("returns 401 Unauthorized without subject", func(t *testing.T) {
			handler := categoriesHandler{commands: nil}
			engine := gin.New()
			engine.GET("/categories", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, "/categories", nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusUnauthorized, recorder.Code)
		})
	})
}
