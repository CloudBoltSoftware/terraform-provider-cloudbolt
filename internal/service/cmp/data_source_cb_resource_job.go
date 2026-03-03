package cmp

import (
	"context"
	"strings"

	"github.com/cloudboltsoftware/cloudbolt-go-sdk/cbclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceCloudBoltResourceJobs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCloudBoltResourceJobsRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The global id of a CloudBolt Resource, required if \"url_path\" is not provided",
			},
			"url_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The relative API URL path for the CloudBolt Resource, required if \"id\" is not provided",
			},
			"job_info": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Job information for the CloudBolt Resource",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"title": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"start_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"end_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"output": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"error": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"progress_messages": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourceCloudBoltResourceJobsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*cbclient.CloudBoltClient)
	id := d.Get("id").(string)
	urlPath := d.Get("url_path").(string)

	if id == "" && urlPath == "" {
		return diag.Errorf("Either id or url_path is required")
	}

	var resourceJobInfo *cbclient.CloudBoltResourceJobInfo
	var err error
	if urlPath != "" {
		jobInfoPath := strings.TrimRight(urlPath, "/") + "/jobsInfo/"
		resourceJobInfo, err = apiClient.GetResourceJobInfo(jobInfoPath)
	} else {
		resourceJobInfo, err = apiClient.GetResourceJobInfoById(id)
	}

	if err != nil {
		return diag.FromErr(err)
	}

	if id != "" {
		d.SetId(id)
	} else {
		d.SetId(urlPath)
	}

	jobInfoList := make([]map[string]interface{}, len(*resourceJobInfo))
	for i, job := range *resourceJobInfo {
		jobInfoList[i] = map[string]interface{}{
			"title":             job.Title,
			"start_date":        job.StartDate,
			"end_date":          job.EndDate,
			"status":            job.Status,
			"output":            job.Output,
			"error":             job.Error,
			"progress_messages": job.ProgressMessages,
		}
	}
	d.Set("job_info", jobInfoList)

	return nil
}
