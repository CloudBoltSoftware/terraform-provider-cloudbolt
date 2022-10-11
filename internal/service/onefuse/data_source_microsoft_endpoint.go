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

func DataSourceMicrosoftEndpoint() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMicrosoftEndpointRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the OneFuse Microsoft Endpoint",
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the OneFuse Microsoft Endpoint",
			},
		},
	}
}

func dataSourceMicrosoftEndpointRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*cbclient.CloudBoltClient)
	name := d.Get("name").(string)

	msEndpoint, err := apiClient.GetMicrosoftEndpoint(name)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(msEndpoint.ID))
	d.Set("name", msEndpoint.Name)

	return nil
}
