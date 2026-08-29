package nginx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Metadata(t *testing.T) {
	t.Run("SNISupportType", func(t *testing.T) {
		t.Run("returns StaticSupportType when tlsSniEnabled is true", func(t *testing.T) {
			metadata := newMetadata()
			metadata.tlsSniEnabled = true
			assert.Equal(t, StaticSupportType, metadata.SNISupportType())
		})

		t.Run("returns NoneSupportType when tlsSniEnabled is false", func(t *testing.T) {
			metadata := newMetadata()
			metadata.tlsSniEnabled = false
			assert.Equal(t, NoneSupportType, metadata.SNISupportType())
		})
	})

	t.Run("StreamSupportType", func(t *testing.T) {
		t.Run("returns StaticSupportType when stream module is present", func(t *testing.T) {
			metadata := newMetadata()
			metadata.Modules = []string{"stream"}
			assert.Equal(t, StaticSupportType, metadata.StreamSupportType())
		})

		t.Run("returns DynamicSupportType when ngx_stream_module is present", func(t *testing.T) {
			metadata := newMetadata()
			metadata.Modules = []string{"ngx_stream_module"}
			assert.Equal(t, DynamicSupportType, metadata.StreamSupportType())
		})

		t.Run("returns NoneSupportType when neither module is present", func(t *testing.T) {
			metadata := newMetadata()
			metadata.Modules = []string{"other"}
			assert.Equal(t, NoneSupportType, metadata.StreamSupportType())
		})
	})

	t.Run("RunCodeSupportType", func(t *testing.T) {
		t.Run(
			"returns DynamicSupportType when all required modules are present",
			func(t *testing.T) {
				metadata := newMetadata()
				metadata.Modules = []string{
					"ngx_http_js_module",
					"ngx_http_lua_module",
					"ndk_http_module",
				}
				assert.Equal(t, DynamicSupportType, metadata.RunCodeSupportType())
			},
		)

		t.Run("returns NoneSupportType when any required module is missing", func(t *testing.T) {
			metadata := newMetadata()
			metadata.Modules = []string{
				"ngx_http_js_module",
				"ngx_http_lua_module",
			}
			assert.Equal(t, NoneSupportType, metadata.RunCodeSupportType())
		})
	})
}
