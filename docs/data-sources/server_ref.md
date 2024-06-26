---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "cloudbolt_server_ref Data Source - terraform-provider-cloudbolt"
subcategory: "Cloud Management Platform"
description: |-
  
---

# cloudbolt_server_ref (Data Source)

Use this data source to retreive reference information for a CloudBolt Server by ID, Hostname, or API URL path.

## Example Usage
```hcl
data "cloudbolt_server_ref" "server" {
    url_path = "/api/v3/cmp/servers/SVR-d5y4u971/"
}

data "cloudbolt_server_ref" "server_id" {
    id = "SVR-d5y4u971"
}

data "cloudbolt_server_ref" "server_name" {
    hostname = "myhostname"
}
```

<!-- schema generated by tfplugindocs -->
## Argument Reference

### Optional

- `id` (String) The global id of a CloudBolt Server, required if "hostname" or "url_path" is not provided
- `hostname` (String) The hostname of the CloudBolt Server, required if "id" or "url_path" is not provided
- `url_path` (String) The relative API URL path for the CloudBolt Server, required if "id" or "hostname" is not provided

### Read-Only

- `ip_address` (String) Server IP Address
- `status` (String) CloudBolt Server Status
- `mac` (String) Server MAC Address
- `power_status` (String) Server Power Status
- `cpu_count` (Number) CPU Count
- `memory_size_gb` (String) Total Memory in GB
- `disk_size_gb` (Number) Total Disk Size in GB
- `date_added_to_cloudbolt` (String) Date the server was added to CloudBolt
- `notes` (String) Server Notes
- `labels` (List of String) Server Labels
- `os_family` (String) Server OS Family
- `attributes` (Map of String) CloudBolt Rsource attributes
- `rate_breakdown` (Map of String) Server Rate Breakdown
- `tech_specific_attributes` (Map of String) Resource Handler technical specific attributes
- `disks` (List of Object) "Server disks (see [below for nested schema](#nestedobjatt--servers--disks))
- `networks` (List of Map of String) Server NICs
 
<a id="nestedobjatt--servers--disks"></a>
### Nested Schema for `servers.disks`

Read-Only:

- `disk_size_gb` (Number) Disk Size in GB
- `name` (String) Name of Disk
- `uuid` (String) Unique ID of Disk