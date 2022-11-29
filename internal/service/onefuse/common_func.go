package onefuse

import (
	"fmt"
	"time"

	"github.com/cloudboltsoftware/cloudbolt-go-sdk/cbclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func JobStatusStateRefreshFunc(apiClient *cbclient.CloudBoltClient, jobStatusPath string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		jobStatus, err := apiClient.GetJobStatus(jobStatusPath)
		if err != nil {
			return nil, "", err
		}

		if jobStatus.JobState == "Canceled" || jobStatus.JobState == "Failed" {
			return nil, jobStatus.JobState, fmt.Errorf("Job %s failed to reach target state.", jobStatusPath)
		}

		return jobStatus, jobStatus.JobState, nil
	}
}

func GetJobStautusStateChangeConf(apiClient *cbclient.CloudBoltClient, requestTimeout int, jobStatusPath string) resource.StateChangeConf {
	return resource.StateChangeConf{
		Delay:   10 * time.Second,
		Timeout: time.Duration(requestTimeout) * time.Minute,
		Pending: []string{
			"Initialized",
			"In_Progress",
		},
		Target:  []string{"Successful"},
		Refresh: JobStatusStateRefreshFunc(apiClient, jobStatusPath),
	}
}
