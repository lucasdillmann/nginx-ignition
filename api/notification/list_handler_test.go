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
	"dillmann.com.br/nginx-ignition/api/common/pagination"
	corepagination "dillmann.com.br/nginx-ignition/core/common/pagination"
	"dillmann.com.br/nginx-ignition/core/notification"
	"dillmann.com.br/nginx-ignition/core/user"
)

func Test_listHandler(t *testing.T) {
	t.Run("handle", func(t *testing.T) {
		t.Run("returns 200 OK with notification list on success", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			userID := uuid.New()
			item := sampleNotification(userID)
			page := corepagination.New(1, 10, 1, []notification.Notification{*item})
			commands := notification.NewMockedCommands(controller)
			commands.EXPECT().
				ListNotifications(gomock.Any(), userID, gomock.Any(), gomock.Any(), gomock.Any()).
				Return(page, nil)

			handler := listHandler{commands: commands}
			engine := gin.New()
			engine.Use(func(ginContext *gin.Context) {
				ginContext.Set("ABAC:Subject", &authorization.Subject{User: &user.User{ID: userID}})
				ginContext.Next()
			})
			engine.GET("", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, "/?pageSize=10&pageNumber=1", nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusOK, recorder.Code)
			var response pagination.PageDTO[notificationResponse]
			json.Unmarshal(recorder.Body.Bytes(), &response)
			assert.Len(t, response.Contents, 1)
			assert.Equal(t, item.ID, response.Contents[0].ID)
		})

		t.Run("passes search terms to command", func(t *testing.T) {
			searchTerm := new("certificate")
			controller := gomock.NewController(t)
			defer controller.Finish()

			userID := uuid.New()
			item := sampleNotification(userID)
			page := corepagination.New(1, 10, 1, []notification.Notification{*item})
			commands := notification.NewMockedCommands(controller)
			commands.EXPECT().
				ListNotifications(gomock.Any(), userID, gomock.Any(), gomock.Any(), gomock.Eq(searchTerm)).
				Return(page, nil)

			handler := listHandler{commands: commands}
			engine := gin.New()
			engine.Use(func(ginContext *gin.Context) {
				ginContext.Set("ABAC:Subject", &authorization.Subject{User: &user.User{ID: userID}})
				ginContext.Next()
			})
			engine.GET("", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(
				http.MethodGet,
				"/?searchTerms="+*searchTerm+"&pageSize=10&pageNumber=1",
				nil,
			)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusOK, recorder.Code)
		})

		t.Run("returns 401 Unauthorized without subject", func(t *testing.T) {
			handler := listHandler{commands: nil}
			engine := gin.New()
			engine.GET("", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, "/", nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusUnauthorized, recorder.Code)
		})
	})
}
