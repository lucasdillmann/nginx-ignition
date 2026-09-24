package pagination

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_New(t *testing.T) {
	t.Run("creates page with all fields set correctly", func(t *testing.T) {
		pageNumber := 0
		pageSize := 10
		totalItems := 25
		contents := []string{"item1", "item2", "item3"}

		page := New(pageNumber, pageSize, totalItems, contents)

		assert.Equal(t, pageNumber, page.PageNumber)
		assert.Equal(t, pageSize, page.PageSize)
		assert.Equal(t, totalItems, page.TotalItems)
		assert.Equal(t, contents, page.Contents)
	})
}

func Test_Of(t *testing.T) {
	t.Run("creates page with all fields set correctly", func(t *testing.T) {
		contents := []string{"item1", "item2", "item3"}

		page := Of(contents)

		assert.Equal(t, 0, page.PageNumber)
		assert.Equal(t, len(contents), page.PageSize)
		assert.Equal(t, len(contents), page.TotalItems)
		assert.Equal(t, contents, page.Contents)
	})
}
