# CloudBolt Provider Examples

This directory contains a set of examples of using various CloudBolt resources and data sources with Terraform. The examples each have their own README containing more details on what the example does.

To run any example, clone the repository, set the variables (see Setting Variables below fir details) and run terraform apply within the example's own directory.

For example:

```
$ git clone https://github.com/CloudBoltSoftware/terraform-provider-cloudbolt.git
$ cd terraform-provider-cloudbolt/examples/order-blueprint
$ terraform apply
...
## Set Variables

To run an example you will want to copy `terraform.tfvars.dist` to `terraform.tfvars` in each example and fill each in with the appropriate variables.

Each example outlines additional prerequesties that is expected on the CloudBolt server for it to run correctly.

## Troublehsooting

Some common errors that are hard to troubleshoot:

* The user Terraform is using to make API calls does not have the appropriate permissions to access a group, blueprint, resource, etc. Sometimes CloudBolt returns empty results instead of a permission denied.
