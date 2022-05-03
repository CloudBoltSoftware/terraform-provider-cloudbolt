package cmp

import (
	"context"

	"github.com/cloudboltsoftware/cloudbolt-go-sdk/cbclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceCloudBoltGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCloudBoltGroupRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The global id of a CloudBolt Group, required if \"name\" not provided",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name or absolute path to the CloudBolt Group, required if \"id\" not provided",
			},
			"url_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The relative API URL path for the CloudBolt Group.",
			},
		},
	}
}

func dataSourceCloudBoltGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*cbclient.CloudBoltClient)
	name := d.Get("name").(string)
	id := d.Get("id").(string)

	if id == "" && name == "" {
		return diag.Errorf("Either name or id is required")
	}

	var group *cbclient.CloudBoltGroup
	var err error
	if name != "" {
		group, err = apiClient.GetGroup(name)
	} else {
		group, err = apiClient.GetGroupById(id)
	}

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(group.ID)
	d.Set("url_path", group.Links.Self.Href)

	if name == "" {
		d.Set("name", group.Name)
	}

	return nil
}
