package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/api/common/pagination"
	corepagination "dillmann.com.br/nginx-ignition/core/common/pagination"
	"dillmann.com.br/nginx-ignition/core/integration"
)

func Test_ListOptionsHandler_Handle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("returns 200 OK with options list on success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := uuid.New()
		mockOptions := []integration.DriverOption{
			{ID: "opt-1", Name: "Option 1"},
			{ID: "opt-2", Name: "Option 2"},
		}
		page := &corepagination.Page[integration.DriverOption]{
			Contents: mockOptions,
		}

		commands := integration.NewMockedCommands(ctrl)
		commands.EXPECT().
			ListOptions(gomock.Any(), id, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(page, nil)

		handler := listOptionsHandler{commands}
		r := gin.New()
		r.GET("/api/integrations/:id/options", handler.handle)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/integrations/"+id.String()+"/options", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp pagination.PageDTO[integrationOptionResponse]
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Len(t, resp.Contents, 2)
	})

	t.Run("returns 404 Not Found on invalid ID", func(t *testing.T) {
		handler := listOptionsHandler{nil}
		r := gin.New()
		r.GET("/api/integrations/:id/options", handler.handle)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/integrations/invalid/options", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("panics on command error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := uuid.New()
		expectedErr := assert.AnError
		commands := integration.NewMockedCommands(ctrl)
		commands.EXPECT().
			ListOptions(gomock.Any(), id, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, expectedErr)

		handler := listOptionsHandler{commands}
		r := gin.New()
		r.GET("/api/integrations/:id/options", func(c *gin.Context) {
			defer func() {
				if r := recover(); r != nil {
					assert.Equal(t, expectedErr, r)
					panic(r)
				}
			}()
			handler.handle(c)
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/integrations/"+id.String()+"/options", nil)

		assert.Panics(t, func() {
			r.ServeHTTP(w, req)
		})
	})
}
