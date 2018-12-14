package propsd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	// DefaultEndpoint is the default properties endpoint for Propsd
	DefaultEndpoint = "http://localhost:9100/v1/properties"
	// DefaultTimeout is a reasonable timeout for Propsd connections
	DefaultTimeout  = 5 * time.Second
)

// Client represents a Propsd client
type Client struct {
	endpoint string
	httpClient http.Client
}

// NewClient creates a new Propsd client
func NewClient(endpoint string, timeout time.Duration) *Client {
	return &Client{
		endpoint: endpoint,
		httpClient: http.Client{
			Timeout: timeout,
		},
	}
}

// Properties gets the entire property set from Propsd
func (c * Client) Properties() ([]byte, error) {
	req, err := http.NewRequest("GET", c.endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	// Make sure that Close happens!
	defer resp.Body.Close() //nolint errcheck
	// Only 200 is valid - if we get any other code, we are having problems connecting to propsd
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unable to connect to propsd: expected 200, got %d", resp.StatusCode)
	}

	return ioutil.ReadAll(resp.Body)
}

// GetProperty gets a single property by name from Propsd
func (c* Client) GetProperty(key string) ([]byte, error) {
	var result map[string]interface{}

	props, err := c.Properties()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(props, &result)
	if err != nil {
		return nil, err
	}

	val, ok := result[key]
	if !ok {
		return nil, fmt.Errorf("unable to find %s in %v", key, result)
	}

	return []byte(val.(string)), nil
}
