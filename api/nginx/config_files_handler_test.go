package nginx

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/nginx"
)

func Test_ConfigFilesHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Handle", func(t *testing.T) {
		t.Run("returns 200 OK with zip data on success", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockZip := []byte("zip content")
			commands := nginx.NewMockedCommands(ctrl)
			commands.EXPECT().
				GetConfigFiles(gomock.Any(), gomock.Any()).
				Return(mockZip, nil)

			handler := configFilesHandler{
				commands: commands,
			}
			r := gin.New()
			r.GET("/api/nginx/config-files", handler.handle)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(
				"GET",
				"/api/nginx/config-files?basePath=/&configPath=/etc/nginx/",
				nil,
			)
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, "application/zip", w.Header().Get("Content-Type"))
			assert.Equal(t, mockZip, w.Body.Bytes())
		})

		t.Run("panics when command returns error", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			expectedErr := assert.AnError
			commands := nginx.NewMockedCommands(ctrl)
			commands.EXPECT().
				GetConfigFiles(gomock.Any(), gomock.Any()).
				Return(nil, expectedErr)

			handler := configFilesHandler{
				commands: commands,
			}
			r := gin.New()
			r.GET("/api/nginx/config-files", func(c *gin.Context) {
				defer func() {
					if r := recover(); r != nil {
						assert.Equal(t, expectedErr, r)
						panic(r)
					}
				}()
				handler.handle(c)
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/nginx/config-files", nil)

			assert.Panics(t, func() {
				r.ServeHTTP(w, req)
			})
		})
	})
}
