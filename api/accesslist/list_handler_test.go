package accesslist

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/accesslist"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func Test_listHandler(t *testing.T) {
	t.Run("handle", func(t *testing.T) {
		t.Run("returns 200 OK with paginated results", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			page := newAccessListPage()
			commands := accesslist.NewMockedCommands(controller)
			commands.EXPECT().
				List(gomock.Any(), 10, 1, gomock.Any()).
				Return(page, nil)

			recorder := httptest.NewRecorder()
			ginContext, _ := gin.CreateTestContext(recorder)
			ginContext.Request = httptest.NewRequest(
				"GET",
				"/api/access-lists?pageSize=10&pageNumber=1",
				nil,
			)

			handler := listHandler{
				commands: commands,
			}
			handler.handle(ginContext)

			assert.Equal(t, http.StatusOK, recorder.Code)
			var response pagination.Page[accessListResponseDTO]
			err := json.Unmarshal(recorder.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, 1, response.TotalItems)
			assert.Equal(t, "Test", response.Contents[0].Name)
		})

		t.Run("passes search terms to command", func(t *testing.T) {
			searchTerm := "test-term"
			controller := gomock.NewController(t)
			defer controller.Finish()

			page := newAccessListPage()
			commands := accesslist.NewMockedCommands(controller)
			commands.EXPECT().
				List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Eq(&searchTerm)).
				Return(page, nil)

			recorder := httptest.NewRecorder()
			ginContext, _ := gin.CreateTestContext(recorder)
			ginContext.Request = httptest.NewRequest(
				"GET",
				"/api/access-lists?searchTerms="+searchTerm+"&pageSize=10&pageNumber=1",
				nil,
			)

			handler := listHandler{
				commands: commands,
			}
			handler.handle(ginContext)

			assert.Equal(t, http.StatusOK, recorder.Code)
		})

		t.Run("panics when command returns error", func(t *testing.T) {
			expectedErr := errors.New("command error")
			controller := gomock.NewController(t)
			defer controller.Finish()

			commands := accesslist.NewMockedCommands(controller)
			commands.EXPECT().
				List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil, expectedErr)

			recorder := httptest.NewRecorder()
			ginContext, _ := gin.CreateTestContext(recorder)
			ginContext.Request = httptest.NewRequest("GET", "/api/access-lists", nil)

			handler := listHandler{
				commands: commands,
			}
			assert.PanicsWithValue(t, expectedErr, func() {
				handler.handle(ginContext)
			})
		})
	})
}
