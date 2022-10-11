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

func ResourceIPAMReservation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIPAMReservationCreate,
		ReadContext:   resourceIPAMReservationRead,
		UpdateContext: resourceIPAMReservationUpdate,
		DeleteContext: resourceIPAMReservationDelete,
		Schema: map[string]*schema.Schema{
			"hostname": {
				Type:     schema.TypeString,
				Required: true,
			},
			// hostname could potentially be overridden using the hostname override on the policy,
			// and therefore will no longer match the hostname given in the resource
			// so we need a different variable for the computed hostname
			"computed_hostname": {
				Type:     schema.TypeString,
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
			"ip_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"netmask": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"gateway": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"network": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"subnet": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"primary_dns": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"secondary_dns": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"nic_label": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"dns_suffix": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"dns_search_suffix": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
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

func resourceIPAMReservationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*cbclient.CloudBoltClient)

	var ipam_Suffixes []string
	for _, group := range d.Get("dns_search_suffix").([]interface{}) {
		ipam_Suffixes = append(ipam_Suffixes, group.(string))
	}

	newIPAMRecord := cbclient.IPAMReservation{
		Hostname:           d.Get("hostname").(string),
		PolicyID:           d.Get("policy_id").(int),
		WorkspaceURL:       d.Get("workspace_url").(string),
		IPaddress:          d.Get("ip_address").(string),
		Netmask:            d.Get("netmask").(string),
		Subnet:             d.Get("subnet").(string),
		Gateway:            d.Get("gateway").(string),
		Network:            d.Get("network").(string),
		PrimaryDNS:         d.Get("primary_dns").(string),
		SecondaryDNS:       d.Get("secondary_dns").(string),
		DNSSuffix:          d.Get("dns_suffix").(string),
		NicLabel:           d.Get("nic_label").(string),
		TemplateProperties: d.Get("template_properties").(map[string]interface{}),
	}
	jobStatus, err := apiClient.CreateIPAMReservation(&newIPAMRecord)
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

	var ipamRecord *cbclient.IPAMReservation
	ipamRecord, err = apiClient.GetIPAMReservation(jobStatus.Links.ManagedObject.Href)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(ipamRecord.ID))

	return resourceIPAMReservationRead(ctx, d, m)
}

func resourceIPAMReservationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*cbclient.CloudBoltClient)
	var diags diag.Diagnostics

	ipamRecord, err := apiClient.GetIPAMReservationById(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("computed_hostname", ipamRecord.Hostname); err != nil {
		diags = append(diags, diag.Errorf("Cannot set name: %s", ipamRecord.Hostname)...)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	if err := d.Set("workspace_url", ipamRecord.Links.Workspace.Href); err != nil {
		diags = append(diags, diag.Errorf("Cannot set workspace: %s", ipamRecord.Links.Workspace.Href)...)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	if err := d.Set("ip_address", ipamRecord.IPaddress); err != nil {
		diags = append(diags, diag.Errorf("Cannot set IPAddress: %s", ipamRecord.IPaddress)...)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	if err := d.Set("netmask", ipamRecord.Netmask); err != nil {
		diags = append(diags, diag.Errorf("Cannot set Netmask: %s", ipamRecord.Netmask)...)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	if err := d.Set("primary_dns", ipamRecord.PrimaryDNS); err != nil {
		diags = append(diags, diag.Errorf("Cannot set Primmary DNS: %s", ipamRecord.PrimaryDNS)...)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	if err := d.Set("secondary_dns", ipamRecord.SecondaryDNS); err != nil {
		diags = append(diags, diag.Errorf("Cannot set Secondary DNS: %s", ipamRecord.SecondaryDNS)...)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	if err := d.Set("gateway", ipamRecord.Gateway); err != nil {
		diags = append(diags, diag.Errorf("Cannot set Gateway: %s", ipamRecord.Gateway)...)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	if err := d.Set("network", ipamRecord.Network); err != nil {
		diags = append(diags, diag.Errorf("Cannot set Network: %s", ipamRecord.Network)...)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	if err := d.Set("subnet", ipamRecord.Subnet); err != nil {
		diags = append(diags, diag.Errorf("Cannot set Subnet: %s", ipamRecord.Subnet)...)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	if err := d.Set("nic_label", ipamRecord.NicLabel); err != nil {
		diags = append(diags, diag.Errorf("Cannot set NicLabel: %s", ipamRecord.NicLabel)...)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	if err := d.Set("dns_suffix", ipamRecord.DNSSuffix); err != nil {
		diags = append(diags, diag.Errorf("Cannot set DNSSuffix: %s", ipamRecord.DNSSuffix)...)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	ipamPolicyURLSplit := strings.Split(ipamRecord.Links.Policy.Href, "/")
	ipamPolicyID := ipamPolicyURLSplit[len(ipamPolicyURLSplit)-2]
	ipamPolicyIDInt, _ := strconv.Atoi(ipamPolicyID)
	if err := d.Set("policy_id", ipamPolicyIDInt); err != nil {
		diags = append(diags, diag.Errorf("Cannot set policy")...)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	return nil
}

func resourceIPAMReservationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("onefuse.resourceDNSReservationUpdate")
	log.Println("No Op!")
	return nil
}

func resourceIPAMReservationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*cbclient.CloudBoltClient)

	jobStatus, err := apiClient.DeleteIPAMReservation(d.Id())
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
