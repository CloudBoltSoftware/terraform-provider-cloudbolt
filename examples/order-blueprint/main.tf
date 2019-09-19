provider "cloudbolt" {
    cb_protocol = "http"
    cb_host = "localhost"
    cb_port = "8000"
    cb_username = "userly"
    cb_password = "cloudbolt"
}

data "cloudbolt_group_ref" "group" {
    name = "organizationly"
}

data "cloudbolt_object_ref" "blueprint" {
    name = "Just Bools"
    type = "blueprints"
}

resource "cloudbolt_bp_instance" "mycbresource" {
    group = "${data.cloudbolt_group_ref.group.url_path}"
    blueprint = "${data.cloudbolt_object_ref.blueprint.url_path}"
    blueprint_item = {
        name = "bools 2: the boolening"
        parameters = { }
    }
}
