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

func ResourceModuleDeployment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceModuleDeploymentCreate,
		ReadContext:   resourceModuleDeploymentRead,
		UpdateContext: resourceModuleDeploymentUpdate,
		DeleteContext: resourceModuleDeploymentDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: false,
				Computed: true,
			},
			"policy_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "OneFuse Module Policy ID",
			},
			"workspace_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "OneFuse Workspace URL path.",
			},
			"template_properties": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Additional properties that are referenced within the Policy.",
			},
			"provisioning_job_results": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"deprovisioning_job_results": {
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

func resourceModuleDeploymentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*cbclient.CloudBoltClient)

	newModuleDeployment := cbclient.ModuleDeployment{
		PolicyID:           d.Get("policy_id").(int),
		WorkspaceURL:       d.Get("workspace_url").(string),
		TemplateProperties: d.Get("template_properties").(map[string]interface{}),
	}

	jobStatus, err := apiClient.CreateModuleDeployment(&newModuleDeployment)
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

	var moduleDeployment *cbclient.ModuleDeployment
	moduleDeployment, err = apiClient.GetModuleDeployment(jobStatus.Links.ManagedObject.Href)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(moduleDeployment.ID))

	return resourceModuleDeploymentRead(ctx, d, m)
}

func resourceModuleDeploymentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*cbclient.CloudBoltClient)
	var diags diag.Diagnostics

	moduleDeployment, err := apiClient.GetModuleDeploymentById(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("workspace_url", moduleDeployment.Links.Workspace.Href); err != nil {
		diags = append(diags, diag.Errorf("Cannot set workspace: %s", moduleDeployment.Links.Workspace.Href)...)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	if err := d.Set("name", moduleDeployment.Name); err != nil {
		diags = append(diags, diag.Errorf("Cannot set name: %s", moduleDeployment.Name)...)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	provisioningJobResultsJson, err := json.Marshal(moduleDeployment.ProvisioningJobResults)
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

	deprovisioningJobResultsJson, err := json.Marshal(moduleDeployment.DeprovisioningJobResults)
	if err != nil {
		diags = append(diags, diag.Errorf("Unable to Marshal deprovisioning_job_results into string")...)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	deprovisioningJobResultsString := string(deprovisioningJobResultsJson)
	if err := d.Set("deprovisioning_job_results", deprovisioningJobResultsString); err != nil {
		diags = append(diags, diag.Errorf("Cannot set deprovisioning_job_results: %s", deprovisioningJobResultsString)...)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	ModulePolicyURLSplit := strings.Split(moduleDeployment.Links.Policy.Href, "/")
	ModulePolicyID := ModulePolicyURLSplit[len(ModulePolicyURLSplit)-2]
	ModulePolicyIDInt, _ := strconv.Atoi(ModulePolicyID)
	if err := d.Set("policy_id", ModulePolicyIDInt); err != nil {
		diags = append(diags, diag.Errorf("Cannot set policy")...)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	return nil
}

func resourceModuleDeploymentUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("onefuse.resourceModuleDeploymentUpdate")
	log.Println("No Op!")
	return nil
}

func resourceModuleDeploymentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*cbclient.CloudBoltClient)

	jobStatus, err := apiClient.DeleteModuleDeployment(d.Id())
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
