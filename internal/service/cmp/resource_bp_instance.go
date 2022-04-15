package cmp

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/cloudboltsoftware/cloudbolt-go-sdk/cbclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceBPInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceBPInstanceCreate,
		Read:   resourceBPInstanceRead,
		Update: resourceBPInstanceUpdate,
		Delete: resourceBPInstanceDelete,

		Schema: map[string]*schema.Schema{
			"group": {
				Type:     schema.TypeString,
				Required: true,
			},
			"blueprint_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"parameters": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"resource_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"deployment_item": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"environment": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"osbuild": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"parameters": {
							Type:     schema.TypeMap,
							Required: true,
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
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mac": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"power_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"date_added_to_cloudbolt": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"memory_size_gb": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"disk_size_gb": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"notes": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"labels": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"os_family": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"attributes": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"rate_breakdown": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"tech_specific_attributes": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"disks": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"uuid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"disk_size_gb": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
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
						},
					},
				},
			},
			"server_hostname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"server_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"attributes": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
		},
	}
}

func resourceBPInstanceCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*cbclient.CloudBoltClient)

	// log.Printf("[!!] apiClient in resourceBPInstanceCreate: %+v", apiClient)

	bpItems := make([]map[string]interface{}, 0)
	bpItemList := d.Get("deployment_item").(*schema.Set).List()
	bpParams := d.Get("parameters").(map[string]interface{})
	for _, v := range bpItemList {
		m := v.(map[string]interface{})
		log.Printf("[!!] m in resourceBPInstanceCreate: %+v", m)

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
		return err
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
		return fmt.Errorf("Error waiting for Order (%s) to complete. Error: %s", order.ID, err)
	}

	// Retrieve the updated order to obtain Resource ID
	order, err = apiClient.GetOrder(order.ID)
	if err != nil {
		return err
	}

	var resourceId string
	var servers []string = make([]string, 0)
	for _, j := range order.Links.Jobs {
		job, joberr := apiClient.GetJob(j.Href)
		if joberr != nil {
			return joberr
		}

		if job.Type == "deploy_blueprint" {
			log.Println("[!!] Deploying Blueprint")
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
		return fmt.Errorf("Error Order (%s) does not have a Resource or Server", order.ID)
	}

	if resourceId != "" {
		d.SetId(resourceId)
	} else {
		d.SetId(strings.Join(servers, "_"))
	}

	return resourceBPInstanceRead(d, m)
}

func parseAttributes(attributes []map[string]interface{}) (map[string]interface{}, error) {
	resAttributes := make(map[string]interface{}, 0)

	for _, attr := range attributes {
		attrType, _ := attr["type"].(string)
		attrName, _ := attr["name"].(string)

		switch attrType {
		case "BOOL":
			log.Printf("resourceBPInstanceRead: BOOLEAN")
			attrValue, _ := attr["value"].(bool)
			resAttributes[attrName] = strconv.FormatBool(attrValue)
		case "DEC", "INT":
			log.Printf("resourceBPInstanceRead: FLOAT")
			attrValue, _ := attr["value"].(float64)
			resAttributes[attrName] = fmt.Sprintf("%g", attrValue)
			log.Printf("resourceBPInstanceRead: Converted Value - %+v", resAttributes[attrName])
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

func resourceBPInstanceRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*cbclient.CloudBoltClient)
	instanceType := d.Get("instance_type").(string)

	log.Printf("[!!] apiClient in resourceBPInstanceRead: %+v", apiClient)

	if instanceType == "Resource" {
		log.Printf("resourceBPInstanceRead: instanceType - %s", instanceType)
		res, err := apiClient.GetResource(d.Id())
		if err != nil {
			return err
		}

		log.Printf("resourceBPInstanceRead: Adding Servers")
		var servers []map[string]interface{}
		for _, s := range res.Links.Servers {
			svr, svrerr := apiClient.GetServer(s.Href)
			if svrerr != nil {
				return fmt.Errorf("Error getting Servers for Resource: %s", svrerr)
			}

			log.Printf("resourceBPInstanceRead: sbr - %+v", svr)

			server, _ := parseServer(svr)
			log.Printf("resourceBPInstanceRead: server - %+v", server)
			servers = append(servers, server)
		}

		if servers != nil {
			log.Printf("resourceBPInstanceRead: Set Server")
			log.Printf("resourceBPInstanceRead: server - %+v", servers)
			d.Set("servers", servers)
		}

		resAttributes, _ := parseAttributes(res.Attributes)

		log.Printf("resourceBPInstanceRead: resAttributes - %+v", resAttributes)
		d.Set("attributes", resAttributes)
	} else {
		serverIds := strings.Split(d.Id(), "_")

		servers := make([]map[string]interface{}, 0)
		for _, serverId := range serverIds {
			svr, svrerr := apiClient.GetServerById(serverId)
			if svrerr != nil {
				return fmt.Errorf("Error getting Server: %s", svrerr)
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

func resourceBPInstanceUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceBPInstanceRead(d, m)
}

func resourceBPInstanceDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*cbclient.CloudBoltClient)
	instanceType := d.Get("instance_type").(string)

	log.Printf("[!!] apiClient in resourceBPInstanceDelete: %+v", apiClient)

	if instanceType == "Resource" {
		res, err := apiClient.GetResource(d.Id())
		if err != nil {
			return err
		}

		log.Printf("Resource Result: %+v", res)
		var delActionPath string
		for _, v := range res.Links.Actions {
			log.Printf("Action Link: %+v", v)
			if v.Title == "Delete" {
				delActionPath = v.Href
				break
			}
		}

		if delActionPath == "" {
			return fmt.Errorf("Error deleting resource (%s).", d.Id())
		}

		job, delerr := apiClient.SubmitAction(delActionPath, res.Links.Self.Href)
		if delerr != nil {
			return delerr
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
			return fmt.Errorf("Error waiting for Job (%s) to complete: %s", job.Links.Self.Href, err)
		}
	} else {
		serverIds := strings.Split(d.Id(), "_")

		for _, serverId := range serverIds {
			decomResult, err := apiClient.DecomServer(serverId)
			if err != nil {
				return err
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
				return fmt.Errorf("Error waiting for Decom Server (%s) to complete: %s", decomResult.Links.Self.Href, err)
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
