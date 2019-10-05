package cloudbolt

import (
	// "log"

	"github.com/cloudboltsoftware/cloudbolt-go-sdk/cbclient"
	"github.com/hashicorp/terraform/helper/schema"
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
		d.Get("cb_password").(string),
	)
	// if err != nil {
	// 	log.Printf("[!!] API Client produced an error: %+v", err)
	// }

	config := Config{
		APIClient: apiClient,
	}

	// log.Printf("config: %+v", config)

	return config, nil
}
