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

func ResourceDNSReservation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDNSReservationCreate,
		ReadContext:   resourceDNSReservationRead,
		UpdateContext: resourceDNSReservationUpdate,
		DeleteContext: resourceDNSReservationDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
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
			"value": {
				Type:     schema.TypeString,
				Required: true,
			},
			"zones": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
			},
			"template_properties": {
				Type:     schema.TypeMap,
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

func resourceDNSReservationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*cbclient.CloudBoltClient)

	var dnsZones []string
	for _, group := range d.Get("zones").([]interface{}) {
		dnsZones = append(dnsZones, group.(string))
	}

	newDNSRecord := cbclient.DNSReservation{
		Name:               d.Get("name").(string),
		PolicyID:           d.Get("policy_id").(int),
		WorkspaceURL:       d.Get("workspace_url").(string),
		Value:              d.Get("value").(string),
		Zones:              dnsZones,
		TemplateProperties: d.Get("template_properties").(map[string]interface{}),
	}

	jobStatus, err := apiClient.CreateDNSReservation(&newDNSRecord)
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

	var dnsRecord *cbclient.DNSReservation
	dnsRecord, err = apiClient.GetDNSReservation(jobStatus.Links.ManagedObject.Href)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(dnsRecord.ID))

	return resourceDNSReservationRead(ctx, d, m)
}

func resourceDNSReservationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*cbclient.CloudBoltClient)
	var diags diag.Diagnostics

	dnsRecord, err := apiClient.GetDNSReservationById(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", dnsRecord.Name); err != nil {
		diags = append(diags, diag.Errorf("Cannot set name: %s", dnsRecord.Name)...)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	if err := d.Set("workspace_url", dnsRecord.Links.Workspace.Href); err != nil {
		diags = append(diags, diag.Errorf("Cannot set workspace: %s", dnsRecord.Links.Workspace.Href)...)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	dnsPolicyURLSplit := strings.Split(dnsRecord.Links.Policy.Href, "/")
	dnsPolicyID := dnsPolicyURLSplit[len(dnsPolicyURLSplit)-2]
	dnsPolicyIDInt, _ := strconv.Atoi(dnsPolicyID)
	if err := d.Set("policy_id", dnsPolicyIDInt); err != nil {
		diags = append(diags, diag.Errorf("Cannot set policy")...)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	return nil
}

func resourceDNSReservationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("onefuse.resourceDNSReservationUpdate")
	log.Println("No Op!")
	return nil
}

func resourceDNSReservationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*cbclient.CloudBoltClient)

	jobStatus, err := apiClient.DeleteDNSReservation(d.Id())
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
