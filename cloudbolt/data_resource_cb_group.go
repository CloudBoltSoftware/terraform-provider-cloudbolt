package cloudbolt

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceCloudBoltGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudBoltGroupRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name or absolute path to the CloudBolt Group",
			},
			"url_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The relative API URL path for the CloudBolt Group.",
			},
		},
	}
}

func dataSourceCloudBoltGroupRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("in dataSoruceCloudBoltGroupRead")
	apiClient := m.(Config).APIClient

	log.Printf("[!!] apiClient: %+v", apiClient)

	group, err := apiClient.GetGroup(d.Get("name").(string))

	log.Printf("[!!] group : %+v", group)

	if err != nil {
		return fmt.Errorf("Error loading CloudBolt Group: %s", err)
	}

	d.SetId(group.ID)
	d.Set("url_path", group.Links.Self.Href)

	return nil
}
