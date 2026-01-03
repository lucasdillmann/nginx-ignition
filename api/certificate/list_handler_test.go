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
	corepagination "dillmann.com.br/nginx-ignition/core/common/pagination"
)

func Test_ListHandler_Handle(t *testing.T) {
	mockCerts := []certificate.Certificate{
		{ProviderID: "p1"},
		{ProviderID: "p2"},
	}

	t.Run("returns 200 OK with certificate list on success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		page := &corepagination.Page[certificate.Certificate]{
			Contents: mockCerts,
		}

		commands := certificate.NewMockedCommands(ctrl)
		commands.EXPECT().
			List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(page, nil)

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/api/certificates?pageSize=10&pageNumber=1", nil)

		handler := listHandler{commands}
		handler.handle(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp pagination.PageDTO[certificateResponse]
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Len(t, resp.Contents, 2)
	})

	t.Run("panics when command returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		expectedErr := errors.New("list error")
		commands := certificate.NewMockedCommands(ctrl)
		commands.EXPECT().
			List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, expectedErr)

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/api/certificates", nil)

		handler := listHandler{commands}
		assert.PanicsWithValue(t, expectedErr, func() {
			handler.handle(ctx)
		})
	})
}
