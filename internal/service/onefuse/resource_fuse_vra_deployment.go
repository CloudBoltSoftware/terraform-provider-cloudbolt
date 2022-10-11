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

func ResourceVraDeployment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVraDeploymentCreate,
		ReadContext:   resourceVraDeploymentRead,
		UpdateContext: resourceVraDeploymentUpdate,
		DeleteContext: resourceVraDeploymentDelete,
		Schema: map[string]*schema.Schema{
			"policy_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"workspace_url": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"deployment_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"template_properties": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"deployment_info": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"blueprint_name": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"project_name": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
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

func resourceVraDeploymentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*cbclient.CloudBoltClient)

	newVraDeployment := cbclient.VraDeployment{
		PolicyID:           d.Get("policy_id").(int),
		WorkspaceURL:       d.Get("workspace_url").(string),
		DeploymentName:     d.Get("deployment_name").(string),
		TemplateProperties: d.Get("template_properties").(map[string]interface{}),
	}

	jobStatus, err := apiClient.CreateVraDeployment(&newVraDeployment)
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

	var vraDeployment *cbclient.VraDeployment
	vraDeployment, err = apiClient.GetVraDeployment(jobStatus.Links.ManagedObject.Href)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(vraDeployment.ID))

	return resourceVraDeploymentRead(ctx, d, m)
}

func resourceVraDeploymentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*cbclient.CloudBoltClient)
	var diags diag.Diagnostics

	vraDeployment, err := apiClient.GetVraDeploymentById(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("workspace_url", vraDeployment.Links.Workspace.Href); err != nil {
		diags = append(diags, diag.Errorf("Cannot set workspace: %s", vraDeployment.Links.Workspace.Href)...)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	if err := d.Set("deployment_name", vraDeployment.Name); err != nil {
		diags = append(diags, diag.Errorf("Cannot set deployment name: %s", vraDeployment.Name)...)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	deploymentInfoJSON, err := json.Marshal(vraDeployment.DeploymentInfo)
	if err != nil {
		diags = append(diags, diag.Errorf("Unable to Marshal deployment_info into string")...)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}
	deploymentInfoString := string(deploymentInfoJSON)
	if err := d.Set("deployment_info", deploymentInfoString); err != nil {
		diags = append(diags, diag.Errorf("Cannot set deployment_info: %s", deploymentInfoString)...)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	if err := d.Set("blueprint_name", vraDeployment.BlueprintName); err != nil {
		diags = append(diags, diag.Errorf("Cannot set blueprint name: %s", vraDeployment.BlueprintName)...)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	if err := d.Set("project_name", vraDeployment.ProjectName); err != nil {
		diags = append(diags, diag.Errorf("Cannot set blueprint name: %s", vraDeployment.ProjectName)...)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	vraPolicyURLSplit := strings.Split(vraDeployment.Links.Policy.Href, "/")
	vraPolicyID := vraPolicyURLSplit[len(vraPolicyURLSplit)-2]
	vraPolicyIDInt, _ := strconv.Atoi(vraPolicyID)
	if err := d.Set("policy_id", vraPolicyIDInt); err != nil {
		diags = append(diags, diag.Errorf("Cannot set policy")...)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	return nil
}

func resourceVraDeploymentUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("onefuse.resourceVraDeploymentUpdate")
	log.Println("No Op!")
	return nil
}

func resourceVraDeploymentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*cbclient.CloudBoltClient)

	jobStatus, err := apiClient.DeleteVraDeployment(d.Id())
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
