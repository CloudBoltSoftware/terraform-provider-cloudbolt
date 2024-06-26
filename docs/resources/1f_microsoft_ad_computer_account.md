---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "cloudbolt_1f_microsoft_ad_computer_account Resource - terraform-provider-cloudbolt"
subcategory: "OneFuse"
description: |-
  
---

# cloudbolt_1f_microsoft_ad_computer_account (Resource)

Provides a OneFuse Microsoft Active Directory resource to create and delete computer accounts.

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
// Data Source for Microsoft Active Directory Policy - lookup Policy ID by Policy Name
data "cloudbolt_1f_ad_policy" "ad_policy" {
  name = "my_policy_name"  // Replace with Policy Name
}

// Resource for Microsoft Active Directory
resource "cloudbolt_1f_microsoft_ad_computer_account" "my_ad_computer_account" {
  policy_id           = data.onefuse_ad_policy.ad_policy.id        // Refers to Data Source to get Policy ID
  workspace_url       = ""
  template_properties = var.template_properties
  name                = "computer_name"
}
```

<!-- schema generated by tfplugindocs -->
## Argument Reference

### Required

- `policy_id` (Number) OneFuse Module Policy ID.

### Optional

- `final_ou` (String)
- `id` (String) The ID of this resource.
- `name` (String) Computer Account Name.
- `request_timeout` (Number) Timeout in minutes, Default (30)
- `template_properties` (Map of String) Additional properties that are referenced within the Policy.
- `workspace_url` (String) OneFuse Workspace URL path.


