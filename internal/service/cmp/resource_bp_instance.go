package cmp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"log"
	"runtime/debug"

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
			"request_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     30,
				Description: "Timeout in minutes, Default (30)",
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
							Optional:    true,
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
	var diags diag.Diagnostics
	defer withPanicRecovery(&diags, "Create")

	apiClient := m.(*cbclient.CloudBoltClient)

	bpItems := make([]map[string]interface{}, 0)
	bpItemList := d.Get("deployment_item").(*schema.Set).List()
	bpParams := normalizeParameters(d.Get("parameters").(map[string]interface{}))
	for _, v := range bpItemList {
		m := v.(map[string]interface{})
		itemParams := normalizeParameters(m["parameters"].(map[string]interface{}))
		bpItem := map[string]interface{}{
			"bp-item-name":    m["name"].(string),
			"bp-item-paramas": itemParams,
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

	requestTimeout := d.Get("request_timeout").(int)
	stateChangeConf := resource.StateChangeConf{
		Delay:   10 * time.Second,
		Timeout: time.Duration(requestTimeout) * time.Minute,
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

	// Populate Terraform state by reading the resource
	readDiags := resourceBPInstanceRead(ctx, d, m)
	diags = append(diags, readDiags...)

	return diags
}

func parseAttributes(attributes []map[string]interface{}) (map[string]interface{}, error) {
	resAttributes := make(map[string]interface{}, 0)

	for _, attr := range attributes {
		attrName, _ := attr["name"].(string)
		resAttributes[attrName] = convertValueToString(attr["value"])
	}

	return resAttributes, nil
}

func convertValueToString(value interface{}) string {
	var stringValue string

	boolValue, ok := value.(bool)
	if ok {
		stringValue = strconv.FormatBool(boolValue)
	}

	if stringValue == "" {
		intValue, ok := value.(int)
		if ok {
			stringValue = fmt.Sprint(intValue)
		}
	}

	if stringValue == "" {
		floatValue, ok := value.(float64)
		if ok {
			stringValue = fmt.Sprintf("%g", floatValue)
		}
	}

	if stringValue == "" {
		strValue, ok := value.(string)
		if ok {
			stringValue = strValue
		}
	}

	if stringValue == "" {
		interfaceArrValue, ok := value.([]interface{})
		if ok {
			var buffer bytes.Buffer
			buffer.WriteString("[")
			for i, v := range interfaceArrValue {
				if i != 0 {
					buffer.WriteString("|")
				}
				buffer.WriteString(convertValueToString(v))
			}
			buffer.WriteString("]")
			stringValue = buffer.String()
		}
	}

	return stringValue
}

func convertValuesToString(attributes map[string]interface{}) map[string]interface{} {
	stringValues := make(map[string]interface{}, 0)

	for k, v := range attributes {
		stringValue := convertValueToString(v)
		if stringValue != "" {
			stringValues[k] = convertValueToString(v)
		}
	}

	return stringValues
}

func getParameters(d *schema.ResourceData) (map[string]interface{}, bool) {
	rawParams, ok := d.GetOk("parameters")
	if !ok {
		return nil, false
	}

	params, ok := rawParams.(map[string]interface{})

	return params, ok
}

func setParameters(attributes map[string]interface{}, currentParameters map[string]interface{}) map[string]interface{} {
	parameters := make(map[string]interface{})
	for k, v := range currentParameters {
		if attributeValue, ok := attributes[k]; ok {
			parameters[k] = attributeValue
		} else {
			parameters[k] = convertValueToString(v)
		}
	}

	return parameters
}

func getDeploymentItems(d *schema.ResourceData) ([]map[string]interface{}, bool) {
	rawDepItems, ok := d.GetOk("deployment_item")
	if !ok {
		return nil, false
	}

	depSet, ok := rawDepItems.(*schema.Set)
	if !ok {
		return nil, false
	}

	var deploymentItems []map[string]interface{}
	for _, item := range depSet.List() {
		if itemMap, ok := item.(map[string]interface{}); ok {
			deploymentItems = append(deploymentItems, itemMap)
		}
	}

	return deploymentItems, true
}

func setDeploymentItems(attributes map[string]interface{}, currentDepItems []map[string]interface{}) []interface{} {
	var depItems []interface{}

	for _, currentDepItem := range currentDepItems {
		depItem := make(map[string]interface{})

		// Handle known fields, optionally override from `attributes`
		if name, ok := currentDepItem["name"].(string); ok {
			depItem["name"] = name
		}
		if env, ok := currentDepItem["environment"].(string); ok {
			depItem["environment"] = env
		}
		if osbuild, ok := currentDepItem["osbuild"].(string); ok {
			depItem["osbuild"] = osbuild
		}
		if params, ok := currentDepItem["parameters"].(map[string]interface{}); ok {
			depItem["parameters"] = setParameters(attributes, params)
		}

		depItems = append(depItems, depItem)
	}

	return depItems
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
		server["tech_specific_attributes"] = convertValuesToString(svr.TechSpecificAttributes)
	}

	svrAttributes, _ := parseAttributes(svr.Attributes)
	server["attributes"] = svrAttributes

	return server, nil
}

func resourceBPInstanceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	defer withPanicRecovery(&diags, "Read")

	apiClient := m.(*cbclient.CloudBoltClient)
	instanceType := d.Get("instance_type").(string)
	allAttributes := make(map[string]interface{})

	if instanceType == "Resource" {
		res, err := apiClient.GetResource(d.Id())
		if err != nil {
			if errors.Is(err, cbclient.ErrNotFound) {
				d.SetId("")
				return diags
			}

			return diag.FromErr(err)
		}

		if res.Status == "HISTORICAL" {
			d.SetId("")
			return diags
		}

		var servers []map[string]interface{}
		for _, s := range res.Links.Servers {
			svr, svrerr := apiClient.GetServer(s.Href)
			if svrerr != nil {
				return diag.Errorf("Error getting Servers for Resource: %s", svrerr.Error())
			}

			server, _ := parseServer(svr)
			servers = append(servers, server)

			if len(res.Links.Servers) == 1 {
				if attributes, ok := server["attributes"].(map[string]interface{}); ok {
					for k, v := range attributes {
						allAttributes[k] = v
					}
				}

				if techSpecificAttributes, ok := server["tech_specific_attributes"].(map[string]interface{}); ok {
					for k, v := range techSpecificAttributes {
						allAttributes[k] = v
					}
				}
				allAttributes["cpu_cnt"] = convertValueToString(svr.CPUCount)

				mem_size := convertValueToString(svr.MemorySizeGB)
				mem_size = strings.TrimRight(mem_size, "0")
				mem_size = strings.TrimRight(mem_size, ".")
				allAttributes["mem_size"] = mem_size
			}
		}

		d.Set("servers", servers)

		resAttributes, _ := parseAttributes(res.Attributes)
		for k, v := range resAttributes {
			allAttributes[k] = v
		}

		d.Set("attributes", resAttributes)
	} else {
		serverIds := strings.Split(d.Id(), "_")

		servers := make([]map[string]interface{}, 0)
		for _, serverId := range serverIds {
			svr, svrerr := apiClient.GetServerById(serverId)
			if svrerr != nil {
				return diag.Errorf("Error getting Server: %s", svrerr.Error())
			}

			server, _ := parseServer(svr)
			servers = append(servers, server)

			if len(serverIds) == 1 {
				if attributes, ok := server["attributes"].(map[string]interface{}); ok {
					for k, v := range attributes {
						allAttributes[k] = v
					}
				}

				if techSpecificAttributes, ok := server["tech_specific_attributes"].(map[string]interface{}); ok {
					for k, v := range techSpecificAttributes {
						allAttributes[k] = v
					}
				}

				allAttributes["cpu_cnt"] = convertValueToString(svr.CPUCount)

				mem_size := convertValueToString(svr.MemorySizeGB)
				mem_size = strings.TrimRight(mem_size, "0")
				mem_size = strings.TrimRight(mem_size, ".")
				allAttributes["mem_size"] = mem_size
			}
		}

		if servers != nil {
			d.Set("servers", servers)
		}
	}

	if parameters, ok := getParameters(d); ok {
		updatedParameters := setParameters(allAttributes, parameters)
		d.Set("parameters", updatedParameters)
	}

	if depItems, ok := getDeploymentItems(d); ok {
		updatedDepItems := setDeploymentItems(allAttributes, depItems)
		d.Set("deployment_item", updatedDepItems)
	}

	return diags
}

func resourceBPInstanceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	defer withPanicRecovery(&diags, "Update")

	instanceType := d.Get("instance_type").(string)
	if instanceType != "Resource" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "The CloudBolt provider does not support Terraform config updates for Servers.",
		})
		return diags
	}

	requestTimeout := d.Get("request_timeout").(int)
	if d.HasChange("parameters") || d.HasChange("deployment_item") {
		apiClient := m.(*cbclient.CloudBoltClient)
		actionPath, geterr := getResourceActionPath(apiClient, d.Id(), "Terraform Provider Update", true)
		if geterr != nil {
			return diag.FromErr(geterr)
		}

		if actionPath == "" {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "CloudBolt blueprint does not have a management action named \"Terraform Provider Update\", this action is required to apply terraform configuration changes.",
			})
			return diags
		}

		tfConfigParams := make(map[string]interface{}, 0)
		bpItemList := d.Get("deployment_item").(*schema.Set).List()
		bpParams := normalizeParameters(d.Get("parameters").(map[string]interface{}))

		if bpParams != nil {
			tfConfigParams["parameters"] = bpParams
		}

		for _, v := range bpItemList {
			m := v.(map[string]interface{})
			itemParams := normalizeParameters(m["parameters"].(map[string]interface{}))
			tfConfigParams[m["name"].(string)] = itemParams
		}

		parametersJSON, jsonerr := json.Marshal(tfConfigParams)
		if jsonerr != nil {
			fmt.Println(jsonerr)
		}

		fmt.Println(string(parametersJSON))

		parameters := map[string]interface{}{
			"tf_config_parameters": string(parametersJSON),
		}

		runActionResult, upderr := apiClient.SubmitAction(actionPath, d.Id(), parameters)
		if upderr != nil {
			return diag.FromErr(upderr)
		}

		if runActionResult.Results.Status != "" {
			if runActionResult.Results.Status != "SUCCESS" {
				var b strings.Builder

				// First line (single line, status included)
				fmt.Fprintf(
					&b,
					"Action failed (status=%s).\n\n",
					runActionResult.Results.Status,
				)

				// Errors (multiline)
				if strings.TrimSpace(runActionResult.Results.ErrorMessage) != "" {
					b.WriteString("Errors:\n")
					for _, line := range strings.Split(
						strings.TrimRight(runActionResult.Results.ErrorMessage, "\n"),
						"\n",
					) {
						fmt.Fprintf(&b, "  • %s\n", line)
					}
					b.WriteString("\n")
				}

				// Output (multiline)
				if strings.TrimSpace(runActionResult.Results.OutputMessage) != "" {
					b.WriteString("Output:\n")
					for _, line := range strings.Split(
						strings.TrimRight(runActionResult.Results.OutputMessage, "\n"),
						"\n",
					) {
						fmt.Fprintf(&b, "  • %s\n", line)
					}
				}

				return diag.Errorf(b.String())
			}
		} else {
			stateChangeConf := resource.StateChangeConf{
				Delay:   10 * time.Second,
				Timeout: time.Duration(requestTimeout) * time.Minute,
				Pending: []string{"INIT", "QUEUED", "PENDING", "RUNNING", "TO_CANCEL"},
				Target:  []string{"SUCCESS"},
			}

			var runProcessType string
			if runActionResult.Results.Job.Links.Self.Href != "" {
				stateChangeConf.Pending = []string{"INIT", "QUEUED", "PENDING", "RUNNING", "TO_CANCEL"}
				stateChangeConf.Refresh = JobStateRefreshFunc(apiClient, runActionResult.Results.Job.Links.Self.Href)
				runProcessType = "job"
			} else if runActionResult.Results.Order.Links.Self.Href != "" {
				stateChangeConf.Pending = []string{"ACTIVE"}
				stateChangeConf.Refresh = OrderStateRefreshFunc(apiClient, runActionResult.Results.Order.ID)
				runProcessType = "order"
			}

			_, err := stateChangeConf.WaitForState()
			if err != nil && runActionResult.Results.Job.Links.Self.Href != "" {
				return diag.Errorf("Error waiting for Job (%s) to complete: %s", runActionResult.Results.Job.Links.Self.Href, err)
			}

			if err != nil && runActionResult.Results.Order.Links.Self.Href != "" {
				return diag.Errorf("Error waiting for Order (%s) to complete: %s", runActionResult.Results.Order.Links.Self.Href, err)
			}

			if err != nil {
				return diag.Errorf(
					"Timed out after %d minutes waiting for %s to complete. Error: %s",
					requestTimeout,
					runProcessType,
					err,
				)
			}
		}
	}

	// Populate Terraform state by reading the resource
	readDiags := resourceBPInstanceRead(ctx, d, m)
	diags = append(diags, readDiags...)

	return diags
}

func resourceBPInstanceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	defer withPanicRecovery(&diags, "Delete")

	apiClient := m.(*cbclient.CloudBoltClient)
	instanceType := d.Get("instance_type").(string)

	requestTimeout := d.Get("request_timeout").(int)
	if instanceType == "Resource" {
		delActionPath, err := getResourceActionPath(apiClient, d.Id(), "Delete", false)
		if err != nil {
			return diag.FromErr(err)
		}

		if delActionPath == "" {
			return diag.Errorf("Error deleting resource (%s).", d.Id())
		}

		runActionResult, delerr := apiClient.SubmitAction(delActionPath, d.Id(), nil)
		if delerr != nil {
			return diag.FromErr(delerr)
		}

		if runActionResult.Results.Status != "" {
			if runActionResult.Results.Status != "SUCCESS" {
				var message string
				if runActionResult.Results.ErrorMessage != "" {
					message = runActionResult.Results.ErrorMessage
				} else {
					message = runActionResult.Results.OutputMessage
				}

				return diag.Errorf("Action Failed Status: %s Error: %s", runActionResult.Results.Status, message)
			}
		} else {
			stateChangeConf := resource.StateChangeConf{
				Delay:   10 * time.Second,
				Timeout: time.Duration(requestTimeout) * time.Minute,
				Target:  []string{"SUCCESS"},
			}
			if runActionResult.Results.Job.Links.Self.Href != "" {
				stateChangeConf.Pending = []string{"INIT", "QUEUED", "PENDING", "RUNNING", "TO_CANCEL"}
				stateChangeConf.Refresh = JobStateRefreshFunc(apiClient, runActionResult.Results.Job.Links.Self.Href)
			} else if runActionResult.Results.Order.Links.Self.Href != "" {
				stateChangeConf.Pending = []string{"ACTIVE"}
				stateChangeConf.Refresh = OrderStateRefreshFunc(apiClient, runActionResult.Results.Order.ID)
			}

			_, err = stateChangeConf.WaitForState()
			if err != nil && runActionResult.Results.Job.Links.Self.Href != "" {
				return diag.Errorf("Error waiting for Job (%s) to complete: %s", runActionResult.Results.Job.Links.Self.Href, err)
			}

			if err != nil && runActionResult.Results.Order.Links.Self.Href != "" {
				return diag.Errorf("Error waiting for Order (%s) to complete: %s", runActionResult.Results.Order.Links.Self.Href, err)
			}
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
			stateChangeConf.Timeout = time.Duration(requestTimeout) * time.Minute
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

	return diags
}

func withPanicRecovery(
	diags *diag.Diagnostics,
	operation string,
) {
	if r := recover(); r != nil {
		log.Printf(
			"[ERROR] [provider.cloudbolt] panic in %s: %v\n%s",
			operation,
			r,
			debug.Stack(),
		)

		*diags = append(*diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "CloudBolt provider crashed",
			Detail:   fmt.Sprintf("panic during %s: %v", operation, r),
		})
	}
}

func OrderStateRefreshFunc(apiClient *cbclient.CloudBoltClient, orderId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		order, err := apiClient.GetOrder(orderId)

		if err != nil {
			return nil, "ERROR", err
		}

		if order.Status == "FAILURE" {
			// return nil, order.Status, fmt.Errorf("Order %s failed to reach target state.", orderId)
			status, err := apiClient.GetOrderStatus(orderId)
			if err != nil {
				return nil, "FAILURE", fmt.Errorf(
					"Order %s failed, but status details could not be retrieved: %w",
					orderId,
					err,
				)
			}

			var b strings.Builder

			fmt.Fprintf(&b, "Order %s failed.\n\n", orderId)

			if len(status.OutputMessages) > 0 {
				b.WriteString("Outputs:\n")
				for _, msg := range status.OutputMessages {
					fmt.Fprintf(&b, "  • %s\n", msg)
				}
				b.WriteString("\n")
			}

			if len(status.ErrorMessages) > 0 {
				b.WriteString("Errors:\n")
				for _, msg := range status.ErrorMessages {
					fmt.Fprintf(&b, "  • %s\n", msg)
				}
			}

			return nil, order.Status, fmt.Errorf(b.String())
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
			var b strings.Builder

			fmt.Fprintf(&b, "Job %s failed to reach target state.\n\n", job.ID)

			if strings.TrimSpace(job.Errors) != "" {
				b.WriteString("Errors:\n")
				for _, line := range strings.Split(strings.TrimRight(job.Errors, "\n"), "\n") {
					fmt.Fprintf(&b, "  • %s\n", line)
				}
				b.WriteString("\n")
			}

			if strings.TrimSpace(job.Output) != "" {
				b.WriteString("Outputs:\n")
				for _, line := range strings.Split(strings.TrimRight(job.Output, "\n"), "\n") {
					fmt.Fprintf(&b, "  • %s\n", line)
				}
			}

			return nil, job.Status, fmt.Errorf(b.String())
		}

		return job, job.Status, nil
	}
}

func getResourceActionPath(apiClient *cbclient.CloudBoltClient, resourcePath string, resourceActionName string, prefixFilter bool) (string, error) {
	var actionPath string
	res, err := apiClient.GetResource(resourcePath)
	if err != nil {
		return actionPath, err
	}

	for _, v := range res.Links.Actions {
		if prefixFilter {
			if strings.HasPrefix(v.Title, resourceActionName) {
				actionPath = v.Href
				break
			}
		} else {
			if v.Title == resourceActionName {
				actionPath = v.Href
				break
			}
		}
	}

	return actionPath, nil
}

func normalizeParameters(params map[string]interface{}) map[string]interface{} {
	normalizedParams := make(map[string]interface{}, 0)

	for k, v := range params {
		value, ok := v.(string)
		if !ok {
			normalizedParams[k] = v
			continue
		}

		if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
			parameterValues := strings.Split(value[1:len(value)-1], "|")
			normalizedParams[k] = parameterValues
		} else {
			normalizedParams[k] = v
		}
	}

	return normalizedParams
}
