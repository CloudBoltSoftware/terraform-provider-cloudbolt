# Troubleshooting

You may find that the Terraform Provider doesn't always provide the most useful error messages.
Here are some common errors and their workarounds.

## unexpected state 'CART'

```text
$ terraform destroy
... snip ...
cloudbolt_bp_instance.mycbresource: Destroying... [id=/api/v2/resources/service/35/]
cloudbolt_bp_instance.mycbresource: Still destroying... [id=/api/v2/resources/service/35/, 10s elapsed]

Error: Error waiting for Job (/api/v2/orders/81/) to complete: unexpected state 'CART', wanted target 'SUCCESS'. last error: %!s(<nil>)
```

This means that the "Delete" Resource Action requires approval in CloudBolt.

To resolve this...

1. Go to `https://<your-cloudbolt-instance>/actions/resource_actions/`
2. Edit the "Delete" action.
3. De-select "Requires Approval".

You may also want to log in to CloudBolt as the user Terraform is running as and clear your Cart.

## Parameter named 'some_input_name' does not exist

```text
$ terraform apply
... snip ...
cloudbolt_bp_instance.mycbresource: Creating...

Error: received a server error: {"status_code":500,"error":"Parameter named 'some_input_name' does not exist","detail":"An unexpected error occurred."}
```

This means that a Parameter name used in the Terraform Plan does not match what is in CloudBolt.

This can happen for a few reasons, but most often it is a result of a Plugin's generated parameter name differing from the expected name.

1. Navigate to `https://<your-cloudbolt-instance>/actions/{action_id}/`.
2. Edit the Action Inputs so their "Name" matches the one used in the Terraform Plan.
