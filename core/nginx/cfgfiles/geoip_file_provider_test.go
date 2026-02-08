package cfgfiles

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"dillmann.com.br/nginx-ignition/core/common/configuration"
)

func Test_geoIPFileProvider(t *testing.T) {
	t.Run("Provide", func(t *testing.T) {
		t.Run("successfully downloads and caches if no cache exists", func(t *testing.T) {
			tempDir := t.TempDir()
			config := configuration.NewWithOverrides(map[string]string{
				"nginx-ignition.database.data-path": tempDir,
			})

			var countryDownloadURL, cityDownloadURL string
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				case "/releases":
					fmt.Fprintf(
						w,
						`[{"tag_name": "v1.0.0", "assets": [{"name": "GeoLite2-Country.mmdb", "browser_download_url": "%s"}, {"name": "GeoLite2-City.mmdb", "browser_download_url": "%s"}]}]`,
						countryDownloadURL,
						cityDownloadURL,
					)
				case "/download-country":
					w.Write([]byte("fake-country-data"))
				case "/download-city":
					w.Write([]byte("fake-city-data"))
				default:
					panic("unexpected request: " + r.URL.Path)
				}
			}))
			defer ts.Close()
			countryDownloadURL = ts.URL + "/download-country"
			cityDownloadURL = ts.URL + "/download-city"

			provider := &geoIPFileProvider{config: config}
			// Override the URL for testing
			oldURL := geoIPReleasesURL
			geoIPReleasesURL = ts.URL + "/releases"
			defer func() { geoIPReleasesURL = oldURL }()

			ctx := newProviderContext(t)
			ctx.cfg = newSettings()
			ctx.cfg.Nginx.Stats.Enabled = true
			files, err := provider.provide(ctx)

			assert.NoError(t, err)
			assert.Len(t, files, 2)
			assert.Equal(t, geoIPCountryFileName, files[0].Name)
			assert.Equal(t, "fake-country-data", files[0].Contents)
			assert.Equal(t, geoIPCityFileName, files[1].Name)
			assert.Equal(t, "fake-city-data", files[1].Contents)

			// Check if files were cached
			assert.FileExists(t, filepath.Join(tempDir, geoIPCountryFileName))
			assert.FileExists(t, filepath.Join(tempDir, geoIPCityFileName))
			assert.FileExists(t, filepath.Join(tempDir, geoIPVersionFileName))

			versionData, _ := os.ReadFile(filepath.Join(tempDir, geoIPVersionFileName))
			assert.Equal(t, "v1.0.0", string(versionData))
		})

		t.Run("uses cached version if up to date", func(t *testing.T) {
			tempDir := t.TempDir()
			config := configuration.NewWithOverrides(map[string]string{
				"nginx-ignition.database.data-path": tempDir,
			})

			_ = os.WriteFile(
				filepath.Join(tempDir, geoIPCountryFileName),
				[]byte("cached-country"),
				0o644,
			)
			_ = os.WriteFile(
				filepath.Join(tempDir, geoIPCityFileName),
				[]byte("cached-city"),
				0o644,
			)
			_ = os.WriteFile(filepath.Join(tempDir, geoIPVersionFileName), []byte("v1.0.0"), 0o644)

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				fmt.Fprint(w, `[{"tag_name": "v1.0.0", "assets": []}]`)
			}))
			defer ts.Close()

			provider := &geoIPFileProvider{config: config}
			oldURL := geoIPReleasesURL
			geoIPReleasesURL = ts.URL
			defer func() { geoIPReleasesURL = oldURL }()

			ctx := newProviderContext(t)
			ctx.cfg = newSettings()
			ctx.cfg.Nginx.Stats.Enabled = true
			files, err := provider.provide(ctx)

			assert.NoError(t, err)
			assert.Len(t, files, 2)
			assert.Equal(t, "cached-country", files[0].Contents)
			assert.Equal(t, "cached-city", files[1].Contents)
		})

		t.Run("falls back to cache if API fails", func(t *testing.T) {
			tempDir := t.TempDir()
			config := configuration.NewWithOverrides(map[string]string{
				"nginx-ignition.database.data-path": tempDir,
			})

			_ = os.WriteFile(
				filepath.Join(tempDir, geoIPCountryFileName),
				[]byte("cached-country"),
				0o644,
			)
			_ = os.WriteFile(
				filepath.Join(tempDir, geoIPCityFileName),
				[]byte("cached-city"),
				0o644,
			)

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			}))
			defer ts.Close()

			provider := &geoIPFileProvider{config: config}
			oldURL := geoIPReleasesURL
			geoIPReleasesURL = ts.URL
			defer func() { geoIPReleasesURL = oldURL }()

			ctx := newProviderContext(t)
			ctx.cfg = newSettings()
			ctx.cfg.Nginx.Stats.Enabled = true
			files, err := provider.provide(ctx)

			assert.NoError(t, err)
			assert.Len(t, files, 2)
			assert.Equal(t, "cached-country", files[0].Contents)
			assert.Equal(t, "cached-city", files[1].Contents)
		})

		t.Run("returns error if API fails and no cache exists", func(t *testing.T) {
			tempDir := t.TempDir()
			config := configuration.NewWithOverrides(map[string]string{
				"nginx-ignition.database.data-path": tempDir,
			})

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			}))
			defer ts.Close()

			provider := &geoIPFileProvider{config: config}
			oldURL := geoIPReleasesURL
			geoIPReleasesURL = ts.URL
			defer func() { geoIPReleasesURL = oldURL }()

			ctx := newProviderContext(t)
			ctx.cfg = newSettings()
			ctx.cfg.Nginx.Stats.Enabled = true
			_, err := provider.provide(ctx)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "Failed to fetch latest GeoIP release")
		})

		t.Run("returns nothing if stats are disabled", func(t *testing.T) {
			provider := &geoIPFileProvider{config: configuration.New()}
			ctx := newProviderContext(t)
			ctx.cfg = newSettings()
			ctx.cfg.Nginx.Stats.Enabled = false

			files, err := provider.provide(ctx)
			assert.NoError(t, err)
			assert.Nil(t, files)
		})
	})
}
