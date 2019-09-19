provider "cloudbolt" {
    cb_protocol = "${var.CB_PROTOCOL}"
    cb_host = "${var.CB_HOST}"
    cb_port = "${var.CB_PORT}"
    cb_username = "${var.CB_USERNAME}"
    cb_password = "${var.CB_PASSWORD}"
}

data "cloudbolt_group_ref" "group" {
    // Get the user-defined group
    name = "${var.CB_GROUP}"
}

data "cloudbolt_object_ref" "blueprint" {
    // Get the 'Custom Server' blueprint object
    name = "Custom Server"
    type = "blueprints"
}

data "cloudbolt_object_ref" "environment" {
    name = "(AWS) US East (Ohio) vpc-eb0ce282"
    type = "environments"
}

resource "random_id" "cb_server_hostname" {
    // Generate a random ID for the server hostname
    byte_length = "4"
    prefix = "cloudbolt-terraform-example-"
}

resource "cloudbolt_bp_instance" "cb_server" {
    group = "${data.cloudbolt_group_ref.group.url_path}"
    blueprint = "${data.cloudbolt_object_ref.blueprint.url_path}"

    blueprint_item = {
        name = "Server"
        environment = "${data.cloudbolt_object_ref.environment.url_path}"
        parameters = {
            // hostname = "${random_id.cb_server_hostname.b64_url}"
            ebs_volume_type = "gp2"
            instance_type = "t2.nano"
            key_name = "AMJ"
        }
    }
}
