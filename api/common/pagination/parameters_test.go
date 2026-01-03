package pagination

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"dillmann.com.br/nginx-ignition/api/common/apierror"
)

func Test_ExtractPaginationParameters(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("returns defaults when no params provided", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/", nil)

		pageSize, pageNumber, searchTerms, err := ExtractPaginationParameters(c)

		assert.NoError(t, err)
		assert.Equal(t, 25, pageSize)
		assert.Equal(t, 0, pageNumber)
		assert.Nil(t, searchTerms)
	})

	t.Run("returns values when valid params provided", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/?pageSize=50&pageNumber=2&searchTerms=test", nil)

		pageSize, pageNumber, searchTerms, err := ExtractPaginationParameters(c)

		assert.NoError(t, err)
		assert.Equal(t, 50, pageSize)
		assert.Equal(t, 2, pageNumber)
		assert.NotNil(t, searchTerms)
		assert.Equal(t, "test", *searchTerms)
	})

	t.Run("returns empty searchTerms as nil", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/?searchTerms=+++++", nil)

		_, _, searchTerms, err := ExtractPaginationParameters(c)

		assert.NoError(t, err)
		assert.Nil(t, searchTerms)
	})

	t.Run("returns error when pageSize is invalid int", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/?pageSize=invalid", nil)

		_, _, _, err := ExtractPaginationParameters(c)

		assert.Error(t, err)
		var apiErr *apierror.APIError
		assert.ErrorAs(t, err, &apiErr)
		assert.Equal(t, http.StatusBadRequest, apiErr.StatusCode)
	})

	t.Run("returns error when pageSize is out of range", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/?pageSize=10001", nil)

		_, _, _, err := ExtractPaginationParameters(c)

		assert.Error(t, err)
		var apiErr *apierror.APIError
		assert.ErrorAs(t, err, &apiErr)
		assert.Equal(t, http.StatusBadRequest, apiErr.StatusCode)
	})

	t.Run("returns error when pageNumber is invalid int", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/?pageNumber=invalid", nil)

		_, _, _, err := ExtractPaginationParameters(c)

		assert.Error(t, err)
		var apiErr *apierror.APIError
		assert.ErrorAs(t, err, &apiErr)
		assert.Equal(t, http.StatusBadRequest, apiErr.StatusCode)
	})

	t.Run("returns error when pageNumber is negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/?pageNumber=-1", nil)

		_, _, _, err := ExtractPaginationParameters(c)

		assert.Error(t, err)
		var apiErr *apierror.APIError
		assert.ErrorAs(t, err, &apiErr)
		assert.Equal(t, http.StatusBadRequest, apiErr.StatusCode)
	})
}
