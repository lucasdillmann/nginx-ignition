package nginx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_service_metadata(t *testing.T) {
	nginxService := &service{}

	t.Run("ExtractVersion", func(t *testing.T) {
		t.Run("extracts version correctly", func(t *testing.T) {
			assert.Equal(t, "1.25.3", nginxService.extractVersion("nginx version: nginx/1.25.3"))
			assert.Equal(t, "unknown", nginxService.extractVersion("invalid"))
		})
	})

	t.Run("ExtractBuildDetails", func(t *testing.T) {
		t.Run("extracts build details correctly", func(t *testing.T) {
			output := "built by gcc 12.2.0\nbuilt with OpenSSL 3.0.11"
			assert.Equal(
				t,
				"by gcc 12.2.0; with OpenSSL 3.0.11",
				nginxService.extractBuildDetails(output),
			)
		})
	})

	t.Run("ExtractTLSSNIEnabled", func(t *testing.T) {
		t.Run("detects SNI support", func(t *testing.T) {
			assert.True(t, nginxService.extractTLSSNIEnabled("TLS SNI support enabled"))
			assert.False(t, nginxService.extractTLSSNIEnabled("no support"))
		})
	})

	t.Run("ExtractStaticModules", func(t *testing.T) {
		t.Run("extracts static modules from configure arguments", func(t *testing.T) {
			args := "--with-http_ssl_module --with-pcre --with-http_v2_module"
			modules := nginxService.extractStaticModules(args)
			assert.ElementsMatch(t, []string{
				"http_ssl_module",
				"pcre",
				"http_v2_module",
			}, modules)
		})
	})

	t.Run("ExtractDynamicModules", func(t *testing.T) {
		t.Run("extracts dynamic modules from configure arguments", func(t *testing.T) {
			args := "--with-http_xslt_module=dynamic --add-dynamic-module=/path/to/module_name"
			modules := nginxService.extractDynamicModules(args)
			assert.ElementsMatch(t, []string{
				"http_xslt_module",
				"module_name",
			}, modules)
		})
	})

	t.Run("MergeModules", func(t *testing.T) {
		t.Run("merges and deduplicates modules", func(t *testing.T) {
			merged := nginxService.mergeModules([]string{
				"mod1",
				"mod2",
			}, []string{
				"mod2",
				"mod3",
			}, []string{
				"mod4",
			})
			assert.ElementsMatch(t, []string{
				"mod1",
				"mod2",
				"mod3",
				"mod4",
			}, merged)
		})
	})
}
