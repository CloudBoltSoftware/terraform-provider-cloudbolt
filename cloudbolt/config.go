package cloudbolt

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/laltomar/cloudbolt-go-sdk/cbclient"
)

type Config struct {
	APIClient cbclient.CloudBoltClient
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	apiClient, _ := cbclient.New(
		d.Get("cb_protocol").(string),
		d.Get("cb_host").(string),
		d.Get("cb_port").(string),
		d.Get("cb_username").(string),
		d.Get("cb_password").(string))

	config := Config{
		APIClient: apiClient,
	}

	return config, nil
}
