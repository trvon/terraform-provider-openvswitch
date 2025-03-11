package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/carletes/terraform-provider-openvswitch/openvswitch"
	"github.com/hashicorp/terraform/plugin"
	otfplugin "github.com/opentofu/opentofu/plugin"
)

func main() {
	// Check for OpenTofu in the process name
	isOpenTofu := false
	for _, arg := range os.Args {
		if strings.Contains(strings.ToLower(arg), "tofu") {
			isOpenTofu = true
			break
		}
	}

	if isOpenTofu {
		otfplugin.Serve(&otfplugin.ServeOpts{
			ProviderFunc: openvswitch.ProviderOpenTofu(),
		})
	} else {
		plugin.Serve(&plugin.ServeOpts{
			ProviderFunc: openvswitch.Provider,
		})
	}
}