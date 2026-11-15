package provider

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/configuration"
)

func newTestGeoIPProvider(
	config *configuration.Configuration,
	validate func(data []byte) error,
) *geoIPProvider {
	if validate == nil {
		validate = func([]byte) error { return nil }
	}

	return &geoIPProvider{config: config, validate: validate}
}

func Test_geoIPProvider(t *testing.T) {
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

			provider := newTestGeoIPProvider(config, nil)
			// Override the URL for testing
			oldURL := geoIPReleasesURL
			geoIPReleasesURL = ts.URL + "/releases"
			defer func() { geoIPReleasesURL = oldURL }()

			ctx := newProviderContext(t)
			ctx.Cfg = newSettings()
			ctx.Cfg.Nginx.Stats.Enabled = true
			files, err := provider.Provide(ctx)

			assert.NoError(t, err)
			assert.Len(t, files, 2)
			assert.Equal(t, geoIPCountryFileName, files[0].Name)
			assert.True(t, files[0].Binary)
			assert.Equal(t, "fake-country-data", files[0].Contents)
			assert.Equal(t, geoIPCityFileName, files[1].Name)
			assert.True(t, files[1].Binary)
			assert.Equal(t, "fake-city-data", files[1].Contents)
			assert.True(t, ctx.GeoIPAvailable)

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

			provider := newTestGeoIPProvider(config, nil)
			oldURL := geoIPReleasesURL
			geoIPReleasesURL = ts.URL
			defer func() { geoIPReleasesURL = oldURL }()

			ctx := newProviderContext(t)
			ctx.Cfg = newSettings()
			ctx.Cfg.Nginx.Stats.Enabled = true
			files, err := provider.Provide(ctx)

			assert.NoError(t, err)
			assert.Len(t, files, 2)
			assert.True(t, files[0].Binary)
			assert.Equal(t, "cached-country", files[0].Contents)
			assert.True(t, files[1].Binary)
			assert.Equal(t, "cached-city", files[1].Contents)
			assert.True(t, ctx.GeoIPAvailable)
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

			provider := newTestGeoIPProvider(config, nil)
			oldURL := geoIPReleasesURL
			geoIPReleasesURL = ts.URL
			defer func() { geoIPReleasesURL = oldURL }()

			ctx := newProviderContext(t)
			ctx.Cfg = newSettings()
			ctx.Cfg.Nginx.Stats.Enabled = true
			files, err := provider.Provide(ctx)

			assert.NoError(t, err)
			assert.Len(t, files, 2)
			assert.Equal(t, "cached-country", files[0].Contents)
			assert.Equal(t, "cached-city", files[1].Contents)
			assert.True(t, ctx.GeoIPAvailable)
		})

		t.Run("falls back to cache if the downloaded database is invalid", func(t *testing.T) {
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
			_ = os.WriteFile(filepath.Join(tempDir, geoIPVersionFileName), []byte("v0.9.0"), 0o644)

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
					w.Write([]byte("corrupt-country-data"))
				case "/download-city":
					w.Write([]byte("corrupt-city-data"))
				default:
					panic("unexpected request: " + r.URL.Path)
				}
			}))
			defer ts.Close()
			countryDownloadURL = ts.URL + "/download-country"
			cityDownloadURL = ts.URL + "/download-city"

			provider := newTestGeoIPProvider(config, func([]byte) error {
				return errors.New("corrupt data")
			})
			oldURL := geoIPReleasesURL
			geoIPReleasesURL = ts.URL + "/releases"
			defer func() { geoIPReleasesURL = oldURL }()

			ctx := newProviderContext(t)
			ctx.Cfg = newSettings()
			ctx.Cfg.Nginx.Stats.Enabled = true
			files, err := provider.Provide(ctx)

			assert.NoError(t, err)
			assert.Len(t, files, 2)
			assert.Equal(t, "cached-country", files[0].Contents)
			assert.Equal(t, "cached-city", files[1].Contents)
			assert.True(t, ctx.GeoIPAvailable)

			// The corrupt data must not overwrite the cache
			cachedCountry, _ := os.ReadFile(filepath.Join(tempDir, geoIPCountryFileName))
			assert.Equal(t, "cached-country", string(cachedCountry))
			cachedVersion, _ := os.ReadFile(filepath.Join(tempDir, geoIPVersionFileName))
			assert.Equal(t, "v0.9.0", string(cachedVersion))
		})

		t.Run(
			"disables geoip if the downloaded database is invalid and no cache exists",
			func(t *testing.T) {
				tempDir := t.TempDir()
				config := configuration.NewWithOverrides(map[string]string{
					"nginx-ignition.database.data-path": tempDir,
				})

				var countryDownloadURL, cityDownloadURL string
				ts := httptest.NewServer(
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						switch r.URL.Path {
						case "/releases":
							fmt.Fprintf(
								w,
								`[{"tag_name": "v1.0.0", "assets": [{"name": "GeoLite2-Country.mmdb", "browser_download_url": "%s"}, {"name": "GeoLite2-City.mmdb", "browser_download_url": "%s"}]}]`,
								countryDownloadURL,
								cityDownloadURL,
							)
						case "/download-country":
							w.Write([]byte("corrupt-country-data"))
						case "/download-city":
							w.Write([]byte("corrupt-city-data"))
						default:
							panic("unexpected request: " + r.URL.Path)
						}
					}),
				)
				defer ts.Close()
				countryDownloadURL = ts.URL + "/download-country"
				cityDownloadURL = ts.URL + "/download-city"

				provider := newTestGeoIPProvider(config, func([]byte) error {
					return errors.New("corrupt data")
				})
				oldURL := geoIPReleasesURL
				geoIPReleasesURL = ts.URL + "/releases"
				defer func() { geoIPReleasesURL = oldURL }()

				ctx := newProviderContext(t)
				ctx.Cfg = newSettings()
				ctx.Cfg.Nginx.Stats.Enabled = true
				files, err := provider.Provide(ctx)

				assert.NoError(t, err)
				assert.Empty(t, files)
				assert.False(t, ctx.GeoIPAvailable)

				// Nothing should have been cached
				assert.NoFileExists(t, filepath.Join(tempDir, geoIPCountryFileName))
				assert.NoFileExists(t, filepath.Join(tempDir, geoIPCityFileName))
				assert.NoFileExists(t, filepath.Join(tempDir, geoIPVersionFileName))
			},
		)

		t.Run("disables geoip if API fails and no cache exists", func(t *testing.T) {
			tempDir := t.TempDir()
			config := configuration.NewWithOverrides(map[string]string{
				"nginx-ignition.database.data-path": tempDir,
			})

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			}))
			defer ts.Close()

			provider := newTestGeoIPProvider(config, nil)
			oldURL := geoIPReleasesURL
			geoIPReleasesURL = ts.URL
			defer func() { geoIPReleasesURL = oldURL }()

			ctx := newProviderContext(t)
			ctx.Cfg = newSettings()
			ctx.Cfg.Nginx.Stats.Enabled = true
			files, err := provider.Provide(ctx)

			assert.NoError(t, err)
			assert.Empty(t, files)
			assert.False(t, ctx.GeoIPAvailable)
		})

		t.Run("returns nothing if stats are disabled", func(t *testing.T) {
			provider := newTestGeoIPProvider(configuration.New(), nil)
			ctx := newProviderContext(t)
			ctx.Cfg = newSettings()
			ctx.Cfg.Nginx.Stats.Enabled = false

			files, err := provider.Provide(ctx)
			assert.NoError(t, err)
			assert.Nil(t, files)
			assert.False(t, ctx.GeoIPAvailable)
		})
	})
}
