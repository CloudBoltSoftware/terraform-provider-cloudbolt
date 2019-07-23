package cloudbolt

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceCloudBoltGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudBoltGroupRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name or absolulte path to the CloudBolt Group",
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
	apiClient := m.(Config).APIClient

	group, err := apiClient.GetGroup(d.Get("name").(string))

	if err != nil {
		return fmt.Errorf("Error loading CloudBolt Group: %s", err)
	}

	d.SetId(group.ID)
	d.Set("url_path", group.Links.Self.Href)

	return nil
}
