package cfgfiles

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_mimeTypesFileProvider(t *testing.T) {
	t.Run("Provide", func(t *testing.T) {
		provider := &mimeTypesFileProvider{}
		ctx := newProviderContext()
		files, err := provider.provide(ctx)

		assert.NoError(t, err)
		assert.Len(t, files, 1)
		assert.Equal(t, "mime.types", files[0].Name)
		assert.Contains(t, files[0].Contents, "text/html")
		assert.Contains(t, files[0].Contents, "image/png")
	})
}
