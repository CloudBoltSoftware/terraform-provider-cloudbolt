---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "cloudbolt_1f_ipam_policy Data Source - terraform-provider-cloudbolt"
subcategory: "OneFuse"
description: |-
  
---

# cloudbolt_1f_ipam_policy (Data Source)

Use this data source to retreive reference information for a OneFuse IPAM Policy.

## Example Usage
```hcl
// Data Source for IPAM Policy - lookup Policy ID by Policy Name
data "cloudbolt_1f_ipam_policy" "ipam_policy" {
  name = "my_policy_name"  // Replace with Policy Name
}
```

<!-- schema generated by tfplugindocs -->
## Argument Reference

### Required

- `name` (String) The name of the OneFuse IPAM Policy

### Read-Only

- `description` (String) The description for the OneFuse IPAM Policy
- `id` (String) The ID of the OneFuse IPAM Policy


