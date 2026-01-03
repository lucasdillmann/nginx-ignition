package healthcheck

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/common/healthcheck"
)

func Test_LivenessHandler(t *testing.T) {
	t.Run("Handle", func(t *testing.T) {
		t.Run("returns 200 OK when healthy", func(t *testing.T) {
			hc := healthcheck.New()
			// No providers means healthy by default

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = httptest.NewRequest("GET", "/health/liveness", nil)

			handler := livenessHandler{
				healthCheck: hc,
			}
			handler.handle(ctx)

			assert.Equal(t, http.StatusOK, w.Code)
			var resp statusDTO
			json.Unmarshal(w.Body.Bytes(), &resp)
			assert.True(t, resp.Healthy)
		})

		t.Run("returns 503 Service Unavailable when unhealthy", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			provider := healthcheck.NewMockedProvider(ctrl)
			provider.EXPECT().ID().Return("db").AnyTimes()
			provider.EXPECT().Check(gomock.Any()).Return(assert.AnError)

			hc := healthcheck.New()
			hc.Register(provider)

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = httptest.NewRequest("GET", "/health/liveness", nil)

			handler := livenessHandler{
				healthCheck: hc,
			}
			handler.handle(ctx)

			assert.Equal(t, http.StatusServiceUnavailable, w.Code)
			var resp statusDTO
			json.Unmarshal(w.Body.Bytes(), &resp)
			assert.False(t, resp.Healthy)
		})
	})
}
