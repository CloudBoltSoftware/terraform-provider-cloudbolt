// Copyright 2020 CloudBolt Software
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package onefuse

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/cloudboltsoftware/cloudbolt-go-sdk/cbclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceRenderedTemplate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRenderedTemplateRead,

		Schema: map[string]*schema.Schema{
			"template": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "OneFuse Template string.",
			},
			"template_properties": {
				Type:        schema.TypeMap,
				Required:    true,
				Description: "Additional properties that are referenced and rendered within the Template.",
			},
			"value": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceRenderedTemplateRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*cbclient.CloudBoltClient)

	renderedTemplate, err := apiClient.RenderTemplate(d.Get("template").(string), d.Get("template_properties").(map[string]interface{}))

	if err != nil {
		return diag.FromErr(err)
	}

	// a resource needs an ID, otherwise it will be destroyed, so here is a fun hack to make up an ID bc we dont have one
	inputVars := fmt.Sprint(d.Get("template_properties").(map[string]interface{}))
	concatVars := inputVars + d.Get("template").(string)
	id := sha256.Sum256([]byte(concatVars))

	d.SetId(fmt.Sprintf("%x", id))
	d.Set("value", renderedTemplate.Value)

	return nil
}
