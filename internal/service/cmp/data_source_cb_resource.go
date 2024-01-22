package cmp

import (
	"context"

	"github.com/cloudboltsoftware/cloudbolt-go-sdk/cbclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceCloudBoltResource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCloudBoltResourceRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The global id of a CloudBolt Resource, required if \"name\" not provided",
			},
			"url_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The relative API URL path for the CloudBolt Resource.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of a CloudBolt Resource, required if \"id\" not provided",
			},
			"create_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date the CloudBolt Resource was created",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CloudBolt Resource Status",
			},
			"attributes": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Description: "CloudBolt Resource attributes",
			},
		},
	}
}

func dataSourceCloudBoltResourceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*cbclient.CloudBoltClient)
	urlPath := d.Get("url_path").(string)
	id := d.Get("id").(string)

	if id == "" && urlPath == "" {
		return diag.Errorf("Either id or url_path is required")
	}

	var resource *cbclient.CloudBoltResource
	var err error
	if urlPath != "" {
		resource, err = apiClient.GetResource(urlPath)
	} else {
		resource, err = apiClient.GetResourceById(id)
	}

	if err != nil {
		return diag.FromErr(err)
	}

	if resource.ID == "" {
		return diag.Errorf("Resource not found.")
	}

	d.SetId(resource.ID)
	d.Set("url_path", resource.Links.Self.Href)
	d.Set("name", resource.Name)
	d.Set("create_date", resource.Created)
	d.Set("status", resource.Status)
	resAttributes, _ := parseAttributes(resource.Attributes)
	d.Set("attributes", resAttributes)

	return nil
}
