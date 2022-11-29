// Copyright 2020 CloudBolt Software
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package onefuse

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"strings"

	"github.com/cloudboltsoftware/cloudbolt-go-sdk/cbclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceAnsibleTowerDeployment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAnsibleTowerDeploymentCreate,
		ReadContext:   resourceAnsibleTowerDeploymentRead,
		UpdateContext: resourceAnsibleTowerDeploymentUpdate,
		DeleteContext: resourceAnsibleTowerDeploymentDelete,
		Schema: map[string]*schema.Schema{
			"policy_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "OneFuse Module Policy ID",
			},
			"workspace_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "OneFuse Workspace URL path,",
			},
			"limit": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Ansible Tower Policy Limit. Pattern matches hosts. or example, dev-* will match all host that start with \"dev-\"",
			},
			"hosts": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "List of specific hosts to execute against. Similar to limit, but will match against exact hostnames.",
			},
			"template_properties": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"inventory_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"provisioning_job_results": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"request_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     25,
				Description: "Timeout in minutes, Default (30)",
			},
		},
	}
}

func resourceAnsibleTowerDeploymentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("onefuse.resourceAnsibleTowerDeploymentCreate")

	var hosts []string
	for _, group := range d.Get("hosts").([]interface{}) {
		hosts = append(hosts, group.(string))
	}

	newAnsibleTowerDeployment := cbclient.AnsibleTowerDeployment{
		PolicyID:           d.Get("policy_id").(int),
		WorkspaceURL:       d.Get("workspace_url").(string),
		Hosts:              hosts,
		Limit:              d.Get("limit").(string),
		TemplateProperties: d.Get("template_properties").(map[string]interface{}),
	}

	apiClient := m.(*cbclient.CloudBoltClient)
	jobStatus, err := apiClient.CreateAnsibleTowerDeployment(&newAnsibleTowerDeployment)
	if err != nil {
		return diag.FromErr(err)
	}

	requestTimeout := d.Get("request_timeout").(int)
	stateChangeConf := GetJobStautusStateChangeConf(apiClient, requestTimeout, jobStatus.Links.Self.Href)

	_, err = stateChangeConf.WaitForState()
	if err != nil {
		return diag.Errorf("Error waiting for Job (%s) to complete. Error: %s", jobStatus.Links.Self.Href, err)
	}

	// Retrieve the updated Job Status to obtain the Managed Object
	jobStatus, err = apiClient.GetJobStatus(jobStatus.Links.Self.Href)
	if err != nil {
		return diag.FromErr(err)
	}

	var ansibleDeployment *cbclient.AnsibleTowerDeployment
	ansibleDeployment, err = apiClient.GetAnsibleTowerDeployment(jobStatus.Links.ManagedObject.Href)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(ansibleDeployment.ID))

	return resourceAnsibleTowerDeploymentRead(ctx, d, m)
}

func resourceAnsibleTowerDeploymentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*cbclient.CloudBoltClient)
	var diags diag.Diagnostics

	ansibleDeployment, err := apiClient.GetAnsibleTowerDeploymentById(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("workspace_url", ansibleDeployment.Links.Workspace.Href); err != nil {
		diags = append(diags, diag.Errorf("Cannot set workspace: %s", ansibleDeployment.Links.Workspace.Href)...)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	if err := d.Set("hosts", ansibleDeployment.Hosts); err != nil {
		hosts := strings.Join(ansibleDeployment.Hosts[:], ",")

		diags = append(diags, diag.Errorf("Cannot set hosts: %s", hosts)...)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	if err := d.Set("limit", ansibleDeployment.Limit); err != nil {
		diags = append(diags, diag.Errorf("Cannot set limit: %s", ansibleDeployment.Limit)...)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	if err := d.Set("inventory_name", ansibleDeployment.InventoryName); err != nil {
		diags = append(diags, diag.Errorf("Cannot set inventory name: %s", ansibleDeployment.InventoryName)...)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	provisioningJobResultsJson, err := json.Marshal(ansibleDeployment.ProvisioningJobResults)
	if err != nil {
		diags = append(diags, diag.Errorf("Unable to Marshal provisioning_job_results into string")...)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}
	provisioningJobResultsString := string(provisioningJobResultsJson)
	if err := d.Set("provisioning_job_results", provisioningJobResultsString); err != nil {
		diags = append(diags, diag.Errorf("Cannot set provisioning_job_results: %s", provisioningJobResultsString)...)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	ansibleTowerPolicyURLSplit := strings.Split(ansibleDeployment.Links.Policy.Href, "/")
	ansibleTowerPolicyID := ansibleTowerPolicyURLSplit[len(ansibleTowerPolicyURLSplit)-2]
	ansibleTowerPolicyIDInt, _ := strconv.Atoi(ansibleTowerPolicyID)
	if err := d.Set("policy_id", ansibleTowerPolicyIDInt); err != nil {
		diags = append(diags, diag.Errorf("Cannot set policy")...)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	return nil
}

func resourceAnsibleTowerDeploymentUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("onefuse.resourceAnsibleTowerDeploymentUpdate")
	log.Println("No Op!")
	return nil
}

func resourceAnsibleTowerDeploymentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("onefuse.resourceAnsibleTowerDeploymentDelete")
	apiClient := m.(*cbclient.CloudBoltClient)

	jobStatus, err := apiClient.DeleteAnsibleTowerDeployment(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	requestTimeout := d.Get("request_timeout").(int)
	stateChangeConf := GetJobStautusStateChangeConf(apiClient, requestTimeout, jobStatus.Links.Self.Href)

	_, err = stateChangeConf.WaitForState()
	if err != nil {
		return diag.Errorf("Error waiting for Job (%s) to complete. Error: %s", jobStatus.Links.Self.Href, err)
	}

	// Retrieve the updated Job Status to obtain the Managed Object
	jobStatus, err = apiClient.GetJobStatus(jobStatus.Links.Self.Href)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
