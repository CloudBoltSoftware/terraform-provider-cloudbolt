# Order Blueprint test

This test orders a simple Terraform Blueprint.
The blueprint does not contain any server builds, so this can be run without incurring the cost of provisioning new resources.

## Prerequisites

This requires a few resources to exist on your CloudBolt instance:

* A CloudBolt host running at "http://localhost:8000"
* A CloudBolt environment named "TerraformEnvironment01"
* A ClouBolt group named "TerraformGroup01"
* A CloudBolt user named "TerraformUser01" with API access and a member of TerraformGroup01.
* A shared CloudBolt Plugin name "TerraformAction01"
* A CloudBolt Catalog Item named "TerraformCatalogItem01" deployable by TerraformGroup01 and with the following parameters:
  * Param1 (String)
  * Param2 (Integer)
  * Param3 (Boolean)
  * Param4 (Multi-select string with options 'a' 'b' and 'c')
  * Param5 (Ranged Integer from 1 to 10)
  * Param6 (Optional String)

The Username, Password, Group, and Environment can be changed by passing different values for their corresponding Terraform Variables.
For more information, see the Terraform Variables documentation: https://www.terraform.io/docs/configuration-0-11/variables.html#environment-variables

## Terraform Variables

The file `terraform.tfvars.dist` contains the necessary varaible inputs with values listed above.
Feel free to change these, or pass your own variables to `terraform` with `-var 'var=value'` to override them.

e.g., if you want to run this example on your CloudBolt instance, you might change `CB_PROTOCOL` to `https`, `CB_HOST` to `my.cb.internal` and `CB_PORT` to `443`.

The only part of this example which is not configurable is the name of the Blueprint and Action.

## TeraformAction01

An example of the Action which generates the above 6 parameters is as follows:

```py
from common.methods import set_progress


def run(job, *args, **kwargs):
    set_progress("Running Terarform Action 01")

    p1 = {{ param1 }}
    p2 = {{ param2 }}
    p3 = {{ param3 }}
    p4 = {{ param4 }}
    p5 = {{ param5 }}
    p6 = {{ param6 }}

    set_progress("Variables recieved:")
    set_progress(f"Parm 1: { p1 }")
    set_progress(f"Parm 2: { p2 }")
    set_progress(f"Parm 3: { p3 }")
    set_progress(f"Parm 4: { p4 }")
    set_progress(f"Parm 5: { p5 }")
    set_progress(f"Parm 6: { p6 }")
```

Note that the auto-generated Action Input Parameters will have names like `param1_a116`.
These names should be manually cleaned up to be `param1` to make them easier to use in the API.

## TerraformUser01

The Terraform User should have the password "TerraformPassword01", but this can be overridden.
