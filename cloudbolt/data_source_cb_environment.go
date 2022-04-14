package cloudbolt

import (
	"fmt"

	"github.com/cloudboltsoftware/cloudbolt-go-sdk/cbclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCloudBoltEnvironment() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudBoltEnvironmentRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The global id of a CloudBolt Environment",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of a CloudBolt Environment",
			},
			"url_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The relative API URL path for the CloudBolt Environment.",
			},
		},
	}
}

func dataSourceCloudBoltEnvironmentRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(Config).APIClient
	name := d.Get("name").(string)
	id := d.Get("id").(string)

	if id == "" && name == "" {
		return fmt.Errorf("Either name or id  is required")
	}

	var environment *cbclient.CloudBoltReferenceFields
	var err error
	if name != "" {
		environment, err = apiClient.GetEnvironment(name)
	} else {
		environment, err = apiClient.GetEnvironmentById(id)
	}

	if err != nil {
		return fmt.Errorf("Error loading CloudBolt Environment: %s", err)
	}

	d.SetId(environment.ID)
	d.Set("url_path", environment.Links.Self.Href)

	if name == "" {
		d.Set("name", environment.Name)
	}

	return nil
}
