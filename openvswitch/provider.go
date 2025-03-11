package openvswitch

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	otfschema "github.com/opentofu/opentofu/helper/schema"
	otfplugin "github.com/opentofu/opentofu/plugin"
	otftf "github.com/opentofu/opentofu/terraform"
)

// Provider returns a schema.Provider for OpenVSwitch.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{},

		ResourcesMap: map[string]*schema.Resource{
			"openvswitch_bridge": resourceBridge(),
			"openvswitch_port":   resourcePort(),
		},

		DataSourcesMap: map[string]*schema.Resource{},
	}
}

// ProviderOpenTofu returns a schema.Provider for OpenVSwitch for OpenTofu.
func ProviderOpenTofu() otfplugin.ProviderFunc {
	return func() otftf.ResourceProvider {
		return &otfschema.Provider{
			Schema: map[string]*otfschema.Schema{},

			ResourcesMap: map[string]*otfschema.Resource{
				"openvswitch_bridge": resourceBridgeOTF(),
				"openvswitch_port":   resourcePortOTF(),
			},

			DataSourcesMap: map[string]*otfschema.Resource{},
		}
	}
}