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

	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/core/settings"
)

func Test_ListHandler(t *testing.T) {
	page := pagination.New(1, 10, 1, []host.Host{
		{
			DomainNames: []string{"Test"},
		},
	})

	t.Run("Handle", func(t *testing.T) {
		t.Run("returns 200 OK with paginated results", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			commands := host.NewMockedCommands(ctrl)
			commands.EXPECT().
				List(gomock.Any(), 10, 1, gomock.Any()).
				Return(page, nil)

			settingsCommands := settings.NewMockedCommands(ctrl)
			settingsCommands.EXPECT().
				Get(gomock.Any()).
				Return(&settings.Settings{}, nil)

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = httptest.NewRequest("GET", "/api/hosts?pageSize=10&pageNumber=1", nil)

			handler := listHandler{
				hostCommands:     commands,
				settingsCommands: settingsCommands,
			}
			handler.handle(ctx)

			assert.Equal(t, http.StatusOK, w.Code)
			var response pagination.Page[hostResponseDTO]
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, 1, response.TotalItems)
			assert.Equal(t, "Test", response.Contents[0].DomainNames[0])
		})

		t.Run("panics when command returns error", func(t *testing.T) {
			expectedErr := errors.New("command error")
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			commands := host.NewMockedCommands(ctrl)
			commands.EXPECT().
				List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil, expectedErr)

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = httptest.NewRequest("GET", "/api/hosts", nil)

			handler := listHandler{
				hostCommands:     commands,
				settingsCommands: nil,
			}
			assert.PanicsWithValue(t, expectedErr, func() {
				handler.handle(ctx)
			})
		})
	})
}
