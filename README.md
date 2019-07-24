# terraform-provider-cloudbolt
Sample Terraform Resource Provider to initiate CloudBolt Blueprint Orders.

## Prerequisites
- Install and Configure golang
- Install and Configure terraform

## Installation
```go
go get github.com/laltomar/cloudbolt-go-sdk
go get github.com/laltomar/terraform-provider-cloudbolt

cd ${GOPATH}/src/github.com/laltomar/terraform-provider-cloudbolt 

mkdir ~/.terraform.d/plugins

go build -o terraform-provider-cloudbolt
mv terraform-provider-cloudbolt ~/.terraform.d/plugins/.

```
## Sample Terraform Configuration
```go
provider "cloudbolt" {
    cb_protocol = "https"
    cb_host = "localhost"
    cb_port = "8443"
    cb_username = "cbadmin"
    cb_password = "cbadmin"
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
```go
➜  terraform-sample ls -ltr
total 8
-rw-r--r--  1 laltomar  staff  709 Jul 23 16:01 main.tf
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
