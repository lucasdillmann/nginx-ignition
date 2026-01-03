package host

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/core/settings"
)

func Test_GetHandler_Handle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("returns 200 OK with host data on success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := uuid.New()
		mockHost := &host.Host{ID: id, DomainNames: []string{"example.com"}}
		hostCommands := host.NewMockedCommands(ctrl)
		hostCommands.EXPECT().
			Get(gomock.Any(), id).
			Return(mockHost, nil)

		settingsCommands := settings.NewMockedCommands(ctrl)
		settingsCommands.EXPECT().
			Get(gomock.Any()).
			Return(nil, nil)

		handler := getHandler{settingsCommands, hostCommands}
		r := gin.New()
		r.GET("/api/hosts/:id", handler.handle)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/hosts/"+id.String(), nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp hostResponseDTO
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, id, *resp.ID)
	})

	t.Run("returns 404 Not Found on invalid ID", func(t *testing.T) {
		handler := getHandler{nil, nil}
		r := gin.New()
		r.GET("/api/hosts/:id", handler.handle)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/hosts/invalid", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("returns 404 Not Found when host does not exist", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := uuid.New()
		hostCommands := host.NewMockedCommands(ctrl)
		hostCommands.EXPECT().
			Get(gomock.Any(), id).
			Return(nil, nil)

		handler := getHandler{nil, hostCommands}
		r := gin.New()
		r.GET("/api/hosts/:id", handler.handle)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/hosts/"+id.String(), nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("panics when command returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := uuid.New()
		expectedErr := errors.New("get error")
		hostCommands := host.NewMockedCommands(ctrl)
		hostCommands.EXPECT().
			Get(gomock.Any(), id).
			Return(nil, expectedErr)

		handler := getHandler{nil, hostCommands}
		r := gin.New()
		r.GET("/api/hosts/:id", func(c *gin.Context) {
			defer func() {
				if r := recover(); r != nil {
					assert.Equal(t, expectedErr, r)
					panic(r)
				}
			}()
			handler.handle(c)
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/hosts/"+id.String(), nil)

		assert.Panics(t, func() {
			r.ServeHTTP(w, req)
		})
	})
}
