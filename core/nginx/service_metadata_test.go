package nginx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Service_ExtractVersion(t *testing.T) {
	svc := &service{}

	t.Run("extracts version correctly", func(t *testing.T) {
		assert.Equal(t, "1.25.3", svc.extractVersion("nginx version: nginx/1.25.3"))
		assert.Equal(t, "unknown", svc.extractVersion("invalid"))
	})
}

func Test_Service_ExtractBuildDetails(t *testing.T) {
	svc := &service{}

	t.Run("extracts build details correctly", func(t *testing.T) {
		output := "built by gcc 12.2.0\nbuilt with OpenSSL 3.0.11"
		assert.Equal(t, "by gcc 12.2.0; with OpenSSL 3.0.11", svc.extractBuildDetails(output))
	})
}

func Test_Service_ExtractTLSSNIEnabled(t *testing.T) {
	svc := &service{}

	t.Run("detects SNI support", func(t *testing.T) {
		assert.True(t, svc.extractTLSSNIEnabled("TLS SNI support enabled"))
		assert.False(t, svc.extractTLSSNIEnabled("no support"))
	})
}

func Test_Service_ExtractStaticModules(t *testing.T) {
	svc := &service{}

	t.Run("extracts static modules from configure arguments", func(t *testing.T) {
		args := "--with-http_ssl_module --with-pcre --with-http_v2_module"
		modules := svc.extractStaticModules(args)
		assert.ElementsMatch(t, []string{
			"http_ssl_module",
			"pcre",
			"http_v2_module",
		}, modules)
	})
}

func Test_Service_ExtractDynamicModules(t *testing.T) {
	svc := &service{}

	t.Run("extracts dynamic modules from configure arguments", func(t *testing.T) {
		args := "--with-http_xslt_module=dynamic --add-dynamic-module=/path/to/module_name"
		modules := svc.extractDynamicModules(args)
		assert.ElementsMatch(t, []string{
			"http_xslt_module",
			"module_name",
		}, modules)
	})
}

func Test_Service_MergeModules(t *testing.T) {
	svc := &service{}

	t.Run("merges and deduplicates modules", func(t *testing.T) {
		merged := svc.mergeModules([]string{
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
}
