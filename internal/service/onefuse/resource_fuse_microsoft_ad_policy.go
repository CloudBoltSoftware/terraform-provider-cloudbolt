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

func ResourceMicrosoftADPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMicrosoftADPolicyCreate,
		ReadContext:   resourceMicrosoftADPolicyRead,
		UpdateContext: resourceMicrosoftADPolicyUpdate,
		DeleteContext: resourceMicrosoftADPolicyDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"microsoft_endpoint_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"computer_name_letter_case": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Either Lowercase or Uppercase",
			},
			"ou": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"workspace_url": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"security_groups": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"create_ou": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"remove_ou": {
				Type:     schema.TypeBool,
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

func resourceMicrosoftADPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*cbclient.CloudBoltClient)

	var securityGroups []string
	for _, group := range d.Get("security_groups").([]interface{}) {
		securityGroups = append(securityGroups, group.(string))
	}

	newPolicy := cbclient.MicrosoftADPolicy{
		Name:                   d.Get("name").(string),
		Description:            d.Get("description").(string),
		OU:                     d.Get("ou").(string),
		MicrosoftEndpointID:    d.Get("microsoft_endpoint_id").(int),
		ComputerNameLetterCase: d.Get("computer_name_letter_case").(string),
		WorkspaceURL:           d.Get("workspace_url").(string),
		CreateOU:               d.Get("create_ou").(bool),
		RemoveOU:               d.Get("remove_ou").(bool),
		SecurityGroups:         securityGroups,
	}

	adPolicy, err := apiClient.CreateMicrosoftADPolicy(&newPolicy)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(adPolicy.ID))

	return resourceMicrosoftADPolicyRead(ctx, d, m)
}

func resourceMicrosoftADPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*cbclient.CloudBoltClient)
	var diags diag.Diagnostics

	policy, err := apiClient.GetMicrosoftADPolicyByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", policy.Name); err != nil {
		diags = append(diags, diag.Errorf("Cannot set name: %s", policy.Name)...)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	if err := d.Set("description", policy.Description); err != nil {
		diags = append(diags, diag.Errorf("Cannot set description: '%s'", policy.Description)...)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	if err := d.Set("workspace_url", policy.Links.Workspace.Href); err != nil {
		diags = append(diags, diag.Errorf("Cannot set workspace: '%s'", policy.Links.Workspace.Href)...)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	if err := d.Set("computer_name_letter_case", policy.ComputerNameLetterCase); err != nil {
		diags = append(diags, diag.Errorf("Cannot set computer_name_letter_case: '%s'", policy.ComputerNameLetterCase)...)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	if err := d.Set("ou", policy.OU); err != nil {
		diags = append(diags, diag.Errorf("Cannot set OU: '%s'", policy.OU)...)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	if err := d.Set("create_ou", policy.CreateOU); err != nil {
		diags = append(diags, diag.Errorf("Cannot set Create OU: %t", policy.CreateOU)...)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	if err := d.Set("remove_ou", policy.RemoveOU); err != nil {
		diags = append(diags, diag.Errorf("Cannot set Remove OU: %t", policy.RemoveOU)...)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	if err := d.Set("security_groups", policy.SecurityGroups); err != nil {
		diags = append(diags, diag.Errorf("Cannot set Security Groups: %#v", policy.SecurityGroups)...)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	microsoftEndpointURLSplit := strings.Split(policy.Links.MicrosoftEndpoint.Href, "/")
	microsoftEndpointID := microsoftEndpointURLSplit[len(microsoftEndpointURLSplit)-2]
	microsoftEndpointIDInt, err := strconv.Atoi(microsoftEndpointID)
	if err != nil {
		diags = append(diags, diag.Errorf("Expected to convert '%s' to int value.", microsoftEndpointID)...)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}
	if err := d.Set("microsoft_endpoint_id", microsoftEndpointIDInt); err != nil {
		diags = append(diags, diag.Errorf("Cannot set microsoft_endpoint_id")...)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	return nil
}

func resourceMicrosoftADPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("onefuse.resourceMicrosoftADPolicyUpdate")

	// Determine if a change is needed
	changed := (d.HasChange("name") ||
		d.HasChange("description") ||
		d.HasChange("microsoft_endpoint_id") ||
		d.HasChange("computer_name_letter_case") ||
		d.HasChange("workspace_url") ||
		d.HasChange("ou") ||
		d.HasChange("create_ou") ||
		d.HasChange("remove_ou") ||
		d.HasChange("security_groups"))

	if !changed {
		return nil
	}

	apiClient := m.(*cbclient.CloudBoltClient)

	var securityGroups []string
	for _, group := range d.Get("security_groups").([]interface{}) {
		securityGroups = append(securityGroups, group.(string))
	}

	desiredPolicy := cbclient.MicrosoftADPolicy{
		Name:                   d.Get("name").(string),
		Description:            d.Get("description").(string),
		MicrosoftEndpointID:    d.Get("microsoft_endpoint_id").(int),
		ComputerNameLetterCase: d.Get("computer_name_letter_case").(string),
		WorkspaceURL:           d.Get("workspace_url").(string),
		OU:                     d.Get("ou").(string),
		CreateOU:               d.Get("create_ou").(bool),
		RemoveOU:               d.Get("remove_ou").(bool),
		SecurityGroups:         securityGroups,
	}

	_, err := apiClient.UpdateMicrosoftADPolicy(d.Id(), &desiredPolicy)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceMicrosoftADPolicyRead(ctx, d, m)
}

func resourceMicrosoftADPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*cbclient.CloudBoltClient)

	err := apiClient.DeleteMicrosoftADPolicy(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
