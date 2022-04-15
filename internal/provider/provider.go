package provider

import (
	"context"
	"crypto/tls"
	"net/http"
	"time"

	"github.com/cloudboltsoftware/cloudbolt-go-sdk/cbclient"
	"github.com/cloudboltsoftware/terraform-provider-cloudbolt/internal/service/cmp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"cb_protocol": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "CloudBolt API Protocol",
			},
			"cb_host": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "CloudBolt API Host",
			},
			"cb_port": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "CloudBolt API Port",
			},
			"cb_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Timeout in seconds",
			},
			"cb_insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Disable SSL Verification",
			},
			"cb_username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "CloudBolt API Username",
			},
			"cb_password": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "CloudBolt API Password",
			},
			"cb_domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "CloudBolt API Domain",
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"cloudbolt_group_ref":       cmp.DataSourceCloudBoltGroup(),
			"cloudbolt_blueprint_ref":   cmp.DataSourceCloudBoltBlueprint(),
			"cloudbolt_environment_ref": cmp.DataSourceCloudBoltEnvironment(),
			"cloudbolt_osbuild_ref":     cmp.DataSourceCloudBoltOSBuild(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"cloudbolt_bp_instance": cmp.ResourceBPInstance(),
		},

		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	// Default HTTP timeout to 10 seconds
	if d.Get("cb_timeout").(int) <= 0 {
		d.Set("cb_timeout", 10)
	}

	// Default Protocol is HTTPS
	if d.Get("cb_protocol").(string) == "" {
		d.Set("cb_protocol", "https")
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

	return apiClient, nil
}
