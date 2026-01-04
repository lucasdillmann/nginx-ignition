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

func Test_GetHandler(t *testing.T) {
	t.Run("Handle", func(t *testing.T) {
		t.Run("returns 200 OK with backup data on success", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockData := &backup.Backup{
				FileName:    "backup.zip",
				ContentType: "application/zip",
				Contents:    []byte("backup data"),
			}

			commands := backup.NewMockedCommands(ctrl)
			commands.EXPECT().
				Get(gomock.Any()).
				Return(mockData, nil)

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = httptest.NewRequest("GET", "/api/backup", nil)

			handler := getHandler{
				commands: commands,
			}
			handler.handle(ctx)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(
				t,
				"attachment; filename=backup.zip",
				w.Header().Get("Content-Disposition"),
			)
			assert.Equal(t, "application/zip", w.Header().Get("Content-Type"))
			assert.Equal(t, mockData.Contents, w.Body.Bytes())
		})

		t.Run("panics when command returns error", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			expectedErr := errors.New("backup error")
			commands := backup.NewMockedCommands(ctrl)
			commands.EXPECT().
				Get(gomock.Any()).
				Return(nil, expectedErr)

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = httptest.NewRequest("GET", "/api/backup", nil)

			handler := getHandler{
				commands: commands,
			}
			assert.PanicsWithValue(t, expectedErr, func() {
				handler.handle(ctx)
			})
		})
	})
}
