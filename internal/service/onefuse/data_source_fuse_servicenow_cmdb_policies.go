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

func DataSourceServiceNowCMDBPolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceServiceNowCMDBPolicyRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the OneFuse ServiceNow CMDB Policy",
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the OneFuse ServiceNow CMDB Policy",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description for the OneFuse ServiceNow CMDB Policy",
			},
		},
	}
}

func dataSourceServiceNowCMDBPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*cbclient.CloudBoltClient)
	name := d.Get("name").(string)

	snowCMDBPolicy, err := apiClient.GetServiceNowCMDBPolicy(name)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(snowCMDBPolicy.ID))
	d.Set("name", snowCMDBPolicy.Name)
	d.Set("description", snowCMDBPolicy.Description)

	return nil
}
