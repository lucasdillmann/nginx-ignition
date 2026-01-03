package vpn

import (
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

func Test_GetHandler_Handle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("returns 200 OK with vpn data on success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := uuid.New()
		mockVPN := &vpn.VPN{ID: id, Name: "vpn-1"}
		commands := vpn.NewMockedCommands(ctrl)
		commands.EXPECT().
			Get(gomock.Any(), id).
			Return(mockVPN, nil)

		handler := getHandler{commands}
		r := gin.New()
		r.GET("/api/vpns/:id", handler.handle)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/vpns/"+id.String(), nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp vpnResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, id, resp.ID)
	})

	t.Run("returns 404 Not Found on invalid ID", func(t *testing.T) {
		handler := getHandler{nil}
		r := gin.New()
		r.GET("/api/vpns/:id", handler.handle)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/vpns/invalid", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("returns 404 Not Found when vpn does not exist", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := uuid.New()
		commands := vpn.NewMockedCommands(ctrl)
		commands.EXPECT().
			Get(gomock.Any(), id).
			Return(nil, nil)

		handler := getHandler{commands}
		r := gin.New()
		r.GET("/api/vpns/:id", handler.handle)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/vpns/"+id.String(), nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("panics on command error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := uuid.New()
		expectedErr := assert.AnError
		commands := vpn.NewMockedCommands(ctrl)
		commands.EXPECT().
			Get(gomock.Any(), id).
			Return(nil, expectedErr)

		handler := getHandler{commands}
		r := gin.New()
		r.GET("/api/vpns/:id", func(c *gin.Context) {
			defer func() {
				if r := recover(); r != nil {
					assert.Equal(t, expectedErr, r)
					panic(r)
				}
			}()
			handler.handle(c)
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/vpns/"+id.String(), nil)

		assert.Panics(t, func() {
			r.ServeHTTP(w, req)
		})
	})
}
