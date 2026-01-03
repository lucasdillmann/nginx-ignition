package vpn

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/api/common/pagination"
	corepagination "dillmann.com.br/nginx-ignition/core/common/pagination"
	"dillmann.com.br/nginx-ignition/core/vpn"
)

func Test_ListHandler_Handle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("returns 200 OK with vpn list on success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockVPNs := []vpn.VPN{
			{Name: "vpn-1"},
			{Name: "vpn-2"},
		}
		page := &corepagination.Page[vpn.VPN]{
			Contents: mockVPNs,
		}

		commands := vpn.NewMockedCommands(ctrl)
		commands.EXPECT().
			List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(page, nil)

		handler := listHandler{commands}
		r := gin.New()
		r.GET("/api/vpns", handler.handle)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/vpns?pageSize=10&pageNumber=1", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp pagination.PageDTO[vpnResponse]
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Len(t, resp.Contents, 2)
	})

	t.Run("panics on command error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		expectedErr := assert.AnError
		commands := vpn.NewMockedCommands(ctrl)
		commands.EXPECT().
			List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, expectedErr)

		handler := listHandler{commands}
		r := gin.New()
		r.GET("/api/vpns", func(c *gin.Context) {
			defer func() {
				if r := recover(); r != nil {
					assert.Equal(t, expectedErr, r)
					panic(r)
				}
			}()
			handler.handle(c)
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/vpns", nil)

		assert.Panics(t, func() {
			r.ServeHTTP(w, req)
		})
	})
}
