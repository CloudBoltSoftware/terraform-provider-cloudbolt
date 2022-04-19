package cmp

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cloudboltsoftware/cloudbolt-go-sdk/cbclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceBPInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBPInstanceCreate,
		ReadContext:   resourceBPInstanceRead,
		UpdateContext: resourceBPInstanceUpdate,
		DeleteContext: resourceBPInstanceDelete,

		Schema: map[string]*schema.Schema{
			"group": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The relative API URL path for the CloudBolt Group",
			},
			"blueprint_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The global Id for the CloudBolt Blueprint",
			},
			"parameters": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Parameters Name/Value pair",
			},
			"resource_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name for the created CloudBolt Resoucce",
			},
			"deployment_item": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "Set of blueprint deployment items",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The reference name for the blueprint deployment item",
						},
						"environment": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The relative API URL path for the CloudBolt Environment",
						},
						"osbuild": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The relative API URL path for the CloudBolt OS Build",
						},
						"parameters": {
							Type:        schema.TypeMap,
							Required:    true,
							Description: "Parameter Name/Value pair",
						},
					},
				},
			},
			"servers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hostname": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Server Hostname",
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
				},
			},
			"server_hostname": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Server Hostname",
			},
			"server_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Server IP Address",
			},
			"instance_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of deployed instance Resource or Server",
			},
			"attributes": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Description: "CloudBolt Resource attributes",
			},
		},
	}
}

func resourceBPInstanceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*cbclient.CloudBoltClient)

	bpItems := make([]map[string]interface{}, 0)
	bpItemList := d.Get("deployment_item").(*schema.Set).List()
	bpParams := d.Get("parameters").(map[string]interface{})
	for _, v := range bpItemList {
		m := v.(map[string]interface{})
		bpItem := map[string]interface{}{
			"bp-item-name":    m["name"].(string),
			"bp-item-paramas": m["parameters"].(map[string]interface{}),
		}

		env, ok := m["environment"]
		if ok && env != "" {
			bpItem["environment"] = env.(string)
		}

		osb, ok := m["osbuild"]
		if ok && osb != "" {
			bpItem["osbuild"] = osb.(string)
		}

		bpItems = append(bpItems, bpItem)
	}

	order, err := apiClient.DeployBlueprint(d.Get("group").(string), d.Get("blueprint_id").(string), d.Get("resource_name").(string), bpParams, bpItems)
	if err != nil {
		return diag.FromErr(err)
	}

	stateChangeConf := resource.StateChangeConf{
		Delay:   10 * time.Second,
		Timeout: 10 * time.Minute,
		Pending: []string{"ACTIVE"},
		Target:  []string{"SUCCESS"},
		Refresh: OrderStateRefreshFunc(apiClient, order.ID),
	}

	_, err = stateChangeConf.WaitForState()
	if err != nil {
		return diag.Errorf("Error waiting for Order (%s) to complete. Error: %s", order.ID, err)
	}

	// Retrieve the updated order to obtain Resource ID
	order, err = apiClient.GetOrder(order.ID)
	if err != nil {
		return diag.FromErr(err)
	}

	var resourceId string
	var servers []string = make([]string, 0)
	for _, j := range order.Links.Jobs {
		job, joberr := apiClient.GetJob(j.Href)
		if joberr != nil {
			return diag.FromErr(joberr)
		}

		if job.Type == "deploy_blueprint" {
			if len(job.Links.Resource.Href) > 0 {
				resourceId = job.Links.Resource.Href
				d.Set("instance_type", "Resource")
			} else if len(job.Links.Servers) > 0 {
				for _, s := range job.Links.Servers {
					serverHref := strings.TrimRight(s.Href, "/")
					index := strings.LastIndex(serverHref, "/")
					servers = append(servers, serverHref[index+1:])
				}
				d.Set("instance_type", "Server")
			}
			break
		}
	}

	if resourceId == "" && len(servers) == 0 {
		return diag.Errorf("Error Order (%s) does not have a Resource or Server", order.ID)
	}

	if resourceId != "" {
		d.SetId(resourceId)
	} else {
		d.SetId(strings.Join(servers, "_"))
	}

	return resourceBPInstanceRead(ctx, d, m)
}

func parseAttributes(attributes []map[string]interface{}) (map[string]interface{}, error) {
	resAttributes := make(map[string]interface{}, 0)

	for _, attr := range attributes {
		attrType, _ := attr["type"].(string)
		attrName, _ := attr["name"].(string)

		switch attrType {
		case "BOOL":
			attrValue, _ := attr["value"].(bool)
			resAttributes[attrName] = strconv.FormatBool(attrValue)
		case "DEC", "INT":
			attrValue, _ := attr["value"].(float64)
			resAttributes[attrName] = fmt.Sprintf("%g", attrValue)
		default:
			attrValue, _ := attr["value"].(string)
			resAttributes[attrName] = attrValue
		}
	}

	return resAttributes, nil
}

func parseServer(svr *cbclient.CloudBoltServer) (map[string]interface{}, error) {
	server := map[string]interface{}{
		"hostname":                svr.Hostname,
		"ip_address":              svr.IP,
		"status":                  svr.Status,
		"mac":                     svr.Mac,
		"date_added_to_cloudbolt": svr.DateAddedToCloudbolt,
		"cpu_count":               svr.CPUCount,
		"memory_size_gb":          svr.MemorySizeGB,
		"disk_size_gb":            svr.DiskSizeGB,
	}

	if svr.PowerStatus != "" {
		server["power_status"] = svr.PowerStatus
	}

	if svr.Notes != "" {
		server["notes"] = svr.Notes
	}

	if svr.Labels != nil {
		server["labels"] = svr.Labels
	}

	if svr.OsFamily != "" {
		server["os_family"] = svr.OsFamily
	}

	if svr.RateBreakdown != nil {
		server["rate_breakdown"] = svr.RateBreakdown
	}

	if len(svr.Disks) > 0 {
		disks := make([]map[string]interface{}, 0)

		for _, d := range svr.Disks {
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

		server["disks"] = disks
	}

	if len(svr.Networks) > 0 {
		server["networks"] = svr.Networks
	}

	if len(svr.TechSpecificAttributes) > 0 {
		server["tech_specific_attributes"] = svr.TechSpecificAttributes
	}

	svrAttributes, _ := parseAttributes(svr.Attributes)
	server["attributes"] = svrAttributes

	return server, nil
}

func resourceBPInstanceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*cbclient.CloudBoltClient)
	instanceType := d.Get("instance_type").(string)

	if instanceType == "Resource" {
		res, err := apiClient.GetResource(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		var servers []map[string]interface{}
		for _, s := range res.Links.Servers {
			svr, svrerr := apiClient.GetServer(s.Href)
			if svrerr != nil {
				return diag.Errorf("Error getting Servers for Resource: %s", svrerr)
			}

			server, _ := parseServer(svr)
			servers = append(servers, server)
		}

		if servers != nil {
			d.Set("servers", servers)
		}

		resAttributes, _ := parseAttributes(res.Attributes)
		d.Set("attributes", resAttributes)
	} else {
		serverIds := strings.Split(d.Id(), "_")

		servers := make([]map[string]interface{}, 0)
		for _, serverId := range serverIds {
			svr, svrerr := apiClient.GetServerById(serverId)
			if svrerr != nil {
				return diag.Errorf("Error getting Server: %s", svrerr)
			}

			server, _ := parseServer(svr)
			servers = append(servers, server)
		}

		if servers != nil {
			d.Set("servers", servers)
		}
	}

	return nil
}

func resourceBPInstanceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceBPInstanceRead(ctx, d, m)
}

func resourceBPInstanceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*cbclient.CloudBoltClient)
	instanceType := d.Get("instance_type").(string)

	if instanceType == "Resource" {
		res, err := apiClient.GetResource(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		var delActionPath string
		for _, v := range res.Links.Actions {
			if v.Title == "Delete" {
				delActionPath = v.Href
				break
			}
		}

		if delActionPath == "" {
			return diag.Errorf("Error deleting resource (%s).", d.Id())
		}

		job, delerr := apiClient.SubmitAction(delActionPath, res.Links.Self.Href)
		if delerr != nil {
			return diag.FromErr(delerr)
		}

		stateChangeConf := resource.StateChangeConf{
			Delay:   10 * time.Second,
			Timeout: 5 * time.Minute,
			Pending: []string{"INIT", "QUEUED", "PENDING", "RUNNING", "TO_CANCEL"},
			Target:  []string{"SUCCESS"},
			Refresh: JobStateRefreshFunc(apiClient, job.Links.Self.Href),
		}

		_, err = stateChangeConf.WaitForState()
		if err != nil {
			return diag.Errorf("Error waiting for Job (%s) to complete: %s", job.Links.Self.Href, err)
		}
	} else {
		serverIds := strings.Split(d.Id(), "_")

		for _, serverId := range serverIds {
			decomResult, err := apiClient.DecomServer(serverId)
			if err != nil {
				return diag.FromErr(err)
			}

			var stateChangeConf resource.StateChangeConf
			stateChangeConf.Delay = 10 * time.Second
			stateChangeConf.Timeout = 5 * time.Minute
			stateChangeConf.Target = []string{"SUCCESS"}
			if strings.HasPrefix(decomResult.ID, "ORD-") {
				stateChangeConf.Pending = []string{"ACTIVE"}
				stateChangeConf.Refresh = OrderStateRefreshFunc(apiClient, decomResult.ID)
			} else {
				stateChangeConf.Pending = []string{"INIT", "QUEUED", "PENDING", "RUNNING", "TO_CANCEL"}
				stateChangeConf.Refresh = JobStateRefreshFunc(apiClient, decomResult.Links.Self.Href)
			}

			_, err = stateChangeConf.WaitForState()
			if err != nil {
				return diag.Errorf("Error waiting for Decom Server (%s) to complete: %s", decomResult.Links.Self.Href, err)
			}
		}
	}

	return nil
}

func OrderStateRefreshFunc(apiClient *cbclient.CloudBoltClient, orderId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		order, err := apiClient.GetOrder(orderId)

		if err != nil {
			return nil, "ERROR", err
		}

		if order.Status == "FAILURE" {
			return nil, order.Status, fmt.Errorf("Order %s failed to reach target state.", orderId)
		}

		return order, order.Status, nil
	}
}

func JobStateRefreshFunc(apiClient *cbclient.CloudBoltClient, jobPath string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		job, err := apiClient.GetJob(jobPath)
		if err != nil {
			return nil, "", err
		}

		if job.Status == "FAILURE" || job.Status == "WARNING" || job.Status == "CANCELED" {
			return nil, job.Status, fmt.Errorf("Job %s failed to reach target state.", jobPath)
		}

		return job, job.Status, nil
	}
}
