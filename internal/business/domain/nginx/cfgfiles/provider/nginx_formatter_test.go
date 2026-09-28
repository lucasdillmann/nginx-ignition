package provider

import (
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/accesslist"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/binding"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/cache"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/host"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/stream"
)

func Test_nginxSprintf(t *testing.T) {
	t.Run("quotes strings as complete nginx arguments", func(t *testing.T) {
		result := nginxSprintf("directive %s %d;", `host"; return 200 "owned`, 80)
		assert.Equal(t, `directive "host\"; return 200 \"owned" 80;`, result)
	})

	t.Run("escapes nginx strings", func(t *testing.T) {
		result := nginxSprintf("directive %s;", "a\\b\nc\rd\te")
		assert.Equal(t, `directive "a\\b\nc\rd\te";`, result)
	})

	t.Run("preserves trusted directives", func(t *testing.T) {
		result := nginxSprintf("server {\n%s\n}", directiveFragment("return 200;"))
		assert.Equal(t, "server {\nreturn 200;\n}", result)
	})

	t.Run("preserves trusted raw config", func(t *testing.T) {
		result := nginxSprintf("http {\n%s\n}", rawConfigFragment("return 200;\n"))
		assert.Equal(t, "http {\nreturn 200;\n\n}", result)
	})

	t.Run("quotes named strings and unknown arguments", func(t *testing.T) {
		type namedString string
		result := nginxSprintf("directive %s %s;", namedString("value"), uuid.Nil)
		assert.Equal(t, `directive "value" "00000000-0000-0000-0000-000000000000";`, result)
	})
}

func Test_nginxFprintf(t *testing.T) {
	builder := strings.Builder{}
	nginxFprintf(&builder, "directive %s %d;", `host"; return 200 "owned`, 80)
	assert.Equal(t, `directive "host\"; return 200 \"owned" 80;`, builder.String())
}

func Test_nginxFprintfArgs(t *testing.T) {
	builder := strings.Builder{}
	nginxFprintfArgs(&builder, "resolver %s valid=5s;", "8.8.8.8", "8.8.4.4")
	assert.Equal(t, `resolver "8.8.8.8" "8.8.4.4" valid=5s;`, builder.String())
}

func Test_sanitizeHeaderValue(t *testing.T) {
	assert.Equal(t, "value  test", sanitizeHeaderValue("value\r\n test"))
	assert.Equal(t, "value test", sanitizeHeaderValue("value\x00test"))
}

func Test_sanitizeHtpasswdUsername(t *testing.T) {
	assert.Equal(t, "user_name_name", sanitizeHtpasswdUsername("user:name\nname"))
}

func Test_nginxConfigInjection(t *testing.T) {
	t.Run("escapes access list values", func(t *testing.T) {
		accessList := newAccessList()
		accessList.ID = uuid.New()
		accessList.Realm = `Realm"; return 200 "injected;`
		accessList.Credentials = []accesslist.Credentials{{Username: "user"}}
		accessList.Entries = []accesslist.Entry{{
			Outcome:       accesslist.AllowOutcome,
			SourceAddress: []string{`10.0.0.1"; return 200 "injected;`},
		}}

		file := (&accessListProvider{}).buildConfFile(&accessList, newPaths())
		assert.Contains(t, file.Contents, `auth_basic "Realm\"; return 200 \"injected;";`)
		assert.Contains(t, file.Contents, `allow "10.0.0.1\"; return 200 \"injected;";`)
		assert.NotContains(t, file.Contents, "\nreturn 200;")
	})

	t.Run("escapes binding addresses", func(t *testing.T) {
		provider := &hostConfigurationProvider{}
		ctx := newProviderContext(t)
		ctx.Cfg = newSettings()
		h := newHost()
		binding := binding.Binding{
			Type: binding.HTTPBindingType,
			IP:   `127.0.0.1"; return 200 "injected;`,
			Port: 80,
		}

		result, err := provider.buildBinding(
			ctx,
			&h,
			&binding,
			nil,
			"server_name _;",
			"",
			"",
			"",
		)
		assert.NoError(t, err)
		expectedListen := `listen "127.0.0.1\"; return 200 \"injected;:80" ` +
			`default_server;`
		assert.Contains(t, result, expectedListen)
		assert.NotContains(t, result, "\nreturn 200;")
	})

	t.Run("escapes route and header values", func(t *testing.T) {
		provider := &hostConfigurationProvider{}
		ctx := newProviderContext(t)
		h := newHost()
		route := host.Route{
			Priority:   1,
			SourcePath: `/safe"; return 200 "injected;`,
			Response: &host.RouteStaticResponse{
				StatusCode: 200,
				Headers: map[string]string{
					`X-Test`: "value\r\ninjected: true",
				},
			},
		}

		result := provider.buildStaticResponseRoute(ctx, &h, &route)
		assert.Contains(t, result, `location "/safe\"; return 200 \"injected;" {`)
		assert.Contains(t, result, `add_header "X-Test" "value injected: true" always;`)
		assert.NotContains(t, result, "\nreturn 200;")
	})

	t.Run("escapes stream backend addresses", func(t *testing.T) {
		provider := &streamProvider{}
		upstream, err := provider.buildUpstream([]stream.Backend{{
			Address: stream.Address{
				Protocol: stream.TCPProtocol,
				Address:  `10.0.0.1"; return 200 "injected;`,
				Port:     new(8080),
			},
		}}, "test_upstream")
		assert.NoError(t, err)
		assert.Contains(t, *upstream, `server "10.0.0.1\"; return 200 \"injected;:8080";`)
		assert.NotContains(t, *upstream, "\nreturn 200;")
	})

	t.Run("escapes cache expressions and extensions", func(t *testing.T) {
		provider := &hostConfigurationProvider{}
		cacheID := uuid.New()
		cacheConfig := newCache()
		cacheConfig.ID = cacheID
		cacheConfig.BypassRules = []string{`$http_cache_control; return 200 "injected`}
		cacheConfig.FileExtensions = []string{`jpg" ; return 200 "injected`}

		result := provider.buildCacheConfig([]cache.Cache{cacheConfig}, &cacheID)
		expectedBypass := `proxy_cache_bypass "$http_cache_control; return 200 ` +
			`\"injected";`
		assert.Contains(t, result, expectedBypass)
		assert.Contains(t, result, `if ($uri !~* `)
		assert.NotContains(t, result, "\nreturn 200;")
	})
}
