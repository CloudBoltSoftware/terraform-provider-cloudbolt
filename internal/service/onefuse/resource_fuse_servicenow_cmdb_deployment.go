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

func ResourceServicenowCMDBDeployment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServicenowCMDBDeploymentCreate,
		ReadContext:   resourceServicenowCMDBDeploymentRead,
		UpdateContext: resourceServicenowCMDBDeploymentUpdate,
		DeleteContext: resourceServicenowCMDBDeploymentDelete,
		Schema: map[string]*schema.Schema{
			"policy_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "OneFuse Module Policy ID.",
			},
			"workspace_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "OneFuse Workspace URL path.",
			},
			"configuration_items_info": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
				},
				Computed: true,
				Optional: true,
			},
			"execution_details": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"template_properties": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Additional properties that are referenced within the Policy.",
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

func resourceServicenowCMDBDeploymentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*cbclient.CloudBoltClient)

	newServicenowCMDBDeployment := cbclient.ServicenowCMDBDeployment{
		PolicyID:           d.Get("policy_id").(int),
		WorkspaceURL:       d.Get("workspace_url").(string),
		TemplateProperties: d.Get("template_properties").(map[string]interface{}),
	}

	jobStatus, err := apiClient.CreateServicenowCMDBDeployment(&newServicenowCMDBDeployment)
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

	var snowDeployment *cbclient.ServicenowCMDBDeployment
	snowDeployment, err = apiClient.GetServicenowCMDBDeployment(jobStatus.Links.ManagedObject.Href)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(snowDeployment.ID))

	return resourceServicenowCMDBDeploymentRead(ctx, d, m)
}

func resourceServicenowCMDBDeploymentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*cbclient.CloudBoltClient)
	var diags diag.Diagnostics

	snowDeployment, err := apiClient.GetServicenowCMDBDeploymentById(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("workspace_url", snowDeployment.Links.Workspace.Href); err != nil {
		diags = append(diags, diag.Errorf("Cannot set workspace: %s", snowDeployment.Links.Workspace.Href)...)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	if err := d.Set("configuration_items_info", snowDeployment.ConfigurationItemsInfo); err != nil {
		diags = append(diags, diag.Errorf("Cannot set configuration_items_info")...)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	executionDetailsJSON, err := json.Marshal(snowDeployment.ExecutionDetails)
	if err != nil {
		diags = append(diags, diag.Errorf("Unable to Marshal execution_details into string")...)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}
	executionDetailsString := string(executionDetailsJSON)

	if err := d.Set("execution_details", executionDetailsString); err != nil {
		diags = append(diags, diag.Errorf("Cannot set execution_details")...)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	servicenowCMDBPolicyURLSplit := strings.Split(snowDeployment.Links.Policy.Href, "/")
	servicenowCMDBPolicyID := servicenowCMDBPolicyURLSplit[len(servicenowCMDBPolicyURLSplit)-2]
	servicenowCMDBPolicyIDInt, _ := strconv.Atoi(servicenowCMDBPolicyID)
	if err := d.Set("policy_id", servicenowCMDBPolicyIDInt); err != nil {
		diags = append(diags, diag.Errorf("Cannot set policy")...)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	return nil
}

func resourceServicenowCMDBDeploymentUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("onefuse.resourceServicenowCMDBDeploymentUpdate")
	log.Println("No Op!")
	return nil
}

func resourceServicenowCMDBDeploymentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*cbclient.CloudBoltClient)

	jobStatus, err := apiClient.DeleteServicenowCMDBDeployment(d.Id())
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
