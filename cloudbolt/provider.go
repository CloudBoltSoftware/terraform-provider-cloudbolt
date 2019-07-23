package cloudbolt

import (
	"github.com/hashicorp/terraform/helper/schema"
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
		},

		DataSourcesMap: map[string]*schema.Resource{
			"cloudbolt_object_ref": dataSourceCloudBoltObject(),
			"cloudbolt_group_ref":  dataSourceCloudBoltGroup(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"cloudbolt_bp_instance": resourceBPInstance(),
		},

		ConfigureFunc: providerConfigure,
	}
}
