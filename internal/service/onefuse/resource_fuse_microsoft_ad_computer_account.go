package onefuse

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/cloudboltsoftware/cloudbolt-go-sdk/cbclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceMicrosoftADComputerAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMicrosoftADComputerAccountCreate,
		ReadContext:   resourceMicrosoftADComputerAccountRead,
		UpdateContext: resourceMicrosoftADComputerAccountUpdate,
		DeleteContext: resourceMicrosoftADComputerAccountDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Computer Account Name.",
				// Updates not yet supported for Microsoft Active Directory Computer Names.
				ForceNew: true,
				// Suppress diff if both names are the same in Lowercase or Uppercase
				DiffSuppressFunc: func(k string, oldName string, newName string, d *schema.ResourceData) bool {
					if strings.ToLower(oldName) == strings.ToLower(newName) {
						return true
					} else if strings.ToUpper(oldName) == strings.ToUpper(newName) {
						return true
					} else {
						return false
					}
				},
			},
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
			"final_ou": {
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

func resourceMicrosoftADComputerAccountCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*cbclient.CloudBoltClient)

	newComputerAccount := cbclient.MicrosoftADComputerAccount{
		Name:               d.Get("name").(string),
		FinalOU:            d.Get("final_ou").(string),
		PolicyID:           d.Get("policy_id").(int),
		WorkspaceURL:       d.Get("workspace_url").(string),
		TemplateProperties: d.Get("template_properties").(map[string]interface{}),
	}

	jobStatus, err := apiClient.CreateMicrosoftADComputerAccount(&newComputerAccount)
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

	var computerAccount *cbclient.MicrosoftADComputerAccount
	computerAccount, err = apiClient.GetMicrosoftADComputerAccount(jobStatus.Links.ManagedObject.Href)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(computerAccount.ID))

	return resourceDNSReservationRead(ctx, d, m)
}

func resourceMicrosoftADComputerAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*cbclient.CloudBoltClient)
	var diags diag.Diagnostics

	computerAccount, err := apiClient.GetMicrosoftADComputerAccountById(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", computerAccount.Name); err != nil {
		diags = append(diags, diag.Errorf("Cannot set name: %s", computerAccount.Name)...)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	if err := d.Set("final_ou", computerAccount.FinalOU); err != nil {
		diags = append(diags, diag.Errorf("Cannot set final OU: %s", computerAccount.FinalOU)...)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	if err := d.Set("workspace_url", computerAccount.Links.Workspace.Href); err != nil {
		diags = append(diags, diag.Errorf("Cannot set workspace: %s", computerAccount.Links.Workspace.Href)...)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	microsoftADPolicyURLSplit := strings.Split(computerAccount.Links.Policy.Href, "/")
	microsoftADPolicyID := microsoftADPolicyURLSplit[len(microsoftADPolicyURLSplit)-2]
	microsoftADPolicyIDInt, _ := strconv.Atoi(microsoftADPolicyID)
	if err := d.Set("policy_id", microsoftADPolicyIDInt); err != nil {
		diags = append(diags, diag.Errorf("Cannot set policy")...)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	return nil
}

func resourceMicrosoftADComputerAccountUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("onefuse.resourceMicrosoftADComputerAccountUpdate")
	log.Println("No Op!")
	return nil
}

func resourceMicrosoftADComputerAccountDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*cbclient.CloudBoltClient)

	jobStatus, err := apiClient.DeleteMicrosoftADComputerAccount(d.Id())
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
