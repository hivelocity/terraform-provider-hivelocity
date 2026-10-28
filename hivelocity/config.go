package hivelocity

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	hv "github.com/hivelocity/terraform-provider-hivelocity/hivelocity-client-go"
)

// Config is the configuration structure used to instantiate the Hivelocity
// provider.
type Config struct {
	ApiKey  string
	ApiUrl  string
	Referer string
}

// Client wraps the generated Hivelocity API client plus the pieces needed to
// make raw requests for the handful of endpoints the generated client cannot
// model (e.g. responses whose swagger definition has no properties).
type Client struct {
	client     *hv.APIClient
	auth       context.Context
	httpClient *http.Client
	ApiUrl     string
	ApiKey     string
	Referer    string
}

// Client builds a configured Hivelocity API client.
func (c *Config) Client() (*Client, error) {
	httpClient := &http.Client{
		Timeout:   5 * time.Minute,
		Transport: logging.NewLoggingHTTPTransport(http.DefaultTransport),
	}

	apiUrl := strings.TrimRight(c.ApiUrl, "/")

	conf := hv.NewConfiguration()
	conf.BasePath = apiUrl
	conf.HTTPClient = httpClient
	conf.UserAgent = "terraform-provider-hivelocity"
	conf.AddDefaultHeader("Referer", c.Referer)

	hvClient := hv.NewAPIClient(conf)
	authContext := context.WithValue(context.Background(), hv.ContextAPIKey, hv.APIKey{
		Key: c.ApiKey,
	})

	return &Client{
		client:     hvClient,
		auth:       authContext,
		httpClient: httpClient,
		ApiUrl:     apiUrl,
		ApiKey:     c.ApiKey,
		Referer:    c.Referer,
	}, nil
}

// rawRequest performs an authenticated JSON request against the API and returns
// the raw body. It exists for endpoints whose swagger definitions are too loose
// for the generated client (empty models, omitempty on required lists, ...).
func (hv *Client) rawRequest(ctx context.Context, method, path string, body interface{}) ([]byte, *http.Response, error) {
	var reader io.Reader
	if body != nil {
		encoded, err := json.Marshal(body)
		if err != nil {
			return nil, nil, err
		}
		reader = bytes.NewReader(encoded)
	}

	req, err := http.NewRequestWithContext(ctx, method, hv.ApiUrl+path, reader)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("X-API-KEY", hv.ApiKey)
	req.Header.Set("Referer", hv.Referer)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "terraform-provider-hivelocity")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := hv.httpClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("%s %s failed: %w", method, path, err)
	}
	defer func() { _ = resp.Body.Close() }()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp, fmt.Errorf("%s %s: reading response failed: %w", method, path, err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return data, resp, fmt.Errorf("%s %s returned HTTP %d: %s", method, path, resp.StatusCode, strings.TrimSpace(string(data)))
	}

	return data, resp, nil
}

// rawJSON is rawRequest plus JSON decoding of the response into out (when out is
// not nil and the response has a body).
func (hv *Client) rawJSON(ctx context.Context, method, path string, body interface{}, out interface{}) (*http.Response, error) {
	data, resp, err := hv.rawRequest(ctx, method, path, body)
	if err != nil {
		return resp, err
	}
	if out != nil && len(bytes.TrimSpace(data)) > 0 {
		if err := json.Unmarshal(data, out); err != nil {
			return resp, fmt.Errorf("%s %s: decoding response failed: %w", method, path, err)
		}
	}
	return resp, nil
}
