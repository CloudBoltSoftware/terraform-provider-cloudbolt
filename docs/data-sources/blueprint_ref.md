---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "cloudbolt_blueprint_ref Data Source - terraform-provider-cloudbolt"
subcategory: "Cloud Management Platform"
description: |-
  
---

# cloudbolt_blueprint_ref (Data Source)

Use this data source to retreive reference information for a CloudBolt Blueprint by name or ID. 

## Example Usage
```hcl
data "cloudbolt_blueprint_ref" "blueprint" {
    name = "My Blueprint"
}

data "cloudbolt_blueprint_ref" "blueprint_id" {
    id = "BP-abcd1234"
}
```

<!-- schema generated by tfplugindocs -->
## Argument Reference

### Optional

- `id` (String) The global id of a CloudBolt Blueprint, required if "name" not provided
- `name` (String) The name of a CloudBolt Blueprint, required if "id" not provided

### Read-Only

- `url_path` (String) The relative API URL path for the CloudBolt Blueprint.


