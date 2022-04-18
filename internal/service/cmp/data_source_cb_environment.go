package cmp

import (
	"context"

	"github.com/cloudboltsoftware/cloudbolt-go-sdk/cbclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceCloudBoltEnvironment() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCloudBoltEnvironmentRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The global id of a CloudBolt Environment, required if \"name\" not provided",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of a CloudBolt Environment, required if \"id\" not provided",
			},
			"url_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The relative API URL path for the CloudBolt Environment.",
			},
		},
	}
}

func dataSourceCloudBoltEnvironmentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*cbclient.CloudBoltClient)
	name := d.Get("name").(string)
	id := d.Get("id").(string)

	if id == "" && name == "" {
		return diag.Errorf("Either name or id  is required")
	}

	var environment *cbclient.CloudBoltReferenceFields
	var err error
	if name != "" {
		environment, err = apiClient.GetEnvironment(name)
	} else {
		environment, err = apiClient.GetEnvironmentById(id)
	}

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(environment.ID)
	d.Set("url_path", environment.Links.Self.Href)

	if name == "" {
		d.Set("name", environment.Name)
	}

	return nil
}
