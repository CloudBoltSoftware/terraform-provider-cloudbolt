// Copyright 2020 CloudBolt Software
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package onefuse

import (
	"context"
	"strconv"

	"github.com/cloudboltsoftware/cloudbolt-go-sdk/cbclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceStaticPropertySet() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceStaticPropertySetRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the Property Set",
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the Property Set",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description for the OneFuse Property Set",
			},
			"properties": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "The properties of the OneFuse Property Set",
			},
			"raw": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The properties as a string of the OneFuse Property Set",
			},
		},
	}
}

func dataSourceStaticPropertySetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*cbclient.CloudBoltClient)
	name := d.Get("name").(string)

	staticPropertySet, err := apiClient.GetStaticPropertySet(name)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(staticPropertySet.ID))
	d.Set("name", staticPropertySet.Name)
	d.Set("properties", staticPropertySet.Properties)
	d.Set("raw", staticPropertySet.Raw)

	return nil
}
