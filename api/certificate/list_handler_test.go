package certificate

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
	"dillmann.com.br/nginx-ignition/core/certificate"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func Test_listHandler(t *testing.T) {
	t.Run("handle", func(t *testing.T) {
		t.Run("returns 200 OK with certificate list on success", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			page := newCertificatePage()
			commands := certificate.NewMockedCommands(controller)
			commands.EXPECT().
				List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(page, nil)

			recorder := httptest.NewRecorder()
			ginContext, _ := gin.CreateTestContext(recorder)
			ginContext.Request = httptest.NewRequest(
				"GET",
				"/api/certificates?pageSize=10&pageNumber=1",
				nil,
			)

			handler := listHandler{
				commands: commands,
			}
			handler.handle(ginContext)

			assert.Equal(t, http.StatusOK, recorder.Code)
			var response pagination.PageDTO[certificateResponse]
			json.Unmarshal(recorder.Body.Bytes(), &response)
			assert.Len(t, response.Contents, 1)
		})

		t.Run("panics when command returns error", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			expectedErr := errors.New("list error")
			commands := certificate.NewMockedCommands(controller)
			commands.EXPECT().
				List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil, expectedErr)

			recorder := httptest.NewRecorder()
			ginContext, _ := gin.CreateTestContext(recorder)
			ginContext.Request = httptest.NewRequest("GET", "/api/certificates", nil)

			handler := listHandler{
				commands: commands,
			}
			assert.PanicsWithValue(t, expectedErr, func() {
				handler.handle(ginContext)
			})
		})
	})
}
