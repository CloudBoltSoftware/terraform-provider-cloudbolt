package cmp

import (
	"fmt"

	"github.com/cloudboltsoftware/cloudbolt-go-sdk/cbclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceCloudBoltBlueprint() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudBoltBlueprintRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The global id of a CloudBolt Blueprint",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of a CloudBolt Blueprint",
			},
			"url_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The relative API URL path for the CloudBolt Blueprint.",
			},
		},
	}
}

func dataSourceCloudBoltBlueprintRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*cbclient.CloudBoltClient)
	name := d.Get("name").(string)
	id := d.Get("id").(string)

	if id == "" && name == "" {
		return fmt.Errorf("Either name or id  is required")
	}

	var blueprint *cbclient.CloudBoltReferenceFields
	var err error
	if name != "" {
		blueprint, err = apiClient.GetBlueprint(name)
	} else {
		blueprint, err = apiClient.GetBlueprintById(id)
	}

	if err != nil {
		return fmt.Errorf("Error loading CloudBolt Blueprint: %s", err)
	}

	d.SetId(blueprint.ID)
	d.Set("url_path", blueprint.Links.Self.Href)

	if name == "" {
		d.Set("name", blueprint.Name)
	}

	return nil
}
