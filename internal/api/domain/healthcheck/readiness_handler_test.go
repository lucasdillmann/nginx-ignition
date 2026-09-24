package healthcheck

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func Test_readinessHandler(t *testing.T) {
	t.Run("handle", func(t *testing.T) {
		t.Run("returns 200 OK", func(t *testing.T) {
			recorder := httptest.NewRecorder()
			ginContext, _ := gin.CreateTestContext(recorder)
			ginContext.Request = httptest.NewRequest("GET", "/health/readiness", nil)

			handler := readinessHandler{}
			handler.handle(ginContext)

			assert.Equal(t, http.StatusOK, recorder.Code)
			var response map[string]any
			json.Unmarshal(recorder.Body.Bytes(), &response)
			assert.Equal(t, true, response["ready"])
		})
	})
}
