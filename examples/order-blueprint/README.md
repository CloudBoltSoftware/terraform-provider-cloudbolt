# Order Blueprint test

This test orders a simple Terraform Blueprint.
The blueprint does not contain any server builds, so this can be run without incurring the cost of provisioning new resources.

## Prerequisites

### Import Sample blueprint.

* Download the `Terraform Provider Sample Blueprint` blueprint.
* Import the `Terraform Provider Sample Blueprint` blueprint to the target CloudBolt Instance.
* Add Deploument group to the the `Terraform Provider Sample Blueprint` blueprint.


### Terraform Variables

Copy the file `terraform.tfvars.dist` to `terraform.tfvars.dist`, update to match the targeted instance of CloudBolt.

## Order Blueprint
To create and submit an order for the sample bluepint, execute the following.

```
terraform apply
```
