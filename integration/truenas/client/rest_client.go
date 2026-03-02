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
	response := getFromCache(cacheKey, func() []byte {
		output, err := c.executeGetRequest(endpoint, result)
		if err != nil {
			panic(err)
		}

		return output
	})

	if response != nil {
		return json.Unmarshal(response, result)
	}

	return nil
}

func (c *restClient) executeGetRequest(endpoint string, result any) ([]byte, error) {
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(body, result); err != nil {
		return nil, err
	}

	return body, nil
}
