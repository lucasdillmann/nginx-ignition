package vpn

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

	"dillmann.com.br/nginx-ignition/core/vpn"
)

func Test_CreateHandler_Handle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("returns 201 Created on success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		payload := &vpnRequest{
			Name:   "vpn-1",
			Driver: "openvpn",
		}

		commands := vpn.NewMockedCommands(ctrl)
		commands.EXPECT().
			Save(gomock.Any(), gomock.Any()).
			Return(nil)

		handler := createHandler{commands}
		r := gin.New()
		r.POST("/api/vpns", handler.handle)

		jsonPayload, _ := json.Marshal(payload)
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/vpns", bytes.NewBuffer(jsonPayload))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var resp map[string]uuid.UUID
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NotEqual(t, uuid.Nil, resp["id"])
	})

	t.Run("panics on command error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		payload := &vpnRequest{
			Name: "vpn-1",
		}

		expectedErr := assert.AnError
		commands := vpn.NewMockedCommands(ctrl)
		commands.EXPECT().
			Save(gomock.Any(), gomock.Any()).
			Return(expectedErr)

		handler := createHandler{commands}
		r := gin.New()
		r.POST("/api/vpns", func(c *gin.Context) {
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
		req := httptest.NewRequest("POST", "/api/vpns", bytes.NewBuffer(jsonPayload))

		assert.Panics(t, func() {
			r.ServeHTTP(w, req)
		})
	})
}
