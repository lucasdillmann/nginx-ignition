package nginx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetadata_SNISupportType(t *testing.T) {
	t.Run("returns StaticSupportType when tlsSniEnabled is true", func(t *testing.T) {
		m := &Metadata{
			tlsSniEnabled: true,
		}
		assert.Equal(t, StaticSupportType, m.SNISupportType())
	})

	t.Run("returns NoneSupportType when tlsSniEnabled is false", func(t *testing.T) {
		m := &Metadata{
			tlsSniEnabled: false,
		}
		assert.Equal(t, NoneSupportType, m.SNISupportType())
	})
}

func TestMetadata_StreamSupportType(t *testing.T) {
	t.Run("returns StaticSupportType when stream module is present", func(t *testing.T) {
		m := &Metadata{
			Modules: []string{
				"stream",
			},
		}
		assert.Equal(t, StaticSupportType, m.StreamSupportType())
	})

	t.Run("returns DynamicSupportType when ngx_stream_module is present", func(t *testing.T) {
		m := &Metadata{
			Modules: []string{
				"ngx_stream_module",
			},
		}
		assert.Equal(t, DynamicSupportType, m.StreamSupportType())
	})

	t.Run("returns NoneSupportType when neither module is present", func(t *testing.T) {
		m := &Metadata{
			Modules: []string{
				"other",
			},
		}
		assert.Equal(t, NoneSupportType, m.StreamSupportType())
	})
}

func TestMetadata_RunCodeSupportType(t *testing.T) {
	t.Run("returns DynamicSupportType when all required modules are present", func(t *testing.T) {
		m := &Metadata{
			Modules: []string{
				"ngx_http_js_module",
				"ngx_http_lua_module",
				"ndk_http_module",
			},
		}
		assert.Equal(t, DynamicSupportType, m.RunCodeSupportType())
	})

	t.Run("returns NoneSupportType when any required module is missing", func(t *testing.T) {
		m := &Metadata{
			Modules: []string{
				"ngx_http_js_module",
				"ngx_http_lua_module",
			},
		}
		assert.Equal(t, NoneSupportType, m.RunCodeSupportType())
	})
}
