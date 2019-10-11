package cloudbolt

import (
	"fmt"
	// "log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceBPInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceBPInstanceCreate,
		Read:   resourceBPInstanceRead,
		Update: resourceBPInstanceUpdate,
		Delete: resourceBPInstanceDelete,

		Schema: map[string]*schema.Schema{
			"group": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"blueprint": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"resource_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"blueprint_item": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"environment": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"osbuild": &schema.Schema{
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
						"ip": {
							Type:     schema.TypeString,
							Computed: true,
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
		},
	}
}

func resourceBPInstanceCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(Config).APIClient

	// log.Printf("[!!] apiClient in resourceBPInstanceCreate: %+v", apiClient)

	bpItems := make([]map[string]interface{}, 0)
	bpItemList := d.Get("blueprint_item").(*schema.Set).List()

	for _, v := range bpItemList {
		m := v.(map[string]interface{})

		bpItem := map[string]interface{}{
			"bp-item-name":    m["name"].(string),
			"bp-item-paramas": m["parameters"].(map[string]interface{}),
		}

		env, ok := m["environment"]
		if ok {
			bpItem["environment"] = env.(string)
		}

		osb, ok := m["osbuild"]
		if ok {
			bpItem["os-build"] = osb.(string)
		}

		bpItems = append(bpItems, bpItem)
	}

	order, err := apiClient.DeployBlueprint(d.Get("group").(string), d.Get("blueprint").(string), d.Get("resource_name").(string), bpItems)
	if err != nil {
		return err
	}

	stateChangeConf := resource.StateChangeConf{
		Delay:   10 * time.Second,
		Timeout: 5 * time.Minute,
		Pending: []string{"ACTIVE"},
		Target:  []string{"SUCCESS"},
		Refresh: OrderStateRefreshFunc(m.(Config), order.ID),
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
	var serverId string
	for _, j := range order.Links.Jobs {
		job, joberr := apiClient.GetJob(j.Href)
		if joberr != nil {
			return joberr
		}

		if job.Type == "Deploy Blueprint" {
			// log.Println("[!!] Deploying Blueprint")
			if len(job.Links.Resource.Href) > 0 {
				resourceId = job.Links.Resource.Href
				d.Set("instance_type", "Resource")
			} else if len(job.Links.Servers) > 0 {
				serverId = job.Links.Servers[0].Href
				d.Set("instance_type", "Server")
			}
			break
		}
	}

	if resourceId == "" && serverId == "" {
		return fmt.Errorf("Error Order (%s) does not have a Resource or Server", order.ID)
	}

	if resourceId != "" {
		d.SetId(resourceId)
	} else {
		d.SetId(serverId)
	}

	return resourceBPInstanceRead(d, m)
}

func resourceBPInstanceRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(Config).APIClient
	instanceType := d.Get("instance_type").(string)

	// log.Printf("[!!] apiClient in resourceBPInstanceRead: %+v", apiClient)

	if instanceType == "Resource" {
		res, err := apiClient.GetResource(d.Id())
		if err != nil {
			return err
		}

		servers := make([]map[string]interface{}, 0)
		for _, s := range res.Links.Servers {
			svr, svrerr := apiClient.GetServer(s.Href)
			if svrerr != nil {
				return fmt.Errorf("Error getting Servers for Resource: %s", svrerr)
			}

			servers = append(servers, map[string]interface{}{
				"hostname": svr.Hostname,
				"ip":       svr.IP,
			})

			d.Set("server_hostname", svr.Hostname)
			d.Set("server_ip", svr.IP)
		}

		if servers != nil {
			d.Set("servers", servers)
		}
	} else {
		svr, svrerr := apiClient.GetServer(d.Id())
		if svrerr != nil {
			return fmt.Errorf("Error getting Server: %s", svrerr)
		}

		d.Set("server_hostname", svr.Hostname)
		d.Set("server_ip", svr.IP)
	}

	return nil
}

func resourceBPInstanceUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceBPInstanceRead(d, m)
}

func resourceBPInstanceDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(Config).APIClient
	instanceType := d.Get("instance_type").(string)

	// log.Printf("[!!] apiClient in resourceBPInstanceDelete: %+v", apiClient)

	if instanceType == "Resource" {
		res, err := apiClient.GetResource(d.Id())
		if err != nil {
			return err
		}

		var delResPath string
		for _, v := range res.Links.Actions {
			if v.Delete.Href != "" {
				delResPath = v.Delete.Href
				break
			}
		}

		if delResPath == "" {
			return fmt.Errorf("Error deleting resource (%s).", d.Id())
		}

		job, delerr := apiClient.SubmitAction(delResPath)
		if delerr != nil {
			return delerr
		}

		stateChangeConf := resource.StateChangeConf{
			Delay:   10 * time.Second,
			Timeout: 5 * time.Minute,
			Pending: []string{"INIT", "QUEUED", "PENDING", "RUNNING", "TO_CANCEL"},
			Target:  []string{"SUCCESS"},
			Refresh: JobStateRefreshFunc(m.(Config), job.RunActionJob.Self.Href),
		}

		_, err = stateChangeConf.WaitForState()
		if err != nil {
			return fmt.Errorf("Error waiting for Job (%s) to complete: %s", job.RunActionJob.Self.Href, err)
		}
	} else {
		svr, err := apiClient.GetServer(d.Id())
		if err != nil {
			return err
		}

		servers := []string{
			svr.Links.Self.Href,
		}

		order, err := apiClient.DecomOrder(svr.Links.Group.Href, svr.Links.Environment.Href, servers)
		if err != nil {
			return err
		}

		stateChangeConf := resource.StateChangeConf{
			Delay:   10 * time.Second,
			Timeout: 5 * time.Minute,
			Pending: []string{"ACTIVE"},
			Target:  []string{"SUCCESS"},
			Refresh: OrderStateRefreshFunc(m.(Config), order.ID),
		}

		_, err = stateChangeConf.WaitForState()
		if err != nil {
			return fmt.Errorf("Error waiting for Decom Order (%s) to complete: %s", order.ID, err)
		}
	}

	return nil
}

func OrderStateRefreshFunc(config Config, orderId string) resource.StateRefreshFunc {
	apiClient := config.APIClient

	// log.Printf("[!!] apiClient in OrderStateRefreshFunc: %+v", apiClient)

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

func JobStateRefreshFunc(config Config, jobPath string) resource.StateRefreshFunc {
	apiClient := config.APIClient

	// log.Printf("[!!] apiClient in JobStateRefreshFunc: %+v", apiClient)

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
