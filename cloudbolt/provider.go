package cloudbolt

import (
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
			// 	x"cloudbolt_object_ref": dataSourceCloudBoltObject(),
			"cloudbolt_group_ref":       dataSourceCloudBoltGroup(),
			"cloudbolt_blueprint_ref":   dataSourceCloudBoltBlueprint(),
			"cloudbolt_environment_ref": dataSourceCloudBoltEnvironment(),
			"cloudbolt_osbuild_ref":     dataSourceCloudBoltOSBuild(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"cloudbolt_bp_instance": resourceBPInstance(),
		},

		ConfigureFunc: providerConfigure,
	}
}
