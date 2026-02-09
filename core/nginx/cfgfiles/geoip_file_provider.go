package cfgfiles

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/common/log"
)

var geoIPReleasesURL = "https://api.github.com/repos/P3TERX/GeoLite.mmdb/releases?per_page=1&page=0"

const (
	geoIPCountryFileName = "geoip-country.mmdb"
	geoIPCityFileName    = "geoip-city.mmdb"
	geoIPVersionFileName = "geoip.version"

	geoLite2CountryAssetName = "GeoLite2-Country.mmdb"
	geoLite2CityAssetName    = "GeoLite2-City.mmdb"
)

type gitHubRelease struct {
	TagName string               `json:"tag_name"`
	Assets  []gitHubReleaseAsset `json:"assets"`
}

type gitHubReleaseAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

type geoIPFileProvider struct {
	config *configuration.Configuration
}

func newGeoIPFileProvider(config *configuration.Configuration) *geoIPFileProvider {
	return &geoIPFileProvider{
		config: config,
	}
}

type geoIPCachePaths struct {
	country string
	city    string
	version string
}

func (p *geoIPFileProvider) provide(ctx *providerContext) ([]File, error) {
	if !ctx.cfg.Nginx.Stats.Enabled {
		return nil, nil
	}

	dataPath, err := p.config.Get("nginx-ignition.database.data-path")
	if err != nil {
		return nil, err
	}

	cache := geoIPCachePaths{
		country: filepath.Join(dataPath, geoIPCountryFileName),
		city:    filepath.Join(dataPath, geoIPCityFileName),
		version: filepath.Join(dataPath, geoIPVersionFileName),
	}

	latestRelease, err := p.fetchLatestRelease()
	if err != nil {
		return p.fallbackToCacheOrError(cache, "Failed to fetch latest GeoIP release", err)
	}

	if p.isCacheUpToDate(cache, latestRelease.TagName) {
		log.Infof("Cached GeoIP databases are up to date (release %s)", latestRelease.TagName)
		return p.readCachedFiles(cache.country, cache.city)
	}

	countryURL := p.findAssetURL(latestRelease, geoLite2CountryAssetName)
	cityURL := p.findAssetURL(latestRelease, geoLite2CityAssetName)

	if countryURL == "" || cityURL == "" {
		return p.fallbackToCacheOrError(
			cache,
			"GeoIP data files not found in the latest release assets",
			nil,
		)
	}

	countryData, cityData, err := p.downloadBothDatabases(cache, countryURL, cityURL)
	if err != nil {
		return nil, err
	}

	p.updateCache(cache, latestRelease.TagName, countryData, cityData)

	return []File{
		{Name: geoIPCountryFileName, Contents: string(countryData)},
		{Name: geoIPCityFileName, Contents: string(cityData)},
	}, nil
}

func (p *geoIPFileProvider) hasCachedFiles(cache geoIPCachePaths) bool {
	return p.exists(cache.country) && p.exists(cache.city)
}

func (p *geoIPFileProvider) isCacheUpToDate(cache geoIPCachePaths, latestVersion string) bool {
	cachedVersion := p.readCachedVersion(cache.version)
	return cachedVersion == latestVersion && p.hasCachedFiles(cache)
}

func (p *geoIPFileProvider) fallbackToCacheOrError(
	cache geoIPCachePaths,
	message string,
	err error,
) ([]File, error) {
	if p.hasCachedFiles(cache) {
		if err != nil {
			log.Warnf("%s, proceeding with cached version: %s", message, err)
		} else {
			log.Warnf("%s. Proceeding with cached version.", message)
		}

		return p.readCachedFiles(cache.country, cache.city)
	}

	if err != nil {
		return nil, fmt.Errorf("%s and no cached version is available: %w", message, err)
	}

	return nil, errors.New(message + " and no cached version available")
}

func (p *geoIPFileProvider) downloadBothDatabases(
	cache geoIPCachePaths,
	countryURL, cityURL string,
) (countryData, cityData []byte, err error) {
	countryData, err = p.download(countryURL, "Country")
	if err != nil {
		files, cacheErr := p.fallbackToCacheOrError(
			cache,
			"Failed to download latest GeoIP Country data",
			err,
		)
		if cacheErr != nil {
			return nil, nil, cacheErr
		}

		return []byte(files[0].Contents), []byte(files[1].Contents), nil
	}

	cityData, err = p.download(cityURL, "City")
	if err != nil {
		files, cacheErr := p.fallbackToCacheOrError(
			cache,
			"Failed to download latest GeoIP City data",
			err,
		)
		if cacheErr != nil {
			return nil, nil, cacheErr
		}

		return []byte(files[0].Contents), []byte(files[1].Contents), nil
	}

	return countryData, cityData, nil
}

func (p *geoIPFileProvider) updateCache(
	cache geoIPCachePaths,
	version string,
	countryData, cityData []byte,
) {
	_ = os.WriteFile(cache.country, countryData, 0o644)
	_ = os.WriteFile(cache.city, cityData, 0o644)
	_ = os.WriteFile(cache.version, []byte(version), 0o644)
}

func (p *geoIPFileProvider) fetchLatestRelease() (*gitHubRelease, error) {
	log.Info(
		"Checking for GeoIP database updates (courtesy of https://github.com/P3TERX/GeoLite.mmdb)...",
	)
	client := http.Client{Timeout: 2 * time.Second}

	resp, err := client.Get(geoIPReleasesURL)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var releases []gitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return nil, err
	}

	if len(releases) == 0 {
		return nil, errors.New("no releases found")
	}

	return &releases[0], nil
}

func (p *geoIPFileProvider) download(url, dbType string) ([]byte, error) {
	log.Infof("Downloading GeoIP %s database from [%s] ...", dbType, url)
	client := http.Client{Timeout: 60 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func (p *geoIPFileProvider) findAssetURL(release *gitHubRelease, assetName string) string {
	for _, asset := range release.Assets {
		if asset.Name == assetName {
			return asset.BrowserDownloadURL
		}
	}

	return ""
}

func (p *geoIPFileProvider) readCachedVersion(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(data))
}

func (p *geoIPFileProvider) readCachedFiles(countryPath, cityPath string) ([]File, error) {
	countryData, err := os.ReadFile(countryPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read cached GeoIP Country file: %w", err)
	}

	cityData, err := os.ReadFile(cityPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read cached GeoIP City file: %w", err)
	}

	return []File{
		{
			Name:     geoIPCountryFileName,
			Contents: string(countryData),
		},
		{
			Name:     geoIPCityFileName,
			Contents: string(cityData),
		},
	}, nil
}

func (p *geoIPFileProvider) exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
