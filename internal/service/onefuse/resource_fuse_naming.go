package onefuse

import (
	"context"
	"log"
	"strconv"

	"github.com/cloudboltsoftware/cloudbolt-go-sdk/cbclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceCustomNaming() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCustomNameCreate,
		ReadContext:   resourceCustomNameRead,
		UpdateContext: resourceCustomNameUpdate,
		DeleteContext: resourceCustomNameDelete,
		Schema: map[string]*schema.Schema{
			"custom_name_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"naming_policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "OneFuse Module Policy ID.",
			},
			"dns_suffix": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "DNS Suffix to append to the Hostname.",
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "OneFuse Workspace URL path.",
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

func resourceCustomNameCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*cbclient.CloudBoltClient)

	namingPolicyID := d.Get("naming_policy_id").(string)
	workspaceID := d.Get("workspace_id").(string)
	templateProperties := d.Get("template_properties").(map[string]interface{})

	jobStatus, err := apiClient.GenerateCustomName(namingPolicyID, workspaceID, templateProperties)
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

	var customName *cbclient.CustomName
	customName, err = apiClient.GetCustomName(jobStatus.Links.ManagedObject.Href)
	if err != nil {
		return diag.FromErr(err)
	}

	// setting the ID is REALLY necessary here
	// we use the FQDN instead of the numeric ID as it is more likely to remain consistent as a composite key in TF
	d.SetId(customName.Name + "." + customName.DnsSuffix)

	if err := d.Set("custom_name_id", customName.Id); err != nil {
		var diags diag.Diagnostics

		diags = append(diags, diag.Errorf("cannot set custom_name_id")...)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	return resourceCustomNameRead(ctx, d, m)
}

func resourceCustomNameRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*cbclient.CloudBoltClient)
	var diags diag.Diagnostics

	customNameId := strconv.Itoa(d.Get("custom_name_id").(int))
	customName, err := apiClient.GetCustomNameById(customNameId)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", customName.Name); err != nil {
		diags = append(diags, diag.Errorf("Cannot set name: %s", customName.Name)...)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}
	if err := d.Set("dns_suffix", customName.DnsSuffix); err != nil {
		diags = append(diags, diag.Errorf("Cannot set dns_suffix: %s", customName.DnsSuffix)...)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	return nil
}

func resourceCustomNameUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("onefuse.resourceCustomNameUpdate")
	log.Println("No Op!")
	return nil
}

func resourceCustomNameDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*cbclient.CloudBoltClient)

	customNameId := strconv.Itoa(d.Get("custom_name_id").(int))
	jobStatus, err := apiClient.DeleteCustomName(customNameId)
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
