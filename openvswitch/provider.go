package openvswitch

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider returns a schema.Provider for OpenVSwitch.

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{},

		ResourcesMap: map[string]*schema.Resource{
			"openvswitch_bridge": resourceBridge(),
			"openvswitch_port":   resourcePort(),
		},

		DataSourcesMap: map[string]*schema.Resource{},
	}
}
