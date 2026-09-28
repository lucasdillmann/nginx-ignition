package provider

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_mimeTypesProvider(t *testing.T) {
	t.Run("Provide", func(t *testing.T) {
		provider := &mimeTypesProvider{}
		ctx := newProviderContext(t)
		files, err := provider.Provide(ctx)

		assert.NoError(t, err)
		assert.Len(t, files, 1)
		assert.Equal(t, "mime.types", files[0].Name)
		assert.Contains(t, files[0].Contents, "text/html")
		assert.Contains(t, files[0].Contents, "image/png")
	})
}
