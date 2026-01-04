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
	"dillmann.com.br/nginx-ignition/core/common/ptr"
	"dillmann.com.br/nginx-ignition/core/user"
)

func Test_CreateHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Handle", func(t *testing.T) {
		t.Run("returns 201 Created on success", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			payload := &userRequestDTO{
				Name:     ptr.Of("John Doe"),
				Username: ptr.Of("johndoe"),
				Enabled:  ptr.Of(true),
			}

			commands := user.NewMockedCommands(ctrl)
			commands.EXPECT().
				Save(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil)

			handler := createHandler{
				commands: commands,
			}
			r := gin.New()
			r.Use(func(c *gin.Context) {
				c.Set("ABAC:Subject", &authorization.Subject{User: &user.User{ID: uuid.New()}})
				c.Next()
			})
			r.POST("/api/users", handler.handle)

			jsonPayload, _ := json.Marshal(payload)
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/users", bytes.NewBuffer(jsonPayload))
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusCreated, w.Code)
			var resp map[string]uuid.UUID
			json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NotEqual(t, uuid.Nil, resp["id"])
		})

		t.Run("panics on command error", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			payload := &userRequestDTO{
				Name:     ptr.Of("John Doe"),
				Username: ptr.Of("johndoe"),
				Enabled:  ptr.Of(true),
			}

			expectedErr := assert.AnError
			commands := user.NewMockedCommands(ctrl)
			commands.EXPECT().
				Save(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(expectedErr)

			handler := createHandler{
				commands: commands,
			}
			r := gin.New()
			r.Use(func(c *gin.Context) {
				c.Set("ABAC:Subject", &authorization.Subject{User: &user.User{ID: uuid.New()}})
				c.Next()
			})
			r.POST("/api/users", func(c *gin.Context) {
				defer func() {
					if r := recover(); r != nil {
						assert.Equal(t, expectedErr, r)
						panic(r)
					}
				}()
				handler.handle(c)
			})

			jsonPayload, _ := json.Marshal(payload)
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/users", bytes.NewBuffer(jsonPayload))

			assert.Panics(t, func() {
				r.ServeHTTP(w, req)
			})
		})
	})
}
