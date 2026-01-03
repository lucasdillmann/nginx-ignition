package accesslist

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/accesslist"
)

func Test_CreateHandler_Handle(t *testing.T) {
	name := "New List"
	realm := "Realm"
	satisfyAll := true
	defaultOutcome := accesslist.AllowOutcome
	forwardAuth := true

	payload := accessListRequestDTO{
		Name:                        &name,
		Realm:                       &realm,
		SatisfyAll:                  &satisfyAll,
		DefaultOutcome:              &defaultOutcome,
		ForwardAuthenticationHeader: &forwardAuth,
	}
	body, _ := json.Marshal(payload)

	t.Run("returns 201 Created on success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		commands := accesslist.NewMockedCommands(ctrl)
		commands.EXPECT().
			Save(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, al *accesslist.AccessList) error {
				assert.Equal(t, name, al.Name)
				return nil
			})

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("POST", "/api/access-lists", bytes.NewBuffer(body))
		ctx.Request.Header.Set("Content-Type", "application/json")

		handler := createHandler{commands}
		handler.handle(ctx)

		assert.Equal(t, http.StatusCreated, w.Code)
		var resp map[string]string
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NotEmpty(t, resp["id"])
	})

	t.Run("panics on invalid JSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest(
			"POST",
			"/api/access-lists",
			bytes.NewBufferString("invalid json"),
		)
		ctx.Request.Header.Set("Content-Type", "application/json")

		handler := createHandler{nil}
		assert.Panics(t, func() {
			handler.handle(ctx)
		})
	})

	t.Run("panics when command returns error", func(t *testing.T) {
		expectedErr := errors.New("command error")
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		commands := accesslist.NewMockedCommands(ctrl)
		commands.EXPECT().
			Save(gomock.Any(), gomock.Any()).
			Return(expectedErr)

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("POST", "/api/access-lists", bytes.NewBuffer(body))
		ctx.Request.Header.Set("Content-Type", "application/json")

		handler := createHandler{commands}
		assert.PanicsWithValue(t, expectedErr, func() {
			handler.handle(ctx)
		})
	})
}
