package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type restClient struct {
	delegate *http.Client
	baseURL  string
	username string
	password string
}

func newRestClient(baseURL, username, password string) *restClient {
	return &restClient{
		baseURL:  baseURL,
		username: username,
		password: password,
		delegate: &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *restClient) GetAvailableApps() ([]AvailableAppDTO, error) {
	var apps []AvailableAppDTO
	if err := c.get("app", &apps); err != nil {
		return nil, err
	}

	return apps, nil
}

func (c *restClient) get(endpoint string, result any) error {
	cacheKey := fmt.Sprintf("%s:%s:rest:%s", c.baseURL, c.username, endpoint)
	response, err := getFromCache(cacheKey, func() (*[]byte, error) {
		res, err := c.executeGetRequest(endpoint)
		if err != nil {
			return nil, err
		}

		return &res, nil
	})
	if err != nil {
		return err
	}

	if response != nil {
		return json.Unmarshal(*response, result)
	}

	return nil
}

func (c *restClient) executeGetRequest(endpoint string) ([]byte, error) {
	req, err := http.NewRequest("GET", c.baseURL+"/api/v2.0/"+endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.username, c.password)

	//nolint:gosec // G704: req is constructed with a configured base URL and hardcoded endpoints
	resp, err := c.delegate.Do(req)
	if err != nil {
		return nil, err
	}

	//nolint:errcheck
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected HTTP status %d from TrueNAS REST API", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}
