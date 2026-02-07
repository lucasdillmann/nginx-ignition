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

var geoIPReleasesURL = "https://api.github.com/repos/v2fly/geoip/releases?per_page=1&page=0"

const (
	geoIPFileName        = "geoip.dat"
	geoIPVersionFileName = "geoip.version"
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

	cachedDataPath := filepath.Join(dataPath, geoIPFileName)
	cachedVersionPath := filepath.Join(dataPath, geoIPVersionFileName)

	latestRelease, err := p.fetchLatestRelease()
	if err != nil {
		if p.exists(cachedDataPath) {
			log.Warnf(
				"Failed to fetch latest GeoIP release, proceeding with cached version: %s",
				err,
			)
			return p.readCachedFile(cachedDataPath)
		}

		return nil, fmt.Errorf(
			"failed to fetch latest GeoIP release and no cached version is available: %w",
			err,
		)
	}

	cachedVersion := p.readCachedVersion(cachedVersionPath)
	if cachedVersion == latestRelease.TagName && p.exists(cachedDataPath) {
		log.Infof("Cached GeoIP database is up to date (release %s)", cachedVersion)
		return p.readCachedFile(cachedDataPath)
	}

	downloadURL := p.findGeoIPAssetURL(latestRelease)
	if downloadURL == "" {
		if p.exists(cachedDataPath) {
			log.Warnf(
				"GeoIP data file not found in the latest release assets. Proceeding with cached version.",
			)
			return p.readCachedFile(cachedDataPath)
		}

		return nil, errors.New(
			"GeoIP data file not found in the latest release assets and no cached version available",
		)
	}

	data, err := p.download(downloadURL)
	if err != nil {
		if p.exists(cachedDataPath) {
			log.Warnf(
				"Failed to download latest GeoIP data, proceeding with cached version: %s",
				err,
			)
			return p.readCachedFile(cachedDataPath)
		}

		return nil, fmt.Errorf(
			"failed to download latest GeoIP data and no cached version available: %w",
			err,
		)
	}

	_ = os.WriteFile(cachedDataPath, data, 0o644)
	_ = os.WriteFile(cachedVersionPath, []byte(latestRelease.TagName), 0o644)

	return []File{{
		Name:     geoIPFileName,
		Contents: string(data),
	}}, nil
}

func (p *geoIPFileProvider) fetchLatestRelease() (*gitHubRelease, error) {
	log.Info("Checking for GeoIP database updates (courtesy of https://github.com/v2fly/geoip)...")
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

func (p *geoIPFileProvider) download(url string) ([]byte, error) {
	log.Infof("Downloading GeoIP database from [%s] ...", url)
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

func (p *geoIPFileProvider) findGeoIPAssetURL(release *gitHubRelease) string {
	for _, asset := range release.Assets {
		if asset.Name == geoIPFileName {
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

func (p *geoIPFileProvider) readCachedFile(path string) ([]File, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read cached GeoIP file: %w", err)
	}

	return []File{{
		Name:     geoIPFileName,
		Contents: string(data),
	}}, nil
}

func (p *geoIPFileProvider) exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
