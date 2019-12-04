package cloudbolt

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/cloudboltsoftware/cloudbolt-go-sdk/cbclient"
	"github.com/hashicorp/terraform/helper/schema"
)

type Config struct {
	APIClient *cbclient.CloudBoltClient
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	// Default the API version to v2
	if d.Get("cb_api_version").(string) == "" {
		d.Set("cb_api_version", "v2")
	}

	// Default HTTP timeout to 10 seconds
	if d.Get("cb_timeout").(int) <= 0 {
		d.Set("cb_timeout", 10)
	}

	// Default Protocol is HTTPS
	if d.Get("cb_protocol").(string) == "" {
		d.Set("cb_protocol", "https")
	}

	httpClient := &http.Client{
		// (Optional) User requested insecure transport
		Transport: &http.Transport {
			TLSClientConfig: &tls.Config {
				InsecureSkipVerify: d.Get("cb_insecure").(bool), // Default: false
			},
		},
		// (Optional) User requested timeout
		Timeout: time.Duration(d.Get("cb_timeout").(int)) * time.Second, // Default: 10 seconds
	}

	apiClient := cbclient.New(
		d.Get("cb_protocol").(string),
		d.Get("cb_host").(string),
		d.Get("cb_port").(string),
		d.Get("cb_api_version").(string),
		d.Get("cb_username").(string),
		d.Get("cb_password").(string),
		httpClient,
	)

	// TODO: Authenticate with the API client?

	config := Config{
		APIClient: apiClient,
	}

	return config, nil
}
