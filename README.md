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
