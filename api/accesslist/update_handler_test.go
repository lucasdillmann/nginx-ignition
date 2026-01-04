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
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/accesslist"
)

func Test_UpdateHandler(t *testing.T) {
	id := uuid.New()
	name := "Updated List"
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

	t.Run("Handle", func(t *testing.T) {
		t.Run("returns 204 No Content on success", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			commands := accesslist.NewMockedCommands(ctrl)
			commands.EXPECT().
				Save(gomock.Any(), gomock.Any()).
				DoAndReturn(func(_ context.Context, al *accesslist.AccessList) error {
					assert.Equal(t, id, al.ID)
					assert.Equal(t, name, al.Name)
					return nil
				})

			router := gin.New()
			handler := updateHandler{
				commands: commands,
			}
			router.PUT("/api/access-lists/:id", handler.handle)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(
				"PUT",
				"/api/access-lists/"+id.String(),
				bytes.NewBuffer(body),
			)
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusNoContent, w.Code)
		})

		t.Run("returns 404 Not Found when ID is invalid", func(t *testing.T) {
			router := gin.New()
			handler := updateHandler{
				commands: nil,
			}
			router.PUT("/api/access-lists/:id", handler.handle)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(
				"PUT",
				"/api/access-lists/invalid-uuid",
				bytes.NewBuffer(body),
			)
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusNotFound, w.Code)
		})

		t.Run("panics on invalid JSON", func(t *testing.T) {
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Params = gin.Params{
				{
					Key:   "id",
					Value: id.String(),
				},
			}
			ctx.Request = httptest.NewRequest(
				"PUT",
				"/api/access-lists/"+id.String(),
				bytes.NewBufferString("invalid json"),
			)
			ctx.Request.Header.Set("Content-Type", "application/json")

			handler := updateHandler{
				commands: nil,
			}
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
			ctx.Params = gin.Params{
				{
					Key:   "id",
					Value: id.String(),
				},
			}
			ctx.Request = httptest.NewRequest(
				"PUT",
				"/api/access-lists/"+id.String(),
				bytes.NewBuffer(body),
			)
			ctx.Request.Header.Set("Content-Type", "application/json")

			handler := updateHandler{
				commands: commands,
			}
			assert.PanicsWithValue(t, expectedErr, func() {
				handler.handle(ctx)
			})
		})
	})
}
