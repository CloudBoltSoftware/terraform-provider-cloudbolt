# terraform-provider-cloudbolt
Sample Terraform Resource Provider to initiate CloudBolt Blueprint Orders.

## Supported versions of CloudBolt
The Terraform Provider supports version CloudBolt >= 9.1

## Prerequisites
- Install and Configure `golang >= 1.12`
- Install and Configure `terraform <= 0.11`

## Installation

```sh
go get github.com/cloudboltsoftware/terraform-provider-cloudbolt

cd ${GOPATH}/src/github.com/cloudboltsoftware/terraform-provider-cloudbolt 

make install
```

This will build the Terraform Provider and copy it to your `~/.terraform.d/plugins` directory.

## Sample Terraform Configuration

```hcl
provider "cloudbolt" {
    cb_protocol = "https"   // (Optional | Default: https) API protocol
    cb_host = "localhost"   // CloudBolt host
    cb_port = "8443"        // CloudBolt port
    cb_username = "cbadmin" // API user
    cb_password = "cbadmin" // API password
    cb_api_version = "v2"   // (Optional | Default: v2) Which version of the API to use.
    cb_insecure = false     // (Optional | Default: false) Disable SSL verification
    cb_timeout = 10         // (Optional | Default: 10) HTTP timeout
}

data "cloudbolt_group_ref" "group" {
    name = "/My Org/Dev Team 1"
}

data "cloudbolt_object_ref" "blueprint" {
    name = "CentOS_7"
    type = "blueprints"
}

resource "cloudbolt_bp_instance" "mycbresource" {
    group = "${data.cloudbolt_group_ref.group.url_path}"
    blueprint = "${data.cloudbolt_object_ref.blueprint.url_path}"
    blueprint_item = {
        name = "build-item-Build_VM"
        parameters = {
            cpu-cnt = "1",
            mem-size = "1 GB",
            placement-tag = "simulated_vmware",
        }
    }
}
```

## Sample terraform apply

```sh
➜  terraform-sample ls -ltr
total 8
-rw-r--r--  1 myuser staff  709 Jul 23 16:01 main.tf
➜  terraform-sample terraform apply
data.cloudbolt_object_ref.blueprint: Refreshing state...
data.cloudbolt_group_ref.group: Refreshing state...

An execution plan has been generated and is shown below.
Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  + cloudbolt_bp_instance.mycbresource
      id:                                                <computed>
      blueprint:                                         "/api/v2/blueprints/2/"
      blueprint_item.#:                                  "1"
      blueprint_item.479288173.name:                     "build-item-Build_VM"
      blueprint_item.479288173.parameters.%:             "3"
      blueprint_item.479288173.parameters.cpu-cnt:       "1"
      blueprint_item.479288173.parameters.mem-size:      "1 GB"
      blueprint_item.479288173.parameters.placement-tag: "simulated_vmware"
      group:                                             "/api/v2/groups/3/"
      server_hostname:                                   <computed>
      server_ip:                                         <computed>
      servers.#:                                         <computed>


Plan: 1 to add, 0 to change, 0 to destroy.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

cloudbolt_bp_instance.mycbresource: Creating...
  blueprint:                                         "" => "/api/v2/blueprints/2/"
  blueprint_item.#:                                  "0" => "1"
  blueprint_item.479288173.name:                     "" => "build-item-Build_VM"
  blueprint_item.479288173.parameters.%:             "0" => "3"
  blueprint_item.479288173.parameters.cpu-cnt:       "" => "1"
  blueprint_item.479288173.parameters.mem-size:      "" => "1 GB"
  blueprint_item.479288173.parameters.placement-tag: "" => "simulated_vmware"
  group:                                             "" => "/api/v2/groups/3/"
  server_hostname:                                   "" => "<computed>"
  server_ip:                                         "" => "<computed>"
  servers.#:                                         "" => "<computed>"
cloudbolt_bp_instance.mycbresource: Still creating... (10s elapsed)
cloudbolt_bp_instance.mycbresource: Creation complete after 19s (ID: /api/v2/resources/service/7/)

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.
```

## Sample terraform show

```sh
➜  terraform-sample terraform show
cloudbolt_bp_instance.mycbresource:
  id = /api/v2/resources/service/7/
  blueprint = /api/v2/blueprints/2/
  blueprint_item.# = 1
  blueprint_item.479288173.name = build-item-Build_VM
  blueprint_item.479288173.parameters.% = 3
  blueprint_item.479288173.parameters.cpu-cnt = 1
  blueprint_item.479288173.parameters.mem-size = 1 GB
  blueprint_item.479288173.parameters.placement-tag = simulated_vmware
  group = /api/v2/groups/3/
  server_hostname = fakesvr-1563928164.2307684
  server_ip =
  servers.# = 1
  servers.0.hostname = fakesvr-1563928164.2307684
  servers.0.ip =
data.cloudbolt_group_ref.group:
  id = 3
  name = /My Org/Dev Team 1
  url_path = /api/v2/groups/3/
data.cloudbolt_object_ref.blueprint:
  id = 2
  name = CentOS_7
  type = blueprints
  url_path = /api/v2/blueprints/2/
```

## Sample terraform destroy

```sh
  ➜  terraform-sample terraform destroy
data.cloudbolt_object_ref.blueprint: Refreshing state...
data.cloudbolt_group_ref.group: Refreshing state...
cloudbolt_bp_instance.mycbresource: Refreshing state... (ID: /api/v2/resources/service/7/)

An execution plan has been generated and is shown below.
Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  - cloudbolt_bp_instance.mycbresource


Plan: 0 to add, 0 to change, 1 to destroy.

Do you really want to destroy all resources?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes

cloudbolt_bp_instance.mycbresource: Destroying... (ID: /api/v2/resources/service/7/)
cloudbolt_bp_instance.mycbresource: Still destroying... (ID: /api/v2/resources/service/7/, 10s elapsed)
cloudbolt_bp_instance.mycbresource: Destruction complete after 10s

Destroy complete! Resources: 1 destroyed.
```
