package main

import (
	"github.com/cloudboltsoftware/terraform-provider-cloudbolt/cloudbolt"
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return cloudbolt.Provider()
		},
	})
}
