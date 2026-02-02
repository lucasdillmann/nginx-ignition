package logline

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func Test_ExtractSearchParams(t *testing.T) {
	t.Run("returns nil when search query is missing", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		ginContext, _ := gin.CreateTestContext(recorder)
		ginContext.Request = httptest.NewRequest("GET", "/", nil)

		result := ExtractSearchParams(ginContext)

		assert.Nil(t, result)
	})

	t.Run("returns result when only search query is provided", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		ginContext, _ := gin.CreateTestContext(recorder)
		ginContext.Request = httptest.NewRequest("GET", "/?searchTerms=term", nil)

		result := ExtractSearchParams(ginContext)

		assert.NotNil(t, result)
		assert.Equal(t, "term", result.Query)
		assert.Equal(t, 0, result.SurroundingLines)
	})

	t.Run("trims search query whitespace", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		ginContext, _ := gin.CreateTestContext(recorder)
		ginContext.Request = httptest.NewRequest("GET", "/?searchTerms=%20term%20", nil)

		result := ExtractSearchParams(ginContext)

		assert.NotNil(t, result)
		assert.Equal(t, "term", result.Query)
	})

	t.Run("extracts surroundingLines when provided", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		ginContext, _ := gin.CreateTestContext(recorder)
		ginContext.Request = httptest.NewRequest(
			"GET",
			"/?searchTerms=term&surroundingLines=5",
			nil,
		)

		result := ExtractSearchParams(ginContext)

		assert.NotNil(t, result)
		assert.Equal(t, "term", result.Query)
		assert.Equal(t, 5, result.SurroundingLines)
	})

	t.Run("clamps negative surroundingLines to 0", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		ginContext, _ := gin.CreateTestContext(recorder)
		ginContext.Request = httptest.NewRequest(
			"GET",
			"/?searchTerms=term&surroundingLines=-5",
			nil,
		)

		result := ExtractSearchParams(ginContext)

		assert.NotNil(t, result)
		assert.Equal(t, 0, result.SurroundingLines)
	})

	t.Run("clamps large surroundingLines to 10", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		ginContext, _ := gin.CreateTestContext(recorder)
		ginContext.Request = httptest.NewRequest(
			"GET",
			"/?searchTerms=term&surroundingLines=100",
			nil,
		)

		result := ExtractSearchParams(ginContext)

		assert.NotNil(t, result)
		assert.Equal(t, 10, result.SurroundingLines)
	})
}
