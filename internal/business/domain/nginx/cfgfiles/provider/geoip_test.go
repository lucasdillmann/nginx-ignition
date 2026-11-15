package provider

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/configuration"
)

func newGeoIPReleaseURL(t *testing.T, countryData, cityData []byte) string {
	t.Helper()

	var server *httptest.Server
	baseURL := ""

	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/releases":
			_, _ = fmt.Fprintf(
				w,
				`[{"tag_name": "v1.0.0", "assets": [{"name": "GeoLite2-Country.mmdb", "browser_download_url": "%s/download-country"}, {"name": "GeoLite2-City.mmdb", "browser_download_url": "%s/download-city"}]}]`,
				baseURL,
				baseURL,
			)
		case "/download-country":
			_, _ = w.Write(countryData)
		case "/download-city":
			_, _ = w.Write(cityData)
		default:
			http.NotFound(w, r)
		}
	}))

	baseURL = server.URL
	t.Cleanup(server.Close)

	return server.URL + "/releases"
}

func seedGeoIPCache(
	t *testing.T,
	dir string,
	countryData, cityData []byte,
	version string,
	validated bool,
) {
	t.Helper()

	require.NoError(
		t,
		os.WriteFile(filepath.Join(dir, geoIPCountryFileName), countryData, 0o644),
	)
	require.NoError(t, os.WriteFile(filepath.Join(dir, geoIPCityFileName), cityData, 0o644))
	require.NoError(
		t,
		os.WriteFile(filepath.Join(dir, geoIPVersionFileName), []byte(version), 0o644),
	)

	if validated {
		require.NoError(
			t,
			os.WriteFile(filepath.Join(dir, geoIPValidatedFileName), []byte(version), 0o644),
		)
	}
}

func Test_validate(t *testing.T) {
	valid := newGeoIPDatabase(t)

	t.Run("accepts a valid database", func(t *testing.T) {
		assert.NoError(t, validate(valid))
	})

	t.Run("rejects a truncated database", func(t *testing.T) {
		assert.Error(t, validate(valid[:len(valid)/2]))
	})

	t.Run("rejects a database with a corrupt data section", func(t *testing.T) {
		assert.Error(t, validate(corruptGeoIPDataSection(t, valid)))
	})
}

func Test_geoIPProvider(t *testing.T) {
	t.Run("Provide", func(t *testing.T) {
		t.Run(
			"successfully downloads, validates and caches if no cache exists",
			func(t *testing.T) {
				tempDir := t.TempDir()
				config := configuration.NewWithOverrides(map[string]string{
					"nginx-ignition.database.data-path": tempDir,
				})

				valid := newGeoIPDatabase(t)
				releasesURL := newGeoIPReleaseURL(t, valid, valid)

				provider := &geoIPProvider{
					config: config,
				}
				oldURL := geoIPReleasesURL
				geoIPReleasesURL = releasesURL
				defer func() { geoIPReleasesURL = oldURL }()

				ctx := newGeoIPTestContext(t)
				files, err := provider.Provide(ctx)

				assert.NoError(t, err)
				assert.Len(t, files, 2)
				assert.Equal(t, geoIPCountryFileName, files[0].Name)
				assert.True(t, files[0].Binary)
				assert.Equal(t, valid, []byte(files[0].Contents))
				assert.Equal(t, geoIPCityFileName, files[1].Name)
				assert.True(t, files[1].Binary)
				assert.Equal(t, valid, []byte(files[1].Contents))
				assert.True(t, ctx.GeoIPAvailable)

				assert.FileExists(t, filepath.Join(tempDir, geoIPCountryFileName))
				assert.FileExists(t, filepath.Join(tempDir, geoIPCityFileName))
				assert.FileExists(t, filepath.Join(tempDir, geoIPVersionFileName))
				assert.FileExists(t, filepath.Join(tempDir, geoIPValidatedFileName))

				cachedCountry, _ := os.ReadFile(filepath.Join(tempDir, geoIPCountryFileName))
				assert.Equal(t, valid, cachedCountry)

				versionData, _ := os.ReadFile(filepath.Join(tempDir, geoIPVersionFileName))
				assert.Equal(t, "v1.0.0", string(versionData))
			},
		)

		t.Run("uses cached version if up to date and already validated", func(t *testing.T) {
			tempDir := t.TempDir()
			config := configuration.NewWithOverrides(map[string]string{
				"nginx-ignition.database.data-path": tempDir,
			})

			valid := newGeoIPDatabase(t)
			seedGeoIPCache(t, tempDir, valid, valid, "v1.0.0", true)

			releasesURL := newGeoIPReleaseURL(t, nil, nil)

			provider := &geoIPProvider{
				config: config,
			}
			oldURL := geoIPReleasesURL
			geoIPReleasesURL = releasesURL
			defer func() { geoIPReleasesURL = oldURL }()

			ctx := newGeoIPTestContext(t)
			files, err := provider.Provide(ctx)

			assert.NoError(t, err)
			assert.Len(t, files, 2)
			assert.Equal(t, valid, []byte(files[0].Contents))
			assert.True(t, ctx.GeoIPAvailable)
		})

		t.Run(
			"validates a legacy cache without a marker and writes the marker",
			func(t *testing.T) {
				tempDir := t.TempDir()
				config := configuration.NewWithOverrides(map[string]string{
					"nginx-ignition.database.data-path": tempDir,
				})

				valid := newGeoIPDatabase(t)
				seedGeoIPCache(t, tempDir, valid, valid, "v1.0.0", false)

				releasesURL := newGeoIPReleaseURL(t, nil, nil)

				provider := &geoIPProvider{
					config: config,
				}
				oldURL := geoIPReleasesURL
				geoIPReleasesURL = releasesURL
				defer func() { geoIPReleasesURL = oldURL }()

				ctx := newGeoIPTestContext(t)
				files, err := provider.Provide(ctx)

				assert.NoError(t, err)
				assert.Len(t, files, 2)
				assert.True(t, ctx.GeoIPAvailable)

				validatedData, err := os.ReadFile(
					filepath.Join(tempDir, geoIPValidatedFileName),
				)
				assert.NoError(t, err)
				assert.Equal(t, "v1.0.0", string(validatedData))
			},
		)

		t.Run("disables geoip if the cached database is invalid", func(t *testing.T) {
			tempDir := t.TempDir()
			config := configuration.NewWithOverrides(map[string]string{
				"nginx-ignition.database.data-path": tempDir,
			})

			valid := newGeoIPDatabase(t)
			corrupt := corruptGeoIPDataSection(t, valid)
			seedGeoIPCache(t, tempDir, corrupt, valid, "v1.0.0", false)

			releasesURL := newGeoIPReleaseURL(t, nil, nil)

			provider := &geoIPProvider{
				config: config,
			}
			oldURL := geoIPReleasesURL
			geoIPReleasesURL = releasesURL
			defer func() { geoIPReleasesURL = oldURL }()

			ctx := newGeoIPTestContext(t)
			files, err := provider.Provide(ctx)

			assert.NoError(t, err)
			assert.Empty(t, files)
			assert.False(t, ctx.GeoIPAvailable)
			assert.NoFileExists(t, filepath.Join(tempDir, geoIPValidatedFileName))
		})

		t.Run("falls back to cache if the downloaded database is invalid", func(t *testing.T) {
			tempDir := t.TempDir()
			config := configuration.NewWithOverrides(map[string]string{
				"nginx-ignition.database.data-path": tempDir,
			})

			valid := newGeoIPDatabase(t)
			corrupt := corruptGeoIPDataSection(t, valid)
			seedGeoIPCache(t, tempDir, valid, valid, "v0.9.0", true)

			releasesURL := newGeoIPReleaseURL(t, corrupt, corrupt)

			provider := &geoIPProvider{
				config: config,
			}
			oldURL := geoIPReleasesURL
			geoIPReleasesURL = releasesURL
			defer func() { geoIPReleasesURL = oldURL }()

			ctx := newGeoIPTestContext(t)
			files, err := provider.Provide(ctx)

			assert.NoError(t, err)
			assert.Len(t, files, 2)
			assert.Equal(t, valid, []byte(files[0].Contents))
			assert.True(t, ctx.GeoIPAvailable)

			cachedCountry, _ := os.ReadFile(filepath.Join(tempDir, geoIPCountryFileName))
			assert.Equal(t, valid, cachedCountry)
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

				valid := newGeoIPDatabase(t)
				corrupt := corruptGeoIPDataSection(t, valid)
				releasesURL := newGeoIPReleaseURL(t, corrupt, corrupt)

				provider := &geoIPProvider{
					config: config,
				}
				oldURL := geoIPReleasesURL
				geoIPReleasesURL = releasesURL
				defer func() { geoIPReleasesURL = oldURL }()

				ctx := newGeoIPTestContext(t)
				files, err := provider.Provide(ctx)

				assert.NoError(t, err)
				assert.Empty(t, files)
				assert.False(t, ctx.GeoIPAvailable)

				assert.NoFileExists(t, filepath.Join(tempDir, geoIPCountryFileName))
				assert.NoFileExists(t, filepath.Join(tempDir, geoIPCityFileName))
				assert.NoFileExists(t, filepath.Join(tempDir, geoIPVersionFileName))
			},
		)

		t.Run("falls back to cache if API fails", func(t *testing.T) {
			tempDir := t.TempDir()
			config := configuration.NewWithOverrides(map[string]string{
				"nginx-ignition.database.data-path": tempDir,
			})

			valid := newGeoIPDatabase(t)
			seedGeoIPCache(t, tempDir, valid, valid, "v0.9.0", true)

			server := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				}),
			)
			t.Cleanup(server.Close)

			provider := &geoIPProvider{
				config: config,
			}
			oldURL := geoIPReleasesURL
			geoIPReleasesURL = server.URL
			defer func() { geoIPReleasesURL = oldURL }()

			ctx := newGeoIPTestContext(t)
			files, err := provider.Provide(ctx)

			assert.NoError(t, err)
			assert.Len(t, files, 2)
			assert.Equal(t, valid, []byte(files[0].Contents))
			assert.True(t, ctx.GeoIPAvailable)
		})

		t.Run("disables geoip if API fails and no cache exists", func(t *testing.T) {
			tempDir := t.TempDir()
			config := configuration.NewWithOverrides(map[string]string{
				"nginx-ignition.database.data-path": tempDir,
			})

			server := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				}),
			)
			t.Cleanup(server.Close)

			provider := &geoIPProvider{
				config: config,
			}
			oldURL := geoIPReleasesURL
			geoIPReleasesURL = server.URL
			defer func() { geoIPReleasesURL = oldURL }()

			ctx := newGeoIPTestContext(t)
			files, err := provider.Provide(ctx)

			assert.NoError(t, err)
			assert.Empty(t, files)
			assert.False(t, ctx.GeoIPAvailable)
		})

		t.Run("returns nothing if stats are disabled", func(t *testing.T) {
			provider := &geoIPProvider{
				config: configuration.New(),
			}
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
