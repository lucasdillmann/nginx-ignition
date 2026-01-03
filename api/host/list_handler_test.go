package host

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/api/common/pagination"
	corepagination "dillmann.com.br/nginx-ignition/core/common/pagination"
	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/core/settings"
)

func Test_ListHandler_Handle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("returns 200 OK with host list on success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockHosts := []host.Host{
			{DomainNames: []string{"example1.com"}},
			{DomainNames: []string{"example2.com"}},
		}

		page := &corepagination.Page[host.Host]{
			Contents: mockHosts,
		}

		hostCommands := host.NewMockedCommands(ctrl)
		hostCommands.EXPECT().
			List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(page, nil)

		settingsCommands := settings.NewMockedCommands(ctrl)
		settingsCommands.EXPECT().
			Get(gomock.Any()).
			Return(nil, nil)

		handler := listHandler{settingsCommands, hostCommands}
		r := gin.New()
		r.GET("/api/hosts", handler.handle)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/hosts?pageSize=10&pageNumber=1", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp pagination.PageDTO[hostResponseDTO]
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Len(t, resp.Contents, 2)
	})

	t.Run("panics when command returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		expectedErr := errors.New("list error")
		hostCommands := host.NewMockedCommands(ctrl)
		hostCommands.EXPECT().
			List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, expectedErr)

		handler := listHandler{nil, hostCommands}
		r := gin.New()
		r.GET("/api/hosts", func(c *gin.Context) {
			defer func() {
				if r := recover(); r != nil {
					assert.Equal(t, expectedErr, r)
					panic(r)
				}
			}()
			handler.handle(c)
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/hosts", nil)

		assert.Panics(t, func() {
			r.ServeHTTP(w, req)
		})
	})
}
