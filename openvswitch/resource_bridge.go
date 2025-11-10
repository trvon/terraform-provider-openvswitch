package openvswitch

import (
	"fmt"

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
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the bridge to create",
			},
			"ofversion": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "OpenFlow13",
				ValidateFunc: func(v interface{}, k string) (warnings []string, errors []error) {
					value, ok := v.(string)
					if !ok {
						errors = append(errors, fmt.Errorf("%q must be a string", k))
						return warnings, errors
					}
					validVersions := []string{
						"OpenFlow10",
						"OpenFlow11",
						"OpenFlow12",
						"OpenFlow13",
						"OpenFlow14",
						"OpenFlow15",
					}
					for _, version := range validVersions {
						if value == version {
							return nil, nil
						}
					}
					errors = append(errors, fmt.Errorf(
						"%q must be one of: %v", k, validVersions))
					return warnings, errors
				},
				Description: "OpenFlow protocol version (OpenFlow10, OpenFlow11, OpenFlow12, OpenFlow13, OpenFlow14, or OpenFlow15)",
			},
		},
	}
}

func resourceBridgeCreate(d *schema.ResourceData, m interface{}) error {
	bridge, ok := d.Get("name").(string)
	if !ok {
		return fmt.Errorf("name must be a string")
	}

	ofversion, ok := d.Get("ofversion").(string)
	if !ok {
		return fmt.Errorf("ofversion must be a string")
	}

	ver := []string{ofversion}
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
	if err := d.Set("name", bridge); err != nil {
		return fmt.Errorf("error setting name: %w", err)
	}

	// Keep the ofversion attribute in the state if it's already there
	// This ensures we don't lose attributes after apply
	if ofversion, ok := d.GetOk("ofversion"); ok {
		if err := d.Set("ofversion", ofversion); err != nil {
			return fmt.Errorf("error setting ofversion: %w", err)
		}
	}

	return nil
}

func resourceBridgeUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceBridgeRead(d, m)
}

func resourceBridgeDelete(d *schema.ResourceData, m interface{}) error {
	bridge, ok := d.Get("name").(string)
	if !ok {
		return fmt.Errorf("name must be a string")
	}
	return c.VSwitch.DeleteBridge(bridge)
}
