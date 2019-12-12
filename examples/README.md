# terraform-provider-cloudbolt examples

These examples demonstrate how to use the CloudBolt provider.

They are periodically run to verify the Provider is compatible with the latest version of CloudBolt.

## Usage

Navigate to any of the directories here and run `make run` to `terraform init`, `terraform plan` and `terraform apply`.
Run `make destroy` to teardown.

## Infrastructure

To run an example you will want to copy `terraform.tfvars.dist` to `terraform.tfars` in each example and fill each in with the appropriate variables.

Each example outlines additional infrastructure that is expected on the CloudBolt server for it to run correctly.

## Troublehsooting

Some common errors that are hard to troubleshoot:

* The user Terraform is using to make API calls does not have the appropriate permissions to access a group, blueprint, resource, etc. Sometimes CloudBolt returns empty results instead of a permission denied.
