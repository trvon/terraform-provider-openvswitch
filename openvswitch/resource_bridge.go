package openvswitch

import (
	"github.com/digitalocean/go-openvswitch/ovs"
	"github.com/hashicorp/terraform/helper/schema"
	otfschema "github.com/opentofu/opentofu/helper/schema"
)

// Resource Definition
func resourceBridge() *schema.Resource {
	return &schema.Resource{
		Create: resourceBridgeCreate,
		Read:   resourceBridgeRead,
		Update: resourceBridgeUpdate,
		Delete: resourceBridgeDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ofversion": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "OpenFlow13",
			},
		},
	}
}

func resourceBridgeCreate(d *schema.ResourceData, m interface{}) error {
	bridge := d.Get("name").(string)
	ver := []string{d.Get("ofversion").(string)}
	bridge_options := ovs.BridgeOptions{ver}

	err := c.VSwitch.AddBridge(bridge)
	err = c.VSwitch.Set.Bridge(bridge, bridge_options)

	return err
}

func resourceBridgeRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceBridgeUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceBridgeRead(d, m)
}

func resourceBridgeDelete(d *schema.ResourceData, m interface{}) error {
	bridge := d.Get("name").(string)
	return c.VSwitch.DeleteBridge(bridge)
}

// OpenTofu resource implementation
func resourceBridgeOTF() *otfschema.Resource {
	return &otfschema.Resource{
		Create: resourceBridgeCreateOTF,
		Read:   resourceBridgeReadOTF,
		Update: resourceBridgeUpdateOTF,
		Delete: resourceBridgeDeleteOTF,

		Schema: map[string]*otfschema.Schema{
			"name": {
				Type:     otfschema.TypeString,
				Required: true,
			},
			"ofversion": {
				Type:     otfschema.TypeString,
				Optional: true,
				Default:  "OpenFlow13",
			},
		},
	}
}

func resourceBridgeCreateOTF(d *otfschema.ResourceData, m interface{}) error {
	bridge := d.Get("name").(string)
	ver := []string{d.Get("ofversion").(string)}
	bridge_options := ovs.BridgeOptions{ver}

	err := c.VSwitch.AddBridge(bridge)
	err = c.VSwitch.Set.Bridge(bridge, bridge_options)

	return err
}

func resourceBridgeReadOTF(d *otfschema.ResourceData, m interface{}) error {
	return nil
}

func resourceBridgeUpdateOTF(d *otfschema.ResourceData, m interface{}) error {
	return resourceBridgeReadOTF(d, m)
}

func resourceBridgeDeleteOTF(d *otfschema.ResourceData, m interface{}) error {
	bridge := d.Get("name").(string)
	return c.VSwitch.DeleteBridge(bridge)
}
