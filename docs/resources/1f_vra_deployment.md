---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "cloudbolt_1f_vra_deployment Resource - terraform-provider-cloudbolt"
subcategory: "OneFuse"
description: |-
  
---

# cloudbolt_1f_vra_deployment (Resource)

Provides a OneFuse vRealize Automation resource to deploy and destroy Deployments in vRealize Automation.

## Example Usage

### Variables
```hcl
variable "template_properties" {
  type = map
  default = {
    "Environment" = "d"     //p for production or d for development
    "Location"    = "w"     //e for East or w for West
    "Application" = "wp"    //wp for wordpress or iis for IIS
    "OS"          = "l"     //l for Linux or w for Windows
  }
}
```

### Terraform Configuration
```hcl
// Data Source for vRealize Automation Policy - lookup Policy ID by Policy Name
data "cloudbolt_1f_vra_policy" "vra_policy" {
  name = "my_policy_name"  // Replace with Policy Name
}

// Resource for vRealize Automation
resource "cloudbolt_1f_vra_deployment" "my_vra_deployment" {
  policy_id           = data.onefuse_vra_policy.vra_policy.id  // Refers to Data Source to get Policy ID
  workspace_url       = var.workspace_url
  template_properties = var.template_properties
  deployment_name     = "tf_vra_deployment"
  request_timeout     = 20
}
```

<!-- schema generated by tfplugindocs -->
## Argument Reference

### Required

- `deployment_name` (String) Name of the vRA Deployment.
- `policy_id` (Number) OneFuse Module Policy ID.

### Optional

- `blueprint_name` (String)
- `deployment_info` (String)
- `id` (String) The ID of this resource.
- `project_name` (String)
- `request_timeout` (Number) Timeout in minutes, Default (30)
- `template_properties` (Map of String) Additional properties that are referenced within the Policy.
- `workspace_url` (String) OneFuse Workspace URL path.


