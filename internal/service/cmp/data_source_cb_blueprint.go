package cmp

import (
	"context"

	"github.com/cloudboltsoftware/cloudbolt-go-sdk/cbclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceCloudBoltBlueprint() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCloudBoltBlueprintRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The global id of a CloudBolt Blueprint, required if \"name\" not provided",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "he name of a CloudBolt Blueprint, required if \"id\" not provided",
			},
			"url_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The relative API URL path for the CloudBolt Blueprint.",
			},
		},
	}
}

func dataSourceCloudBoltBlueprintRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*cbclient.CloudBoltClient)
	name := d.Get("name").(string)
	id := d.Get("id").(string)

	if id == "" && name == "" {
		return diag.Errorf("Either name or id  is required")
	}

	var blueprint *cbclient.CloudBoltReferenceFields
	var err error
	if name != "" {
		blueprint, err = apiClient.GetBlueprint(name)
	} else {
		blueprint, err = apiClient.GetBlueprintById(id)
	}

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(blueprint.ID)
	d.Set("url_path", blueprint.Links.Self.Href)

	if name == "" {
		d.Set("name", blueprint.Name)
	}

	return nil
}
