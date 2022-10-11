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

func DataSourceNamingPolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNamingPolicyRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the OneFuse Naming Policy",
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the OneFuse Naming Policy",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description for the OneFuse Naming Policy",
			},
		},
	}
}

func dataSourceNamingPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*cbclient.CloudBoltClient)
	name := d.Get("name").(string)

	namingPolicy, err := apiClient.GetNamingPolicy(name)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(namingPolicy.ID))
	d.Set("name", namingPolicy.Name)
	d.Set("description", namingPolicy.Description)

	return nil
}
