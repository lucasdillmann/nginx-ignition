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

func init() {
	gin.SetMode(gin.TestMode)
}

func Test_listHandler(t *testing.T) {
	t.Run("handle", func(t *testing.T) {
		t.Run("returns 200 OK with paginated results", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			page := newHostPage()
			commands := host.NewMockedCommands(controller)
			commands.EXPECT().
				List(gomock.Any(), 10, 1, gomock.Any()).
				Return(page, nil)

			settingsCommands := settings.NewMockedCommands(controller)
			settingsCommands.EXPECT().
				Get(gomock.Any()).
				Return(&settings.Settings{}, nil)

			recorder := httptest.NewRecorder()
			ginContext, _ := gin.CreateTestContext(recorder)
			ginContext.Request = httptest.NewRequest(
				"GET",
				"/api/hosts?pageSize=10&pageNumber=1",
				nil,
			)

			handler := listHandler{
				hostCommands:     commands,
				settingsCommands: settingsCommands,
			}
			handler.handle(ginContext)

			assert.Equal(t, http.StatusOK, recorder.Code)
			var response pagination.Page[hostResponseDTO]
			err := json.Unmarshal(recorder.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, 1, response.TotalItems)
			assert.Equal(t, page.Contents[0].DomainNames[0], response.Contents[0].DomainNames[0])
		})

		t.Run("panics when command returns error", func(t *testing.T) {
			expectedErr := errors.New("command error")
			controller := gomock.NewController(t)
			defer controller.Finish()

			commands := host.NewMockedCommands(controller)
			commands.EXPECT().
				List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil, expectedErr)

			recorder := httptest.NewRecorder()
			ginContext, _ := gin.CreateTestContext(recorder)
			ginContext.Request = httptest.NewRequest("GET", "/api/hosts", nil)

			handler := listHandler{
				hostCommands:     commands,
				settingsCommands: nil,
			}
			assert.PanicsWithValue(t, expectedErr, func() {
				handler.handle(ginContext)
			})
		})
	})
}
