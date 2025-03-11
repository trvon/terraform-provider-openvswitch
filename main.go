package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/trvon/terraform-provider-openvswitch/openvswitch"
)

func main() {
	// The binary works for both Terraform and OpenTofu - they use the same plugin interface
	// To detect if it's being called from OpenTofu, you could look at environment variables
	// or command line args, but we don't need to do anything differently in this case

	// Both Terraform and OpenTofu will properly load and use the same provider schema
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: openvswitch.Provider,
	})
}
