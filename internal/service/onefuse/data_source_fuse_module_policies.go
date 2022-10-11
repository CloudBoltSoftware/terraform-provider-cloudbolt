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

func DataSourceModulePolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIPAMPolicyRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the OneFuse Module Policy",
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the OneFuse Module Policy",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description for the OneFuse Module Policy",
			},
		},
	}
}

func dataSourceModulePolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*cbclient.CloudBoltClient)
	name := d.Get("name").(string)

	modulePolicy, err := apiClient.GetModulePolicy(name)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(modulePolicy.ID))
	d.Set("name", modulePolicy.Name)
	d.Set("description", modulePolicy.Description)

	return nil
}
