package user

import (
	"bytes"
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

func init() {
	gin.SetMode(gin.TestMode)
}

func Test_totpActivateHandler(t *testing.T) {
	t.Run("handle", func(t *testing.T) {
		t.Run("returns 204 No Content on valid code", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			id := uuid.New()
			payload := totpActivateRequestDTO{
				Code: new("123456"),
			}

			commands := user.NewMockedCommands(controller)
			commands.EXPECT().
				ActivateTOTP(gomock.Any(), id, *payload.Code).
				Return(true, nil)

			handler := totpActivateHandler{commands: commands}
			engine := gin.New()
			engine.Use(func(ginContext *gin.Context) {
				ginContext.Set("ABAC:Subject", &authorization.Subject{User: &user.User{ID: id}})
				ginContext.Next()
			})
			engine.POST("/current/totp/activate", handler.handle)

			body, _ := json.Marshal(payload)
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("POST", "/current/totp/activate", bytes.NewBuffer(body))
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusNoContent, recorder.Code)
		})

		t.Run("returns 400 Bad Request on invalid code", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			id := uuid.New()
			payload := totpActivateRequestDTO{
				Code: new("000000"),
			}

			commands := user.NewMockedCommands(controller)
			commands.EXPECT().
				ActivateTOTP(gomock.Any(), id, *payload.Code).
				Return(false, nil)

			handler := totpActivateHandler{commands: commands}
			engine := gin.New()
			engine.Use(func(ginContext *gin.Context) {
				ginContext.Set("ABAC:Subject", &authorization.Subject{User: &user.User{ID: id}})
				ginContext.Next()
			})
			engine.POST("/current/totp/activate", handler.handle)

			body, _ := json.Marshal(payload)
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("POST", "/current/totp/activate", bytes.NewBuffer(body))
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusBadRequest, recorder.Code)
		})
	})
}
