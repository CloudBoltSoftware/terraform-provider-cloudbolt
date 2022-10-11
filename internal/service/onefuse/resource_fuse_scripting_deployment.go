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

func ResourceScriptingDeployment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceScriptingDeploymentCreate,
		ReadContext:   resourceScriptingDeploymentRead,
		UpdateContext: resourceScriptingDeploymentUpdate,
		DeleteContext: resourceScriptingDeploymentDelete,
		Schema: map[string]*schema.Schema{
			"hostname": {
				Type:     schema.TypeString,
				Required: false,
				Computed: true,
			},
			"policy_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"workspace_url": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"template_properties": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"provisioning_details": {
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

func resourceScriptingDeploymentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*cbclient.CloudBoltClient)

	newScriptingDeployment := cbclient.ScriptingDeployment{
		PolicyID:           d.Get("policy_id").(int),
		WorkspaceURL:       d.Get("workspace_url").(string),
		TemplateProperties: d.Get("template_properties").(map[string]interface{}),
	}

	jobStatus, err := apiClient.CreateScriptingDeployment(&newScriptingDeployment)
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

	var scriptingDeployment *cbclient.ScriptingDeployment
	scriptingDeployment, err = apiClient.GetScriptingDeployment(jobStatus.Links.ManagedObject.Href)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(scriptingDeployment.ID))

	return resourceScriptingDeploymentRead(ctx, d, m)
}

func resourceScriptingDeploymentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*cbclient.CloudBoltClient)
	var diags diag.Diagnostics

	scriptingDeployment, err := apiClient.GetScriptingDeploymentById(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("workspace_url", scriptingDeployment.Links.Workspace.Href); err != nil {
		diags = append(diags, diag.Errorf("Cannot set workspace: %s", scriptingDeployment.Links.Workspace.Href)...)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	if err := d.Set("hostname", scriptingDeployment.Hostname); err != nil {
		diags = append(diags, diag.Errorf("Cannot set hostname: %s", scriptingDeployment.Hostname)...)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	provisioningDetailsJson, err := json.Marshal(scriptingDeployment.ProvisioningDetails)
	if err != nil {
		diags = append(diags, diag.Errorf("Unable to Marshal provisioning_details into string")...)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	provisioningDetailsString := string(provisioningDetailsJson)
	if err := d.Set("provisioning_details", provisioningDetailsString); err != nil {
		diags = append(diags, diag.Errorf("Cannot set provisioning_details: %s", provisioningDetailsString)...)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	scriptingPolicyURLSplit := strings.Split(scriptingDeployment.Links.Policy.Href, "/")
	scriptingPolicyID := scriptingPolicyURLSplit[len(scriptingPolicyURLSplit)-2]
	scriptingPolicyIDInt, _ := strconv.Atoi(scriptingPolicyID)
	if err := d.Set("policy_id", scriptingPolicyIDInt); err != nil {
		diags = append(diags, diag.Errorf("Cannot set policy")...)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	return nil
}

func resourceScriptingDeploymentUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("onefuse.resourceScriptingDeploymentUpdate")
	log.Println("No Op!")
	return nil
}

func resourceScriptingDeploymentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*cbclient.CloudBoltClient)

	jobStatus, err := apiClient.DeleteScriptingDeployment(d.Id())
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
