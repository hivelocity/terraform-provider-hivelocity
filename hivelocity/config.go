package hivelocity

import (
	"context"
	"crypto/tls"
	"net/http"

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

// Client wraps the Hivelocity Client
type Client struct {
	client  *hv.APIClient
	auth    context.Context
	ApiUrl  string
	ApiKey  string
	Referer string
}

// Client configures govultr and returns an initialized client
func (c *Config) Client() (*Client, error) {
	transport := &http.Transport{
		TLSNextProto: make(map[string]func(string, *tls.Conn) http.RoundTripper),
	}
	client := http.DefaultClient
	client.Transport = transport

	client.Transport = logging.NewTransport("Hivelocity", client.Transport)

	conf := hv.NewConfiguration()
	conf.BasePath = c.ApiUrl
	conf.DefaultHeader = make(map[string]string)
	conf.DefaultHeader["Referer"] = c.Referer // convert Referer to map[string]string
	hvClient := hv.NewAPIClient(conf)
	authContext := context.WithValue(context.Background(), hv.ContextAPIKey, hv.APIKey{
		Key: c.ApiKey,
	})

	return &Client{client: hvClient, auth: authContext, ApiUrl: c.ApiUrl, ApiKey: c.ApiKey, Referer: c.Referer}, nil
}
