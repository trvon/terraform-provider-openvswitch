package openvswitch

import (
	"github.com/digitalocean/go-openvswitch/ovs"
	"github.com/hashicorp/terraform/helper/schema"
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
	bridge_options := ovs.BridgeOptions{Protocols: ver}

	if err := c.VSwitch.AddBridge(bridge); err != nil {
		return err
	}

	if err := c.VSwitch.Set.Bridge(bridge, bridge_options); err != nil {
		return err
	}

	// Set the ID to the bridge name to ensure Terraform can track the resource
	d.SetId(bridge)
	return resourceBridgeRead(d, m)
}

func resourceBridgeRead(d *schema.ResourceData, m interface{}) error {
	bridge := d.Id()

	// Check if bridge exists by attempting to list its ports
	// If bridge doesn't exist, this will return an error
	_, err := c.VSwitch.ListPorts(bridge)
	if err != nil {
		// Bridge doesn't exist, remove from state
		d.SetId("")
		return nil
	}

	// Bridge exists, set attributes
	d.Set("name", bridge)

	// Keep the ofversion attribute in the state if it's already there
	// This ensures we don't lose attributes after apply
	if ofversion, ok := d.GetOk("ofversion"); ok {
		d.Set("ofversion", ofversion)
	}

	return nil
}

func resourceBridgeUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceBridgeRead(d, m)
}

func resourceBridgeDelete(d *schema.ResourceData, m interface{}) error {
	bridge := d.Get("name").(string)
	return c.VSwitch.DeleteBridge(bridge)
}
