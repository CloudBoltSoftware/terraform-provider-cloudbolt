package cloudbolt

import (
	"fmt"
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
			"blueprint_item": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
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
		},
	}
}

func resourceBPInstanceCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(Config).APIClient

	bpItems := make([]map[string]interface{}, 0)
	bpItemList := d.Get("blueprint_item").(*schema.Set).List()

	for _, v := range bpItemList {
		m := v.(map[string]interface{})

		bpItem := map[string]interface{}{
			"bp-item-name":    m["name"].(string),
			"bp-item-paramas": m["parameters"].(map[string]interface{}),
		}

		bpItems = append(bpItems, bpItem)
	}

	order, err := apiClient.DeployBlueprint(d.Get("group").(string), d.Get("blueprint").(string), bpItems)
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
		return fmt.Errorf("Error waiting for Order (%s) to complete: %s", order.ID, err)
	}

	// Retrieve the updated order to obtain Resource ID
	order, orderr := apiClient.GetOrder(order.ID)
	if orderr != nil {
		return orderr
	}

	var resourceId string
	for _, j := range order.Links.Jobs {
		job, joberr := apiClient.GetJob(j.Href)
		if joberr != nil {
			return joberr
		}

		if job.Type == "Deploy Blueprint" {
			resourceId = job.Links.Resource.Href
			break
		}
	}

	if resourceId == "" {
		return fmt.Errorf("Error Order (%s) does not have a Resource", order.ID)
	}

	d.SetId(resourceId)

	return resourceBPInstanceRead(d, m)
}

func resourceBPInstanceRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(Config).APIClient

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

	return nil
}

func resourceBPInstanceUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceBPInstanceRead(d, m)
}

func resourceBPInstanceDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(Config).APIClient

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

	return nil
}

func OrderStateRefreshFunc(config Config, orderId string) resource.StateRefreshFunc {
	apiClient := config.APIClient

	return func() (interface{}, string, error) {
		order, err := apiClient.GetOrder(orderId)
		if err != nil {
			return nil, "", err
		}

		if order.Status == "FAILURE" {
			return nil, order.Status, fmt.Errorf("Order %s failed to reach target state.", orderId)
		}

		return order, order.Status, nil
	}
}

func JobStateRefreshFunc(config Config, jobPath string) resource.StateRefreshFunc {
	apiClient := config.APIClient

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
