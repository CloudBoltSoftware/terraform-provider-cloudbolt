package cbclient

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type CloudBoltObject struct {
	Links struct {
		Self struct {
			Href  string `json:"href"`
			Title string `json:"title"`
		} `json:"self"`
	} `json:"_links"`
	Name string `json:"name"`
	ID   string `json:"id"`
}

type CloudBoltClient struct {
	BaseURL    url.URL
	HTTPClient *http.Client
	Token      string
}

type CloudBoltResult struct {
	Links struct {
		Self struct {
			Href  string `json:"href"`
			Title string `json:"title"`
		} `json:"self"`
	} `json:"_links"`
	Total    int               `json:"total"`
	Count    int               `json:"count"`
	Embedded []CloudBoltObject `json:"_embedded"`
}

type CloudBoltActionResult struct {
	RunActionJob struct {
		Self struct {
			Href  string `json:"href"`
			Title string `json:"title"`
		} `json:"self"`
	} `json:"run-action-job"`
}

type CloudBoltHALItem struct {
	Href  string `json:"href"`
	Title string `json:"title"`
}

type CloudBoltOrder struct {
	Links struct {
		Self struct {
			Href  string `json:"href"`
			Title string `json:"title"`
		} `json:"self"`
		Group struct {
			Href  string `json:"href"`
			Title string `json:"title"`
		} `json:"group"`
		Owner struct {
			Href  string `json:"href"`
			Title string `json:"title"`
		} `json:"owner"`
		ApprovedBy struct {
			Href  string `json:"href"`
			Title string `json:"title"`
		} `json:"approved-by"`
		Actions struct {
			Href  string `json:"href"`
			Title string `json:"title"`
		} `json:"actions"`
		Jobs []CloudBoltHALItem `json:"jobs"`
	} `json:"_links"`
	Name        string `json:"name"`
	ID          string `json:"id"`
	Status      string `json:"status"`
	Rate        string `json:"rate"`
	CreateDate  string `json:"create-date"`
	ApproveDate string `json:"approve-date"`
	Items       struct {
		DeployItems []struct {
			Blueprint               string `json:"blueprint"`
			BlueprintItemsArguments struct {
				BuildItemBuildServer struct {
					Attributes struct {
						Hostname string `json:"hostname"`
						Quantity int    `json:"quantity"`
					} `json:"attributes"`
					OsBuild     string                 `json:"os-build,omitempty"`
					Environment string                 `json:"environment,omitempty"`
					Parameters  map[string]interface{} `json:"parameters"`
				} `json:"build-item-Server"`
			} `json:"blueprint-items-arguments"`
			ResourceName       string `json:"resource-name"`
			ResourceParameters struct {
			} `json:"resource-parameters"`
		} `json:"deploy-items"`
	} `json:"items"`
}

type CloudBoltJob struct {
	Links struct {
		Self struct {
			Href  string `json:"href"`
			Title string `json:"title"`
		} `json:"self"`
		Owner struct {
			Href  string `json:"href"`
			Title string `json:"title"`
		} `json:"owner"`
		Parent struct {
			Href  string `json:"href"`
			Title string `json:"title"`
		} `json:"parent"`
		Subjobs      []interface{} `json:"subjobs"`
		Prerequisite struct {
		} `json:"prerequisite"`
		DependentJobs []interface{} `json:"dependent-jobs"`
		Order         struct {
			Href  string `json:"href"`
			Title string `json:"title"`
		} `json:"order"`
		Resource struct {
			Href  string `json:"href"`
			Title string `json:"title"`
		} `json:"resource"`
		Servers []struct {
			Href  string `json:"href"`
			Title string `json:"title"`
		} `json:"servers"`
		LogUrls struct {
			RawLog string `json:"raw-log"`
			ZipLog string `json:"zip-log"`
		} `json:"log_urls"`
	} `json:"_links"`
	Status   string `json:"status"`
	Type     string `json:"type"`
	Progress struct {
		TotalTasks int      `json:"total-tasks"`
		Completed  int      `json:"completed"`
		Messages   []string `json:"messages"`
	} `json:"progress"`
	StartDate string `json:"start-date"`
	EndDate   string `json:"end-date"`
	Output    string `json:"output"`
}

type CloudBoltGroup struct {
	Links struct {
		Self struct {
			Href  string `json:"href"`
			Title string `json:"title"`
		} `json:"self"`
		Parent struct {
			Href  string `json:"href"`
			Title string `json:"title"`
		} `json:"parent"`
		Subgroups []struct {
			Href  string `json:"href"`
			Title string `json:"title"`
		} `json:"subgroups"`
		Environments          []interface{} `json:"environments"`
		OrderableEnvironments struct {
			Href  string `json:"href"`
			Title string `json:"title"`
		} `json:"orderable-environments"`
	} `json:"_links"`
	Name         string `json:"name"`
	ID           string `json:"id"`
	Type         string `json:"type"`
	Rate         string `json:"rate"`
	AutoApproval bool   `json:"auto-approval"`
}

type CloudBoltResource struct {
	Links struct {
		Self struct {
			Href  string `json:"href"`
			Title string `json:"title"`
		} `json:"self"`
		Blueprint struct {
			Href  string `json:"href"`
			Title string `json:"title"`
		} `json:"blueprint"`
		Owner struct {
			Href  string `json:"href"`
			Title string `json:"title"`
		} `json:"owner"`
		Group struct {
			Href  string `json:"href"`
			Title string `json:"title"`
		} `json:"group"`
		ResourceType struct {
			Href  string `json:"href"`
			Title string `json:"title"`
		} `json:"resource-type"`
		Servers []struct {
			Href  string `json:"href"`
			Title string `json:"title"`
			Tier  string `json:"tier"`
		} `json:"servers"`
		Actions []struct {
			Delete struct {
				Href  string `json:"href"`
				Title string `json:"title"`
			} `json:"Delete,omitempty"`
			Scale struct {
				Href  string `json:"href"`
				Title string `json:"title"`
			} `json:"Scale,omitempty"`
		} `json:"actions"`
		Jobs struct {
			Href  string `json:"href"`
			Title string `json:"title"`
		} `json:"jobs"`
		History struct {
			Href  string `json:"href"`
			Title string `json:"title"`
		} `json:"history"`
	} `json:"_links"`
	Name        string `json:"name"`
	ID          string `json:"id"`
	Status      string `json:"status"`
	InstallDate string `json:"install-date"`
}

type CloudBoltServer struct {
	Links struct {
		Self struct {
			Href  string `json:"href"`
			Title string `json:"title"`
		} `json:"self"`
		Owner struct {
			Href  string `json:"href"`
			Title string `json:"title"`
		} `json:"owner"`
		Group struct {
			Href  string `json:"href"`
			Title string `json:"title"`
		} `json:"group"`
		Environment struct {
			Href  string `json:"href"`
			Title string `json:"title"`
		} `json:"environment"`
		ResourceHandler struct {
			Href  string `json:"href"`
			Title string `json:"title"`
		} `json:"resource-handler"`
		Actions []struct {
			PowerOn struct {
				Href  string `json:"href"`
				Title string `json:"title"`
			} `json:"power_on,omitempty"`
			PowerOff struct {
				Href  string `json:"href"`
				Title string `json:"title"`
			} `json:"power_off,omitempty"`
			Reboot struct {
				Href  string `json:"href"`
				Title string `json:"title"`
			} `json:"reboot,omitempty"`
			RefreshInfo struct {
				Href  string `json:"href"`
				Title string `json:"title"`
			} `json:"refresh_info,omitempty"`
			Snapshot struct {
				Title string `json:"title"`
				Href  string `json:"href"`
			} `json:"snapshot,omitempty"`
			AdHocScript struct {
				Href  string `json:"href"`
				Title string `json:"title"`
			} `json:"Ad Hoc Script,omitempty"`
		} `json:"actions"`
		ProvisionJob struct {
			Href  string `json:"href"`
			Title string `json:"title"`
		} `json:"provision-job"`
		OsBuild struct {
			Href  string `json:"href"`
			Title string `json:"title"`
		} `json:"os-build"`
		Jobs struct {
			Href  string `json:"href"`
			Title string `json:"title"`
		} `json:"jobs"`
		History struct {
			Href  string `json:"href"`
			Title string `json:"title"`
		} `json:"history"`
	} `json:"_links"`
	Hostname             string        `json:"hostname"`
	PowerStatus          string        `json:"power-status"`
	Status               string        `json:"status"`
	IP                   string        `json:"ip"`
	Mac                  string        `json:"mac"`
	DateAddedToCloudbolt string        `json:"date-added-to-cloudbolt"`
	CPUCnt               int           `json:"cpu-cnt"`
	MemSize              string        `json:"mem-size"`
	DiskSize             string        `json:"disk-size"`
	OsFamily             string        `json:"os-family"`
	Notes                string        `json:"notes"`
	Labels               []interface{} `json:"labels"`
	Credentials          struct {
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"credentials"`
	Disks []struct {
		UUID             string `json:"uuid"`
		DiskSize         int    `json:"disk-size"`
		Name             string `json:"name"`
		Datastore        string `json:"datastore"`
		ProvisioningType string `json:"provisioning-type"`
	} `json:"disks"`
	Networks []struct {
		Name          string      `json:"name"`
		Network       string      `json:"network"`
		Mac           string      `json:"mac"`
		IP            interface{} `json:"ip"`
		PrivateIP     string      `json:"private-ip"`
		AdditionalIps string      `json:"additional-ips"`
	} `json:"networks"`
	Parameters struct {
	} `json:"parameters"`
	TechSpecificDetails struct {
		VmwareLinkedClone bool   `json:"vmware-linked-clone"`
		VmwareCluster     string `json:"vmware-cluster"`
	} `json:"tech-specific-details"`
}

func New(protocol string, host string, port string, username string, password string) (CloudBoltClient, error) {
	var cbClient CloudBoltClient
	cbClient.HTTPClient = &http.Client{
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // TODO make this configurable
			},
		},
	}

	cbClient.BaseURL = url.URL{
		Scheme: protocol,
		Host:   fmt.Sprintf("%s:%s", host, port),
	}

	reqJson, err := json.Marshal(struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		Username: username,
		Password: password,
	})

	if err != nil {
		log.Fatalln(err)
		return cbClient, err
	}

	apiurl := cbClient.BaseURL
	apiurl.Path = "/api/v2/api-token-auth/"

	// log.Printf("[!!] apiurl in New: %+v (%+v)", apiurl.String(), apiurl)

	resp, err := cbClient.HTTPClient.Post(apiurl.String(), "application/json", bytes.NewBuffer(reqJson))
	if err != nil {
		log.Fatalf("Failed to create the API client. %s", err)
	}

	userAuthData := struct {
		Token string `json:"token"`
	}{}

	json.NewDecoder(resp.Body).Decode(&userAuthData)
	cbClient.Token = userAuthData.Token

	// log.Printf("[!!] cbClient: %+v", cbClient)

	return cbClient, nil
}

func (cbClient CloudBoltClient) GetCloudBoltObject(objPath string, objName string) (CloudBoltObject, error) {
	apiurl := cbClient.BaseURL
	apiurl.Path = fmt.Sprintf("/api/v2/%s/", objPath)
	apiurl.RawQuery = fmt.Sprintf("filter=name:%s", objName)

	// log.Printf("[!!] apiurl in GetCloudBoltObject: %+v (%+v)", apiurl.String(), apiurl)

	req, err := http.NewRequest("GET", apiurl.String(), nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", cbClient.Token))
	req.Header.Add("Content-Type", "application/json")

	resp, err := cbClient.HTTPClient.Do(req)
	if err != nil {
		log.Fatalln(err)

		return CloudBoltObject{}, err
	}
	// log.Printf("[!!] HTTP response: %+v", resp)

	var res CloudBoltResult
	json.NewDecoder(resp.Body).Decode(&res)

	// log.Printf("[!!] CloudBoltResult response %+v", res) // HERE IS WHERE THE PANIC IS!!!
	if len(res.Embedded) == 0 {
		return CloudBoltObject{}, err
	}
	return res.Embedded[0], nil
}

func (cbClient CloudBoltClient) verifyGroup(groupPath string, parentPath string) (bool, error) {
	var group CloudBoltGroup
	var parent string
	var nextParentPath string

	apiurl := cbClient.BaseURL
	apiurl.Path = groupPath

	// log.Printf("[!!] apiurl in verifyGroup: %+v (%+v)", apiurl.String(), apiurl)

	req, err := http.NewRequest("GET", apiurl.String(), nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", cbClient.Token))
	req.Header.Add("Content-Type", "application/json")

	// log.Printf("[!!] req: %+v", req)

	resp, err := cbClient.HTTPClient.Do(req)
	// log.Printf("[!!] resp: %+v", resp)
	if err != nil {
		// log.Printf("[!!] request err was not nil: %+v", err)
		log.Fatalln(err)

		return false, err
	}
	if resp.StatusCode >= 300 {
		// log.Printf("[!!] request returned a bad status: %+v", resp.Status)
		log.Fatalln(resp.Status)

		return false, errors.New(resp.Status)
	}

	json.NewDecoder(resp.Body).Decode(&group)

	// log.Printf("[!!] group : %+v", group)

	nextIndex := strings.LastIndex(parentPath, "/")

	// log.Printf("[!!] nextIndex : %+v", nextIndex)

	// log.Printf("[!!] parentPath: %+v", parentPath)
	// log.Printf("[!!] strings.LastIndex(parentPath, '/')+1: %+v", strings.LastIndex(parentPath, "/")+1)
	if nextIndex >= 0 {
		parent = parentPath[strings.LastIndex(parentPath, "/")+1:]
		nextParentPath = parentPath[:strings.LastIndex(parentPath, "/")]
		// log.Printf("[!!] parent: %+v, %+v", parent, nextParentPath)
	} else {
		parent = parentPath
		// log.Printf("[!!] parent: %+v", parent)
	}

	// log.Printf("[!!] group.Links.Parent.Title: %+v", group.Links.Parent.Title)
	if group.Links.Parent.Title != parent {
		return false, nil
	}

	// log.Printf("[!!] nextParentPath: %+v", nextParentPath)
	if nextParentPath != "" {
		return cbClient.verifyGroup(group.Links.Parent.Href, nextParentPath)
	}

	return true, nil
}

func (cbClient CloudBoltClient) GetGroup(groupPath string) (CloudBoltObject, error) {
	var res CloudBoltResult
	var group string
	var parentPath string
	var groupFound bool

	groupPath = strings.Trim(groupPath, "/")
	nextIndex := strings.LastIndex(groupPath, "/")

	// log.Printf("[!!] groupPath: %+v", groupPath)
	// log.Printf("[!!] strings.LastIndex(groupPath, '/')+1: %+v", strings.LastIndex(groupPath, "/")+1)
	if nextIndex >= 0 {
		group = groupPath[strings.LastIndex(groupPath, "/")+1:]
		parentPath = groupPath[:strings.LastIndex(groupPath, "/")]
	} else {
		group = groupPath
	}

	apiurl := cbClient.BaseURL
	apiurl.Path = "/api/v2/groups/"
	apiurl.RawQuery = fmt.Sprintf("filter=name:%s", url.QueryEscape(group))

	// log.Printf("[!!] apiurl in GetGroup: %+v (%+v)", apiurl.String(), apiurl)

	req, err := http.NewRequest("GET", apiurl.String(), nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", cbClient.Token))
	req.Header.Add("Content-Type", "application/json")

	resp, err := cbClient.HTTPClient.Do(req)
	if err != nil {
		log.Fatalln(err)

		return CloudBoltObject{}, err
	}

	json.NewDecoder(resp.Body).Decode(&res)

	for _, v := range res.Embedded {
		groupFound, err = cbClient.verifyGroup(v.Links.Self.Href, parentPath)

		if groupFound {
			return v, nil
		}
	}

	return CloudBoltObject{}, fmt.Errorf("Group (%s): Not Found", groupPath)
}

func (cbClient CloudBoltClient) DeployBlueprint(grpPath string, bpPath string, resourceName string, bpItems []map[string]interface{}) (CloudBoltOrder, error) {
	var order CloudBoltOrder

	deployItems := make([]map[string]interface{}, 0)

	for _, v := range bpItems {
		bpItem := map[string]interface{}{
			"blueprint": bpPath,
			"blueprint-items-arguments": map[string]interface{}{
				v["bp-item-name"].(string): map[string]interface{}{
					"attributes": map[string]interface{}{
						"quantity": 1,
					},
					"parameters": v["bp-item-paramas"].(map[string]interface{}),
				},
			},
			"resource-name": resourceName,
		}

		env, ok := v["environment"]
		if ok {
			bpItem["blueprint-items-arguments"].(map[string]interface{})[v["bp-item-name"].(string)].(map[string]interface{})["environment"] = env
		}

		osb, ok := v["os-build"]
		if ok {
			bpItem["blueprint-items-arguments"].(map[string]interface{})[v["bp-item-name"].(string)].(map[string]interface{})["os-build"] = osb
		}

		deployItems = append(deployItems, bpItem)
	}

	reqData := map[string]interface{}{
		"group": grpPath,
		"items": map[string]interface{}{
			"deploy-items": deployItems,
		},
		"submit-now": "true",
	}

	reqJSON, err := json.Marshal(reqData)
	if err != nil {
		log.Fatalln(err)
		return order, err
	}

	// log.Printf("[!!] JSON payload in POST request to Deploy Blueprint:\n%s", string(reqJSON))

	apiurl := cbClient.BaseURL
	apiurl.Path = "/api/v2/orders/"

	// log.Printf("[!!] apiurl in DeployBlueprint: %+v (%+v)", apiurl.String(), apiurl)

	reqBody := bytes.NewBuffer(reqJSON)

	req, err := http.NewRequest("POST", apiurl.String(), reqBody)
	if err != nil {
		log.Fatalln(err)
		return order, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", cbClient.Token))
	req.Header.Add("Content-Type", "application/json")
	// TODO: Make the API more responsive
	cbClient.HTTPClient.Timeout = 60 * time.Second

	resp, err := cbClient.HTTPClient.Do(req)
	if err != nil {
		log.Fatalln(err)
		return CloudBoltOrder{}, err
	}

	// Handle some common HTTP errors
	switch {
	case resp.StatusCode >= 500:
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		respBody := buf.String()
		return CloudBoltOrder{}, fmt.Errorf("recieved a server error: %s", respBody)
	case resp.StatusCode >= 400:
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		respBody := buf.String()
		return CloudBoltOrder{}, fmt.Errorf("recieved an HTTP client error: %s", respBody)
	default:
		json.NewDecoder(resp.Body).Decode(&order)

		return order, nil
	}
}

func (cbClient CloudBoltClient) GetOrder(orderId string) (CloudBoltOrder, error) {
	var order CloudBoltOrder

	apiurl := cbClient.BaseURL
	apiurl.Path = fmt.Sprintf("/api/v2/orders/%s", orderId)

	// log.Printf("[!!] apiurl in GetOrder: %+v (%+v)", apiurl.String(), apiurl)

	req, err := http.NewRequest("GET", apiurl.String(), nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", cbClient.Token))
	req.Header.Add("Content-Type", "application/json")

	resp, err := cbClient.HTTPClient.Do(req)
	if err != nil {
		log.Fatalln(err)

		return CloudBoltOrder{}, err
	}

	json.NewDecoder(resp.Body).Decode(&order)

	return order, nil
}

func (cbClient CloudBoltClient) GetJob(jobPath string) (CloudBoltJob, error) {
	var job CloudBoltJob

	apiurl := cbClient.BaseURL
	apiurl.Path = jobPath

	// log.Printf("[!!] GetJob: %+v (%+v)", apiurl.String(), apiurl)

	req, err := http.NewRequest("GET", apiurl.String(), nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", cbClient.Token))
	req.Header.Add("Content-Type", "application/json")

	resp, err := cbClient.HTTPClient.Do(req)
	if err != nil {
		log.Fatalln(err)

		return CloudBoltJob{}, err
	}

	json.NewDecoder(resp.Body).Decode(&job)

	return job, nil
}

func (cbClient CloudBoltClient) GetResource(resourcePath string) (CloudBoltResource, error) {
	var res CloudBoltResource

	apiurl := cbClient.BaseURL
	apiurl.Path = resourcePath

	// log.Printf("[!!] apiurl in GetResource: %+v (%+v)", apiurl.String(), apiurl)

	req, err := http.NewRequest("GET", apiurl.String(), nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", cbClient.Token))
	req.Header.Add("Content-Type", "application/json")

	resp, err := cbClient.HTTPClient.Do(req)
	if err != nil {
		log.Fatalln(err)

		return CloudBoltResource{}, err
	}

	json.NewDecoder(resp.Body).Decode(&res)

	return res, nil
}

func (cbClient CloudBoltClient) GetServer(serverPath string) (CloudBoltServer, error) {
	var svr CloudBoltServer

	apiurl := cbClient.BaseURL
	apiurl.Path = serverPath

	// log.Printf("[!!] apiurl in GetServer: %+v (%+v)", apiurl.String(), apiurl)

	req, err := http.NewRequest("GET", apiurl.String(), nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", cbClient.Token))
	req.Header.Add("Content-Type", "application/json")

	resp, err := cbClient.HTTPClient.Do(req)
	if err != nil {
		log.Fatalln(err)
		return CloudBoltServer{}, err
	}

	json.NewDecoder(resp.Body).Decode(&svr)

	return svr, nil
}

func (cbClient CloudBoltClient) SubmitAction(actionPath string) (CloudBoltActionResult, error) {
	var actionRes CloudBoltActionResult

	apiurl := cbClient.BaseURL
	apiurl.Path = actionPath

	// log.Printf("[!!] apiurl in SubmitAction: %+v (%+v)", apiurl.String(), apiurl)

	req, err := http.NewRequest("POST", apiurl.String(), nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", cbClient.Token))
	req.Header.Add("Content-Type", "application/json")

	resp, err := cbClient.HTTPClient.Do(req)
	if err != nil {
		log.Fatalln(err)
		return CloudBoltActionResult{}, err
	}

	json.NewDecoder(resp.Body).Decode(&actionRes)

	return actionRes, nil
}

func (cbClient CloudBoltClient) DecomOrder(grpPath string, envPath string, servers []string) (CloudBoltOrder, error) {
	var order CloudBoltOrder

	decomItems := make([]map[string]interface{}, 0)

	decomItem := make(map[string]interface{})
	decomItem["environment"] = envPath
	decomItem["servers"] = servers

	decomItems = append(decomItems, decomItem)

	reqData := map[string]interface{}{
		"group": grpPath,
		"items": map[string]interface{}{
			"decom-items": decomItems,
		},
		"submit-now": "true",
	}

	reqJson, err := json.Marshal(reqData)
	if err != nil {
		log.Fatalln(err)
		return CloudBoltOrder{}, err
	}

	apiurl := cbClient.BaseURL
	apiurl.Path = "/api/v2/orders/"

	// log.Printf("[!!] apiurl in DecomOrder: %+v (%+v)", apiurl.String(), apiurl)

	req, err := http.NewRequest("POST", apiurl.String(), bytes.NewBuffer(reqJson))
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", cbClient.Token))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := cbClient.HTTPClient.Do(req)
	if err != nil {
		log.Fatalln(err)
		return CloudBoltOrder{}, err
	}

	json.NewDecoder(resp.Body).Decode(&order)

	return order, nil
}
