package cloudbolt

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCloudBoltObject() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudBoltObjectRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name or absolute path to the CloudBolt Object.",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name the CloudBolt Object to search.",
			},
			"url_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The relative API URL path for the CloudBolt Object.",
			},
		},
	}
}

func dataSourceCloudBoltObjectRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(Config).APIClient

	log.Printf("[!!] apiClient in dataSourceCloudBoltObjectRead: %+v", apiClient)

	obj, err := apiClient.GetCloudBoltObject(d.Get("type").(string), d.Get("name").(string))
	if err != nil {
		return fmt.Errorf("Error loading CloudBolt Object: %s", err)
	}

	d.SetId(obj.ID)
	d.Set("url_path", obj.Links.Self.Href)

	log.Printf("[!!] url_path: %+v", obj.Links.Self.Href)

	return nil
}
