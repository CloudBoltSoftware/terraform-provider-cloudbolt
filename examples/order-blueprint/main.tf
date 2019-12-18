provider "cloudbolt" {
  cb_protocol = var.CB_PROTOCOL
  cb_host     = var.CB_HOST
  cb_port     = var.CB_PORT
  cb_username = var.CB_USERNAME
  cb_password = var.CB_PASSWORD
  cb_insecure = true
  cb_timeout  = 500
}

data "cloudbolt_group_ref" "group" {
  name = var.CB_GROUP
}

data "cloudbolt_object_ref" "blueprint" {
  name = "TerraformCatalogItem01"
  type = "blueprints"
}

data "cloudbolt_object_ref" "environment" {
  name = var.CB_ENVIRONMENT
  type = "environments"
}

resource "cloudbolt_bp_instance" "cb_order" {
  group     = data.cloudbolt_group_ref.group.url_path
  blueprint = data.cloudbolt_object_ref.blueprint.url_path

  blueprint_item {
    name        = "TerraformAction01"
    environment = data.cloudbolt_object_ref.environment.url_path
    parameters = {
      param1 = "TerraformInput"
      param2 = 9
      param3 = "foo"
      param4 = "c"
      param5 = 5
    }
  }
}

