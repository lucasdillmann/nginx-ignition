package pagination

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("creates page with all fields set correctly", func(t *testing.T) {
		pageNumber := 1
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
