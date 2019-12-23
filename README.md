# terraform-provider-cloudbolt
Sample Terraform Resource Provider to initiate CloudBolt Blueprint Orders.

## Supported versions of CloudBolt
This Terraform Provider officially supports CloudBolt >= 9.1, but may work with earlier version of the product.
Your mileage may vary.

## Prerequisites
- Install and Configure `golang >= 1.12`
- Install and Configure `terraform >= 0.12` (_may_ work with `terraform <= 0.11.x`)

## Installation

```sh
# Fetch the Terraform Provider
go get github.com/cloudboltsoftware/terraform-provider-cloudbolt

# Navigate to the Provider
cd ${GOPATH}/src/github.com/cloudboltsoftware/terraform-provider-cloudbolt

# Compile the binary and place it in ~/.terraform.d/plugins/
make install
```

This will build the Terraform Provider and copy it to your `~/.terraform.d/plugins` directory.

## Sample Terraform Configuration

To get started with the Terraform Provider for CloudBolt, put the following into a file called `main.tf`.

Fill in the `provider "cloudbolt"` section with details about your CloudBolt instance.

```hcl
provider "cloudbolt" {
  cb_protocol =    "https"      // (Optional | Default: https) API protocol
  cb_host =        "localhost"  // CloudBolt host
  cb_port =        "8443"       // CloudBolt port
  cb_username =    "cbadmin"    // API user
  cb_password =    "cbadmin"    // API password
  cb_api_version = "v2"         // (Optional | Default: v2) Which version of the API to use.
  cb_insecure =    false        // (Optional | Default: false) Disable SSL verification
  cb_timeout =     10           // (Optional | Default: 10) HTTP timeout in seconds
}

data "cloudbolt_group_ref" "group" {
    name = "/My Org/Dev Team 1"
}

data "cloudbolt_object_ref" "blueprint" {
    name = "CentOS_7"
    type = "blueprints"
}

resource "cloudbolt_bp_instance" "mycbresource" {
    group = data.cloudbolt_group_ref.group.url_path
    blueprint = data.cloudbolt_object_ref.blueprint.url_path
    blueprint_item {
        name = "build-item-Build_VM"
        parameters = {
            cpu-cnt = "1"
            mem-size = "1 GB"
            placement-tag = "simulated_vmware"
        }
    }
}
```

## Sample terraform plan

```sh
[cb-terraform-example] $  ls -ltr
total 8
-rw-r--r--  1 myuser staff  709 Jul 23 16:01 main.tf
$ terraform plan
Refreshing Terraform state in-memory prior to plan...
The refreshed state will be used to calculate this plan, but will not be
persisted to local or remote state storage.

data.cloudbolt_object_ref.blueprint: Refreshing state...
data.cloudbolt_group_ref.group: Refreshing state...

------------------------------------------------------------------------

An execution plan has been generated and is shown below.
Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # cloudbolt_bp_instance.mycbresource will be created
  + resource "cloudbolt_bp_instance" "mycbresource" {
      + blueprint       = "/api/v2/blueprints/23/"
      + group           = "/api/v2/groups/6/"
      + id              = (known after apply)
      + instance_type   = (known after apply)
      + server_hostname = (known after apply)
      + server_ip       = (known after apply)
      + servers         = (known after apply)

      + blueprint_item {
          + name       = "build-item-Build_VM"
          + parameters = {
              + "cpu-cnt" = "1"
              + "mem-size" = "1 GB"
              + "placement-tag" = "simulated_vmware"
            }
        }
    }
```

# Sample terraform apply

```sh
[cb-terraform-example] $  terraform apply
data.cloudbolt_object_ref.blueprint: Refreshing state...
data.cloudbolt_group_ref.group: Refreshing state...

An execution plan has been generated and is shown below.
Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # cloudbolt_bp_instance.mycbresource will be created
  + resource "cloudbolt_bp_instance" "mycbresource" {
      + blueprint       = "/api/v2/blueprints/23/"
      + group           = "/api/v2/groups/6/"
      + id              = (known after apply)
      + instance_type   = (known after apply)
      + server_hostname = (known after apply)
      + server_ip       = (known after apply)
      + servers         = (known after apply)

      + blueprint_item {
          + name       = "build-item-Build_VM"
          + parameters = {
              + "cpu-cnt" = "1"
              + "mem-size" = "1 GB"
              + "placement-tag" = "simulated_vmware"
            }
        }
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

cloudbolt_bp_instance.mycbresource: Creating...
cloudbolt_bp_instance.mycbresource: Still creating... [10s elapsed]
cloudbolt_bp_instance.mycbresource: Creation complete after 13s [id=/api/v2/resources/service/25/]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.
```

## Sample terraform show

```sh
[cb-terraform-example] $ terraform show
# cloudbolt_bp_instance.mycbresource:
resource "cloudbolt_bp_instance" "mycbresource" {
    blueprint     = "/api/v2/blueprints/23/"
    group         = "/api/v2/groups/6/"
    id            = "/api/v2/resources/service/25/"
    instance_type = "Resource"
    servers       = []

    blueprint_item {
        name       = "build-item-Build_VM"
        parameters = {
            "cpu-cnt" = "1"
            "mem-size" = "1 GB"
            "placement-tag" = "simulated_vmware"
        }
    }
}

# data.cloudbolt_group_ref.group:
data "cloudbolt_group_ref" "group" {
    id       = "6"
    name     = "/My Org/Dev Team 1"
    url_path = "/api/v2/groups/6/"
}

# data.cloudbolt_object_ref.blueprint:
data "cloudbolt_object_ref" "blueprint" {
    id       = "23"
    name     = "CentOS_7"
    type     = "blueprints"
    url_path = "/api/v2/blueprints/23/"
}
```

## Sample terraform destroy

```sh
[cb-terraform-example] $ terraform destroy
data.cloudbolt_object_ref.blueprint: Refreshing state...
data.cloudbolt_group_ref.group: Refreshing state...
cloudbolt_bp_instance.mycbresource: Refreshing state... [id=/api/v2/resources/service/25/]

An execution plan has been generated and is shown below.
Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  # cloudbolt_bp_instance.mycbresource will be destroyed
  - resource "cloudbolt_bp_instance" "mycbresource" {
      - blueprint     = "/api/v2/blueprints/23/" -> null
      - group         = "/api/v2/groups/6/" -> null
      - id            = "/api/v2/resources/service/25/" -> null
      - instance_type = "Resource" -> null
      - servers       = [] -> null

      - blueprint_item {
          - name       = "build-item-Build_VM" -> null
          - parameters = {
              - "cpu-cnt" = "1"
              - "mem-size" = "1 GB"
              - "placement-tag" = "simulated_vmware"
            } -> null
        }
    }

Plan: 0 to add, 0 to change, 1 to destroy.

Do you really want to destroy all resources?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes

cloudbolt_bp_instance.mycbresource: Destroying... [id=/api/v2/resources/service/25/]
cloudbolt_bp_instance.mycbresource: Still destroying... [id=/api/v2/resources/service/25/, 10s elapsed]
cloudbolt_bp_instance.mycbresource: Destruction complete after 12s

Destroy complete! Resources: 1 destroyed.
```

## Troubleshooting

If you get stuck or confused using the Terraform Provider, check out the TROUBLESHOOTING document in this repo.

If you can't find an answer and believe you have found a limitation or bug, please make a GitHub issue!
