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

func (p *geoIPFileProvider) provide(ctx *providerContext) ([]File, error) {
	if !ctx.cfg.Nginx.Stats.Enabled {
		return nil, nil
	}

	dataPath, err := p.config.Get("nginx-ignition.database.data-path")
	if err != nil {
		return nil, err
	}

	cachedCountryPath := filepath.Join(dataPath, geoIPCountryFileName)
	cachedCityPath := filepath.Join(dataPath, geoIPCityFileName)
	cachedVersionPath := filepath.Join(dataPath, geoIPVersionFileName)

	latestRelease, err := p.fetchLatestRelease()
	if err != nil {
		if p.exists(cachedCountryPath) && p.exists(cachedCityPath) {
			log.Warnf(
				"Failed to fetch latest GeoIP release, proceeding with cached version: %s",
				err,
			)
			return p.readCachedFiles(cachedCountryPath, cachedCityPath)
		}

		return nil, fmt.Errorf(
			"failed to fetch latest GeoIP release and no cached version is available: %w",
			err,
		)
	}

	cachedVersion := p.readCachedVersion(cachedVersionPath)
	if cachedVersion == latestRelease.TagName && p.exists(cachedCountryPath) &&
		p.exists(cachedCityPath) {
		log.Infof("Cached GeoIP databases are up to date (release %s)", cachedVersion)
		return p.readCachedFiles(cachedCountryPath, cachedCityPath)
	}

	countryURL := p.findAssetURL(latestRelease, geoLite2CountryAssetName)
	cityURL := p.findAssetURL(latestRelease, geoLite2CityAssetName)

	if countryURL == "" || cityURL == "" {
		if p.exists(cachedCountryPath) && p.exists(cachedCityPath) {
			log.Warnf(
				"GeoIP data files not found in the latest release assets. Proceeding with cached version.",
			)
			return p.readCachedFiles(cachedCountryPath, cachedCityPath)
		}

		return nil, errors.New(
			"GeoIP data files not found in the latest release assets and no cached version available",
		)
	}

	countryData, err := p.download(countryURL, "Country")
	if err != nil {
		if p.exists(cachedCountryPath) && p.exists(cachedCityPath) {
			log.Warnf(
				"Failed to download latest GeoIP Country data, proceeding with cached version: %s",
				err,
			)
			return p.readCachedFiles(cachedCountryPath, cachedCityPath)
		}

		return nil, fmt.Errorf(
			"failed to download latest GeoIP Country data and no cached version available: %w",
			err,
		)
	}

	cityData, err := p.download(cityURL, "City")
	if err != nil {
		if p.exists(cachedCountryPath) && p.exists(cachedCityPath) {
			log.Warnf(
				"Failed to download latest GeoIP City data, proceeding with cached version: %s",
				err,
			)
			return p.readCachedFiles(cachedCountryPath, cachedCityPath)
		}

		return nil, fmt.Errorf(
			"failed to download latest GeoIP City data and no cached version available: %w",
			err,
		)
	}

	_ = os.WriteFile(cachedCountryPath, countryData, 0o644)
	_ = os.WriteFile(cachedCityPath, cityData, 0o644)
	_ = os.WriteFile(cachedVersionPath, []byte(latestRelease.TagName), 0o644)

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
