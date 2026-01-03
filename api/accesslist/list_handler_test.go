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

func Test_ListHandler_Handle(t *testing.T) {
	page := pagination.New(1, 10, 1, []accesslist.AccessList{{Name: "Test"}})

	t.Run("returns 200 OK with paginated results", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		commands := accesslist.NewMockedCommands(ctrl)
		commands.EXPECT().
			List(gomock.Any(), 10, 1, gomock.Any()).
			Return(page, nil)

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/api/access-lists?pageSize=10&pageNumber=1", nil)

		handler := listHandler{commands}
		handler.handle(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		var response pagination.Page[accessListResponseDTO]
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 1, response.TotalItems)
		assert.Equal(t, "Test", response.Contents[0].Name)
	})

	t.Run("passes search terms to command", func(t *testing.T) {
		searchTerm := "test-term"
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		commands := accesslist.NewMockedCommands(ctrl)
		commands.EXPECT().
			List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Eq(&searchTerm)).
			Return(page, nil)

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest(
			"GET",
			"/api/access-lists?searchTerms="+searchTerm+"&pageSize=10&pageNumber=1",
			nil,
		)

		handler := listHandler{commands}
		handler.handle(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("panics when command returns error", func(t *testing.T) {
		expectedErr := errors.New("command error")
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		commands := accesslist.NewMockedCommands(ctrl)
		commands.EXPECT().
			List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, expectedErr)

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/api/access-lists", nil)

		handler := listHandler{commands}
		assert.PanicsWithValue(t, expectedErr, func() {
			handler.handle(ctx)
		})
	})
}
