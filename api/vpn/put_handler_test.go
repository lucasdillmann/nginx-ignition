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

func Test_PutHandler_Handle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("returns 204 No Content on success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := uuid.New()
		payload := &vpnRequest{
			Name: "updated-vpn",
		}

		commands := vpn.NewMockedCommands(ctrl)
		commands.EXPECT().
			Save(gomock.Any(), gomock.Any()).
			Return(nil)

		handler := putHandler{commands}
		r := gin.New()
		r.PUT("/api/vpns/:id", handler.handle)

		jsonPayload, _ := json.Marshal(payload)
		w := httptest.NewRecorder()
		req := httptest.NewRequest("PUT", "/api/vpns/"+id.String(), bytes.NewBuffer(jsonPayload))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("returns 404 Not Found on invalid ID", func(t *testing.T) {
		handler := putHandler{nil}
		r := gin.New()
		r.PUT("/api/vpns/:id", handler.handle)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("PUT", "/api/vpns/invalid", bytes.NewBufferString("{}"))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("panics on command error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := uuid.New()
		payload := &vpnRequest{
			Name: "updated-vpn",
		}

		expectedErr := assert.AnError
		commands := vpn.NewMockedCommands(ctrl)
		commands.EXPECT().
			Save(gomock.Any(), gomock.Any()).
			Return(expectedErr)

		handler := putHandler{commands}
		r := gin.New()
		r.PUT("/api/vpns/:id", func(c *gin.Context) {
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
		req := httptest.NewRequest("PUT", "/api/vpns/"+id.String(), bytes.NewBuffer(jsonPayload))

		assert.Panics(t, func() {
			r.ServeHTTP(w, req)
		})
	})
}
