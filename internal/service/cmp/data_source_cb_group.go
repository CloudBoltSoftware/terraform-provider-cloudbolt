package cmp

import (
	"fmt"

	"github.com/cloudboltsoftware/cloudbolt-go-sdk/cbclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceCloudBoltGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudBoltGroupRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The global id of a CloudBolt Group",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
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
	apiClient := m.(*cbclient.CloudBoltClient)
	name := d.Get("name").(string)
	id := d.Get("id").(string)

	if id == "" && name == "" {
		return fmt.Errorf("Either name or id  is required")
	}

	var group *cbclient.CloudBoltGroup
	var err error
	if name != "" {
		group, err = apiClient.GetGroup(name)
	} else {
		group, err = apiClient.GetGroupById(id)
	}

	if err != nil {
		return fmt.Errorf("Error loading CloudBolt Group: %s", err)
	}

	d.SetId(group.ID)
	d.Set("url_path", group.Links.Self.Href)

	if name == "" {
		d.Set("name", group.Name)
	}

	return nil
}
