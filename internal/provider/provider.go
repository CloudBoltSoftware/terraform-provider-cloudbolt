package provider

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/cloudboltsoftware/cloudbolt-go-sdk/cbclient"
	"github.com/cloudboltsoftware/terraform-provider-cloudbolt/internal/service/cmp"
	"github.com/cloudboltsoftware/terraform-provider-cloudbolt/internal/service/onefuse"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"cb_protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "https",
				Description:  "CloudBolt API Protocol,  Default (https)",
				ValidateFunc: checkProtocol,
			},
			"cb_host": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "CloudBolt API Host",
			},
			"cb_port": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "443",
				Description: "CloudBolt API Port, Default (443)",
			},
			"cb_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     10,
				Description: "Timeout in seconds, Default (10)",
			},
			"cb_insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Disable SSL Verification, Default (true)",
			},
			"cb_username": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "CloudBolt API Username, required if not provided in environment variable CB_USERNAME",
				DefaultFunc:  schema.EnvDefaultFunc("CB_USERNAME", nil),
				ValidateFunc: checkNotEmptyString,
			},
			"cb_password": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				Description:  "CloudBolt API Password, required if not provided in environment variable CB_PASSWORD",
				DefaultFunc:  schema.EnvDefaultFunc("CB_PASSWORD", nil),
				ValidateFunc: checkNotEmptyString,
			},
			"cb_domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "CloudBolt API Domain, can also be set using environment variable CB_DOMAIN",
				DefaultFunc: schema.EnvDefaultFunc("CB_DOMAIN", ""),
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"cloudbolt_group_ref":                 cmp.DataSourceCloudBoltGroup(),
			"cloudbolt_blueprint_ref":             cmp.DataSourceCloudBoltBlueprint(),
			"cloudbolt_environment_ref":           cmp.DataSourceCloudBoltEnvironment(),
			"cloudbolt_osbuild_ref":               cmp.DataSourceCloudBoltOSBuild(),
			"cloudbolt_resource_handler_ref":      cmp.DataSourceCloudBoltResourceHandler(),
			"cloudbolt_1f_ad_policy":              onefuse.DataSourceADPolicy(),
			"cloudbolt_1f_ansible_tower_policy":   onefuse.DataSourceAnsibleTowerPolicy(),
			"cloudbolt_1f_dns_policy":             onefuse.DataSourceDNSPolicy(),
			"cloudbolt_1f_ipam_policy":            onefuse.DataSourceIPAMPolicy(),
			"cloudbolt_1f_module_policy":          onefuse.DataSourceModulePolicy(),
			"cloudbolt_1f_naming_policy":          onefuse.DataSourceNamingPolicy(),
			"cloudbolt_1f_scripting_policy":       onefuse.DataSourceScriptingPolicy(),
			"cloudbolt_1f_servicenow_cmdb_policy": onefuse.DataSourceServiceNowCMDBPolicy(),
			"cloudbolt_1f_vra_policy":             onefuse.DataSourceVraPolicy(),
			"cloudbolt_1f_rendered_template":      onefuse.DataSourceRenderedTemplate(),
			"cloudbolt_1f_static_property_set":    onefuse.DataSourceStaticPropertySet(),
			"cloudbolt_1f_microsoft_endpoint":     onefuse.DataSourceMicrosoftEndpoint(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"cloudbolt_bp_instance":                      cmp.ResourceBPInstance(),
			"cloudbolt_1f_module_deployment":             onefuse.ResourceModuleDeployment(),
			"cloudbolt_1f_ansible_tower_deployment":      onefuse.ResourceAnsibleTowerDeployment(),
			"cloudbolt_1f_dns_record":                    onefuse.ResourceDNSReservation(),
			"cloudbolt_1f_ipam_record":                   onefuse.ResourceIPAMReservation(),
			"cloudbolt_1f_naming":                        onefuse.ResourceCustomNaming(),
			"cloudbolt_1f_microsoft_ad_policy":           onefuse.ResourceMicrosoftADPolicy(),
			"cloudbolt_1f_microsoft_ad_computer_account": onefuse.ResourceMicrosoftADComputerAccount(),
			"cloudbolt_1f_scripting_deployment":          onefuse.ResourceScriptingDeployment(),
			"cloudbolt_1f_vra_deployment":                onefuse.ResourceVraDeployment(),
			"cloudbolt_1f_servicenow_cmdb_deployment":    onefuse.ResourceServicenowCMDBDeployment(),
		},

		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Need to validate CB_USERNAME enviornment variable is set if cb_username not in the config.
	if d.Get("cb_username").(string) == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "\"cb_username must be set or CB_USERNAME environment variable must set.",
		})
	}

	// Need to validate CB_PASSWORD enviornment variable is set if cb_password not in the config.
	if d.Get("cb_password").(string) == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "\"cb_password must be set or CB_PASSWORD environment variable must set.",
		})
	}

	if diags != nil {
		return nil, diags
	}

	httpClient := &http.Client{
		// (Optional) User requested insecure transport
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: d.Get("cb_insecure").(bool), // Default: false
			},
		},
		// (Optional) User requested timeout
		Timeout: time.Duration(d.Get("cb_timeout").(int)) * time.Second, // Default: 10 seconds
	}

	apiClient := cbclient.New(
		d.Get("cb_protocol").(string),
		d.Get("cb_host").(string),
		d.Get("cb_port").(string),
		d.Get("cb_username").(string),
		d.Get("cb_password").(string),
		d.Get("cb_domain").(string),
		httpClient,
	)

	_, err := apiClient.Authenticate()
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return apiClient, diags
}

func checkNotEmptyString(val interface{}, key string) (warns []string, errs []error) {
	v := val.(string)
	if v == "" {
		errs = append(errs, fmt.Errorf("%q is required and most not be empty.", key))

		return warns, errs
	}

	return warns, errs
}

func checkProtocol(val interface{}, key string) (warns []string, errs []error) {
	v := val.(string)
	if v != "http" && v != "https" {
		errs = append(errs, fmt.Errorf("%q must be either \"https\" or \"http\".", key))

		return warns, errs
	}

	return warns, errs
}
