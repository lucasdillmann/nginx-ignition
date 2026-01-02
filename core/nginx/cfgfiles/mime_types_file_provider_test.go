package cfgfiles

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_MimeTypesFileProvider_Provide(t *testing.T) {
	p := &mimeTypesFileProvider{}
	ctx := &providerContext{
		context: context.Background(),
	}
	files, err := p.provide(ctx)

	assert.NoError(t, err)
	assert.Len(t, files, 1)
	assert.Equal(t, "mime.types", files[0].Name)
	assert.Contains(t, files[0].Contents, "text/html")
	assert.Contains(t, files[0].Contents, "image/png")
}
