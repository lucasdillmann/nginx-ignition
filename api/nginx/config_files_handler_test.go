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

func init() {
	gin.SetMode(gin.TestMode)
}

func Test_configFilesHandler(t *testing.T) {
	t.Run("handle", func(t *testing.T) {
		t.Run("returns 200 OK with zip data on success", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			zipData := []byte("zip content")
			commands := nginx.NewMockedCommands(controller)
			commands.EXPECT().
				GetConfigFiles(gomock.Any(), gomock.Any()).
				Return(zipData, nil)

			handler := configFilesHandler{
				commands: commands,
			}
			engine := gin.New()
			engine.GET("/api/nginx/config-files", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(
				"GET",
				"/api/nginx/config-files?basePath=/&configPath=/etc/nginx/",
				nil,
			)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusOK, recorder.Code)
			assert.Equal(t, "application/zip", recorder.Header().Get("Content-Type"))
			assert.Equal(t, zipData, recorder.Body.Bytes())
		})

		t.Run("panics when command returns error", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			expectedErr := assert.AnError
			commands := nginx.NewMockedCommands(controller)
			commands.EXPECT().
				GetConfigFiles(gomock.Any(), gomock.Any()).
				Return(nil, expectedErr)

			handler := configFilesHandler{
				commands: commands,
			}
			engine := gin.New()
			engine.GET("/api/nginx/config-files", func(ginContext *gin.Context) {
				defer func() {
					if r := recover(); r != nil {
						assert.Equal(t, expectedErr, r)
						panic(r)
					}
				}()
				handler.handle(ginContext)
			})

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/api/nginx/config-files", nil)

			assert.Panics(t, func() {
				engine.ServeHTTP(recorder, request)
			})
		})
	})
}
