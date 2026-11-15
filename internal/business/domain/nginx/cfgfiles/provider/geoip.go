package provider

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

	"github.com/oschwald/maxminddb-golang"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/configuration"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/log"
)

var geoIPReleasesURL = "https://api.github.com/repos/P3TERX/GeoLite.mmdb/releases?per_page=1&page=0"

const (
	geoIPCountryFileName   = "geoip-country.mmdb"
	geoIPCityFileName      = "geoip-city.mmdb"
	geoIPVersionFileName   = "geoip.version"
	geoIPValidatedFileName = "geoip.validated"

	geoLite2CountryAssetName = "GeoLite2-Country.mmdb"
	geoLite2CityAssetName    = "GeoLite2-City.mmdb"

	geoIPValidationSampleSize = 10
)

type gitHubRelease struct {
	TagName string               `json:"tag_name"`
	Assets  []gitHubReleaseAsset `json:"assets"`
}

type gitHubReleaseAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

type geoIPProvider struct {
	config *configuration.Configuration
}

func newGeoIPProvider(config *configuration.Configuration) *geoIPProvider {
	return &geoIPProvider{
		config: config,
	}
}

type geoIPCachePaths struct {
	country   string
	city      string
	version   string
	validated string
}

func (p *geoIPProvider) Provide(ctx *Context) ([]File, error) {
	ctx.GeoIPAvailable = false

	if !ctx.Cfg.Nginx.Stats.Enabled {
		return nil, nil
	}

	dataPath, err := p.config.Get("nginx-ignition.database.data-path")
	if err != nil {
		return nil, err
	}

	cache := geoIPCachePaths{
		country:   filepath.Join(dataPath, geoIPCountryFileName),
		city:      filepath.Join(dataPath, geoIPCityFileName),
		version:   filepath.Join(dataPath, geoIPVersionFileName),
		validated: filepath.Join(dataPath, geoIPValidatedFileName),
	}

	files := p.resolve(cache)

	ctx.GeoIPAvailable = len(files) > 0

	return files, nil
}

func (p *geoIPProvider) resolve(cache geoIPCachePaths) []File {
	latestRelease, err := p.fetchLatestRelease()
	if err != nil {
		return p.useCacheOrSkip(cache, "Failed to fetch latest GeoIP release", err)
	}

	if p.isCacheUpToDate(cache, latestRelease.TagName) {
		log.Infof("Cached GeoIP databases are up to date (release %s)", latestRelease.TagName)
		return p.readCacheOrSkip(cache)
	}

	countryURL := p.findAssetURL(latestRelease, geoLite2CountryAssetName)
	cityURL := p.findAssetURL(latestRelease, geoLite2CityAssetName)

	if countryURL == "" || cityURL == "" {
		return p.useCacheOrSkip(
			cache,
			"GeoIP data files not found in the latest release assets",
			nil,
		)
	}

	countryData, err := p.download(countryURL, "Country")
	if err != nil {
		return p.useCacheOrSkip(cache, "Failed to download latest GeoIP Country data", err)
	}

	cityData, err := p.download(cityURL, "City")
	if err != nil {
		return p.useCacheOrSkip(cache, "Failed to download latest GeoIP City data", err)
	}

	if err := validate(countryData); err != nil {
		return p.useCacheOrSkip(cache, "Downloaded GeoIP Country database is invalid", err)
	}

	if err := validate(cityData); err != nil {
		return p.useCacheOrSkip(cache, "Downloaded GeoIP City database is invalid", err)
	}

	p.updateCache(cache, latestRelease.TagName, countryData, cityData)

	return buildGeoIPFiles(countryData, cityData)
}

func (p *geoIPProvider) hasCachedFiles(cache geoIPCachePaths) bool {
	return p.exists(cache.country) && p.exists(cache.city)
}

func (p *geoIPProvider) isCacheUpToDate(cache geoIPCachePaths, latestVersion string) bool {
	cachedVersion := p.readCachedVersion(cache.version)
	return cachedVersion == latestVersion && p.hasCachedFiles(cache)
}

func (p *geoIPProvider) useCacheOrSkip(
	cache geoIPCachePaths,
	message string,
	err error,
) []File {
	if p.hasCachedFiles(cache) {
		if err != nil {
			log.Warnf("%s, proceeding with cached version: %s", message, err)
		} else {
			log.Warnf("%s. Proceeding with cached version.", message)
		}

		return p.readCacheOrSkip(cache)
	}

	if err != nil {
		log.Warnf(
			"%s and no cached version is available, GeoIP will be disabled: %s",
			message,
			err,
		)
	} else {
		log.Warnf("%s and no cached version is available. GeoIP will be disabled.", message)
	}

	return nil
}

func (p *geoIPProvider) readCacheOrSkip(cache geoIPCachePaths) []File {
	countryData, cityData, err := p.readCachedData(cache)
	if err != nil {
		log.Warnf("Failed to read cached GeoIP databases, GeoIP will be disabled: %s", err)
		return nil
	}

	if !p.isCacheValidated(cache) {
		if err := validate(countryData); err != nil {
			log.Warnf(
				"Cached GeoIP Country database is invalid, GeoIP will be disabled: %s",
				err,
			)

			return nil
		}

		if err := validate(cityData); err != nil {
			log.Warnf(
				"Cached GeoIP City database is invalid, GeoIP will be disabled: %s",
				err,
			)

			return nil
		}

		p.markCacheValidated(cache)
	}

	return buildGeoIPFiles(countryData, cityData)
}

func (p *geoIPProvider) isCacheValidated(cache geoIPCachePaths) bool {
	if !p.exists(cache.validated) {
		return false
	}

	return p.readCachedVersion(cache.validated) == p.readCachedVersion(cache.version)
}

func (p *geoIPProvider) markCacheValidated(cache geoIPCachePaths) {
	version := p.readCachedVersion(cache.version)
	_ = os.WriteFile(cache.validated, []byte(version), 0o644)
}

func validate(data []byte) error {
	reader, err := maxminddb.FromBytes(data)
	if err != nil {
		return err
	}

	defer func() {
		_ = reader.Close()
	}()

	networks := reader.Networks()
	for count := 0; count < geoIPValidationSampleSize && networks.Next(); count++ {
		var record any
		if _, err := networks.Network(&record); err != nil {
			return err
		}
	}

	return networks.Err()
}

func buildGeoIPFiles(countryData, cityData []byte) []File {
	return []File{
		{
			Name:     geoIPCountryFileName,
			Contents: string(countryData),
			Binary:   true,
		},
		{
			Name:     geoIPCityFileName,
			Contents: string(cityData),
			Binary:   true,
		},
	}
}

func (p *geoIPProvider) updateCache(
	cache geoIPCachePaths,
	version string,
	countryData, cityData []byte,
) {
	_ = os.WriteFile(cache.country, countryData, 0o644)
	_ = os.WriteFile(cache.city, cityData, 0o644)
	_ = os.WriteFile(cache.version, []byte(version), 0o644)
	_ = os.WriteFile(cache.validated, []byte(version), 0o644)
}

func (p *geoIPProvider) fetchLatestRelease() (*gitHubRelease, error) {
	log.Info(
		"Checking for GeoIP database updates (courtesy of https://github.com/P3TERX/GeoLite.mmdb)...",
	)

	resp, err := p.executeRequest(geoIPReleasesURL, 10*time.Second)
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

func (p *geoIPProvider) download(url, dbType string) ([]byte, error) {
	log.Infof("Downloading GeoIP %s database from [%s]...", dbType, url)

	resp, err := p.executeRequest(url, 30*time.Minute)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func (p *geoIPProvider) executeRequest(
	url string,
	timeout time.Duration,
) (*http.Response, error) {
	client := http.Client{Timeout: timeout}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Nginx-Ignition")
	req.Header.Set("Accept", "*/*")

	return client.Do(req)
}

func (p *geoIPProvider) findAssetURL(release *gitHubRelease, assetName string) string {
	for _, asset := range release.Assets {
		if asset.Name == assetName {
			return asset.BrowserDownloadURL
		}
	}

	return ""
}

func (p *geoIPProvider) readCachedVersion(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(data))
}

func (p *geoIPProvider) readCachedData(
	cache geoIPCachePaths,
) (countryData, cityData []byte, err error) {
	countryData, err = os.ReadFile(cache.country)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read cached GeoIP Country file: %w", err)
	}

	cityData, err = os.ReadFile(cache.city)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read cached GeoIP City file: %w", err)
	}

	return countryData, cityData, nil
}

func (p *geoIPProvider) exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
