package backup

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/backup"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func Test_getHandler(t *testing.T) {
	t.Run("handle", func(t *testing.T) {
		t.Run("returns 200 OK with backup data on success", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			backupData := newBackup()
			commands := backup.NewMockedCommands(controller)
			commands.EXPECT().
				Get(gomock.Any()).
				Return(backupData, nil)

			recorder := httptest.NewRecorder()
			ginContext, _ := gin.CreateTestContext(recorder)
			ginContext.Request = httptest.NewRequest("GET", "/api/backup", nil)

			handler := getHandler{
				commands: commands,
			}
			handler.handle(ginContext)

			assert.Equal(t, http.StatusOK, recorder.Code)
			assert.Equal(
				t,
				"attachment; filename=backup.zip",
				recorder.Header().Get("Content-Disposition"),
			)
			assert.Equal(t, "application/zip", recorder.Header().Get("Content-Type"))
			assert.Equal(t, backupData.Contents, recorder.Body.Bytes())
		})

		t.Run("panics when command returns error", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			expectedErr := errors.New("backup error")
			commands := backup.NewMockedCommands(controller)
			commands.EXPECT().
				Get(gomock.Any()).
				Return(nil, expectedErr)

			recorder := httptest.NewRecorder()
			ginContext, _ := gin.CreateTestContext(recorder)
			ginContext.Request = httptest.NewRequest("GET", "/api/backup", nil)

			handler := getHandler{
				commands: commands,
			}
			assert.PanicsWithValue(t, expectedErr, func() {
				handler.handle(ginContext)
			})
		})
	})
}
