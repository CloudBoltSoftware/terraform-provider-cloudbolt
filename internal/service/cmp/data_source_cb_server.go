package cmp

import (
	"context"

	"github.com/cloudboltsoftware/cloudbolt-go-sdk/cbclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceCloudBoltServer() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCloudBoltServerRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The global id of a CloudBolt Server, required if \"hostname\" or \"url_path\" is not provided",
			},
			"url_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The relative API URL path for the CloudBolt Server, required if \"id\" or \"hostname\" is not provided",
			},
			"hostname": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The hostname of the CloudBolt Server, required if \"id\" or \"url_path\" is not provided",
			},
			"ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Server IP Address",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CloudBolt Server Status",
			},
			"mac": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Server MAC Address",
			},
			"power_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Server Power Status",
			},
			"date_added_to_cloudbolt": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date the server was added to CloudBolt",
			},
			"cpu_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "CPU Count",
			},
			"memory_size_gb": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Total Memory in GB",
			},
			"disk_size_gb": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total Disk Size in GB",
			},
			"notes": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Server Notes",
			},
			"labels": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Server Labels",
			},
			"os_family": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Server OS Family",
			},
			"attributes": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "CloudBolt Server attributes",
			},
			"rate_breakdown": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Server Rate Breakdown",
			},
			"tech_specific_attributes": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Resource Handler technical specific attributes",
			},
			"disks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Unique ID of Disk",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of Disk",
						},
						"disk_size_gb": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Disk Size in GB",
						},
					},
				},
				Description: "Server disks",
			},
			"networks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				Description: "Server NICs",
			},
		},
	}
}

func dataSourceCloudBoltServerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*cbclient.CloudBoltClient)
	serverPath := d.Get("url_path").(string)
	hostname := d.Get("hostname").(string)
	id := d.Get("id").(string)

	if id == "" && hostname == "" && serverPath == "" {
		return diag.Errorf("Either id, hostname, or url_path is required")
	}
	var server *cbclient.CloudBoltServer
	var err error
	if serverPath != "" {
		server, err = apiClient.GetServer(serverPath)
	} else if hostname != "" {
		server, err = apiClient.GetServerByHostname(hostname)
	} else {
		server, err = apiClient.GetServerById(id)
	}

	if err != nil {
		return diag.FromErr(err)
	}

	if server.ID == "" {
		return diag.Errorf("Server not found.")
	}

	d.SetId(server.ID)
	d.Set("url_path", server.Links.Self)
	d.Set("hostname", server.Hostname)
	d.Set("ip_address", server.IP)
	d.Set("status", server.Status)
	d.Set("mac", server.Mac)
	d.Set("date_added_to_cloudbolt", server.DateAddedToCloudbolt)
	d.Set("cpu_count", server.CPUCount)
	d.Set("memory_size_gb", server.MemorySizeGB)
	d.Set("disk_size_gb", server.DiskSizeGB)

	if server.PowerStatus != "" {
		d.Set("power_status", server.PowerStatus)
	}

	if server.Notes != "" {
		d.Set("notes", server.Notes)
	}

	if server.Labels != nil {
		d.Set("labels", server.Labels)
	}

	if server.OsFamily != "" {
		d.Set("os_family", server.OsFamily)
	}

	if server.RateBreakdown != nil {
		d.Set("rate_breakdown", server.RateBreakdown)
	}

	if len(server.Disks) > 0 {
		disks := make([]map[string]interface{}, 0)

		for _, d := range server.Disks {
			uuid, _ := d["uuid"]
			name, _ := d["name"]
			disk_size_gb, _ := d["diskSize"]

			disk := map[string]interface{}{
				"uuid":         uuid,
				"name":         name,
				"disk_size_gb": disk_size_gb,
			}
			disks = append(disks, disk)
		}

		d.Set("disks", disks)
	}

	if len(server.Networks) > 0 {
		d.Set("networks", server.Networks)
	}

	if len(server.TechSpecificAttributes) > 0 {
		d.Set("tech_specific_attributes", convertValuesToString(server.TechSpecificAttributes))
	}

	svrAttributes, _ := parseAttributes(server.Attributes)
	d.Set("attributes", svrAttributes)

	return nil
}
