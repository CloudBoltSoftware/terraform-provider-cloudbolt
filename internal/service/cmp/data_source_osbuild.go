package cmp

import (
	"fmt"

	"github.com/cloudboltsoftware/cloudbolt-go-sdk/cbclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceCloudBoltOSBuild() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudBoltOSBuildRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The global id of a CloudBolt OS Build",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of a CloudBolt OS Build",
			},
			"url_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The relative API URL path for the CloudBolt OS Build.",
			},
		},
	}
}

func dataSourceCloudBoltOSBuildRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*cbclient.CloudBoltClient)
	name := d.Get("name").(string)
	id := d.Get("id").(string)

	if id == "" && name == "" {
		return fmt.Errorf("Either name or id  is required")
	}

	var osbuild *cbclient.CloudBoltReferenceFields
	var err error
	if name != "" {
		osbuild, err = apiClient.GetOSBuild(name)
	} else {
		osbuild, err = apiClient.GetOSBuildById(id)
	}

	if err != nil {
		return fmt.Errorf("Error loading CloudBolt OS Build: %s", err)
	}

	d.SetId(osbuild.ID)
	d.Set("url_path", osbuild.Links.Self.Href)

	if name == "" {
		d.Set("name", osbuild.Name)
	}

	return nil
}
