terraform {
  required_providers {
    cloudbolt = {
      source  = "hashicorp.com/cbsw/cloudbolt"
      version = "1.0.0"
    }
  }
}

provider "cloudbolt" {
  cb_protocol =    var.CB_PROTOCOL
  cb_host =        var.CB_HOST
  cb_port =        var.CB_PORT
  cb_username =    var.CB_USERNAME
  cb_password =    var.CB_PASSWORD
  cb_insecure =    true
  cb_timeout =     500
}

data "cloudbolt_group_ref" "group" {
    name = var.CB_GROUP
}

data "cloudbolt_blueprint_ref" "blueprint" {
    name = "Terraform Provider Sample Blueprint"
}

resource "cloudbolt_bp_instance" "cb_terraform_sample_resource" {
  group = data.cloudbolt_group_ref.group.url_path
  blueprint_id = data.cloudbolt_blueprint_ref.blueprint.id

  deployment_item {
    name = "plugin-bdi-buphbggq"
    parameters = {
      tf_sample_param1: "[hello|world]",
      tf_sample_param2: 100,
      tf_sample_param3: false,
      tf_sample_param4: 101.1
    }
  }
}
