package client

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type Client struct {
	delegate *http.Client
	cache    *apiCache[[]byte]
	baseURL  string
	username string
	password string
}

func New(baseURL, username, password string, cacheTimeoutSeconds int) *Client {
	return &Client{
		baseURL:  baseURL,
		username: username,
		password: password,
		delegate: &http.Client{
			Timeout: 30 * time.Second,
		},
		cache: newCache[[]byte](cacheTimeoutSeconds),
	}
}

func (c *Client) UpdateCredentials(baseURL, username, password string) {
	c.baseURL = baseURL
	c.username = username
	c.password = password
}

func (c *Client) GetAvailableApps() ([]AvailableAppDTO, error) {
	apps := make([]AvailableAppDTO, 0)
	err := c.get("app", &apps)
	if err != nil {
		return nil, err
	}

	return apps, nil
}

func (c *Client) get(endpoint string, result any) error {
	response := c.cache.get(endpoint, func() []byte {
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

func (c *Client) executeGetRequest(endpoint string, result any) ([]byte, error) {
	req, err := http.NewRequest("GET", c.baseURL+"/api/v2.0/"+endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.username, c.password)

	resp, err := c.delegate.Do(req)
	if err != nil {
		return nil, err
	}

	//nolint:errcheck
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, result)
	if err != nil {
		return nil, err
	}

	return body, nil
}
