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

func init() {
	gin.SetMode(gin.TestMode)
}

func Test_livenessHandler(t *testing.T) {
	t.Run("handle", func(t *testing.T) {
		t.Run("returns 200 OK when healthy", func(t *testing.T) {
			hc := healthcheck.New()

			recorder := httptest.NewRecorder()
			ginContext, _ := gin.CreateTestContext(recorder)
			ginContext.Request = httptest.NewRequest("GET", "/health/liveness", nil)

			handler := livenessHandler{
				healthCheck: hc,
			}
			handler.handle(ginContext)

			assert.Equal(t, http.StatusOK, recorder.Code)
			var response statusDTO
			json.Unmarshal(recorder.Body.Bytes(), &response)
			assert.True(t, response.Healthy)
		})

		t.Run("returns 503 Service Unavailable when unhealthy", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			provider := healthcheck.NewMockedProvider(controller)
			provider.EXPECT().ID().Return("db").AnyTimes()
			provider.EXPECT().Check(gomock.Any()).Return(assert.AnError)

			hc := healthcheck.New()
			hc.Register(provider)

			recorder := httptest.NewRecorder()
			ginContext, _ := gin.CreateTestContext(recorder)
			ginContext.Request = httptest.NewRequest("GET", "/health/liveness", nil)

			handler := livenessHandler{
				healthCheck: hc,
			}
			handler.handle(ginContext)

			assert.Equal(t, http.StatusServiceUnavailable, recorder.Code)
			var response statusDTO
			json.Unmarshal(recorder.Body.Bytes(), &response)
			assert.False(t, response.Healthy)
		})
	})
}
