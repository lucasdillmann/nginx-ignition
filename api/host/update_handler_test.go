package host

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/common/ptr"
	"dillmann.com.br/nginx-ignition/core/host"
)

func Test_UpdateHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	id := uuid.New()
	payload := &hostRequestDTO{
		Enabled:     ptr.Of(true),
		DomainNames: []string{"example.com"},
		FeatureSet:  &featureSetDTO{},
	}
	body, _ := json.Marshal(payload)

	t.Run("Handle", func(t *testing.T) {
		t.Run("returns 204 No Content on success", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			commands := host.NewMockedCommands(ctrl)
			commands.EXPECT().
				Save(gomock.Any(), gomock.Any()).
				Return(nil)

			handler := updateHandler{
				commands: commands,
			}
			r := gin.New()
			r.PUT("/api/hosts/:id", handler.handle)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", "/api/hosts/"+id.String(), bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusNoContent, w.Code)
		})

		t.Run("returns 404 Not Found on invalid ID", func(t *testing.T) {
			handler := updateHandler{
				commands: nil,
			}
			r := gin.New()
			r.PUT("/api/hosts/:id", handler.handle)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", "/api/hosts/invalid", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusNotFound, w.Code)
		})

		t.Run("panics on invalid JSON", func(t *testing.T) {
			handler := updateHandler{
				commands: nil,
			}
			r := gin.New()
			r.PUT("/api/hosts/:id", handler.handle)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(
				"PUT",
				"/api/hosts/"+id.String(),
				bytes.NewBufferString("invalid json"),
			)
			req.Header.Set("Content-Type", "application/json")

			assert.Panics(t, func() {
				r.ServeHTTP(w, req)
			})
		})

		t.Run("panics when command returns error", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			expectedErr := errors.New("update error")
			commands := host.NewMockedCommands(ctrl)
			commands.EXPECT().
				Save(gomock.Any(), gomock.Any()).
				Return(expectedErr)

			handler := updateHandler{
				commands: commands,
			}
			r := gin.New()
			r.PUT("/api/hosts/:id", func(c *gin.Context) {
				defer func() {
					if r := recover(); r != nil {
						assert.Equal(t, expectedErr, r)
						panic(r)
					}
				}()
				handler.handle(c)
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", "/api/hosts/"+id.String(), bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			assert.Panics(t, func() {
				r.ServeHTTP(w, req)
			})
		})
	})
}
