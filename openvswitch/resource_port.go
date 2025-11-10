package openvswitch

import (
	"fmt"
	"log"
	"os/exec"
	"os/user"
	"strings"

	"github.com/digitalocean/go-openvswitch/ovs"
	"github.com/hashicorp/terraform/helper/schema"
)

// OVS Connection
var c = ovs.New(
	ovs.FlowFormat("OXM-OpenFlow14"),
	ovs.Protocols([]string{
		"OpenFlow10",
		"OpenFlow11",
		"OpenFlow12",
		"OpenFlow13",
		"OpenFlow14",
		"OpenFlow15",
	}),
	ovs.Sudo(),
)

// Resource Definition
func resourcePort() *schema.Resource {
	return &schema.Resource{
		Create: resourcePortCreate,
		Read:   resourcePortRead,
		Update: resourcePortUpdate,
		Delete: resourcePortDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the port to create",
			},

			"bridge_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the bridge to attach the port to",
			},

			"action": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "up",
				ValidateFunc: func(v interface{}, k string) (warnings []string, errors []error) {
					value, ok := v.(string)
					if !ok {
						errors = append(errors, fmt.Errorf("%q must be a string", k))
						return warnings, errors
					}
					validActions := []string{
						"up", "down", "stp", "no-stp", "receive", "no-receive",
						"no-receive-stp", "forward", "no-forward", "flood",
						"no-flood", "packet-in", "no-packet-in",
					}
					for _, action := range validActions {
						if value == action {
							return nil, nil
						}
					}
					errors = append(errors, fmt.Errorf(
						"%q must be one of: %v", k, validActions))
					return warnings, errors
				},
				Description: "Port action (up, down, stp, no-stp, receive, no-receive, no-receive-stp, forward, no-forward, flood, no-flood, packet-in, or no-packet-in)",
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

func GetPortAction(action string) ovs.PortAction {
	switch action {
	case ("up"):
		return ovs.PortActionUp
	case ("down"):
		return ovs.PortActionDown
	case ("stp"):
		return ovs.PortActionSTP
	case ("no-stp"):
		return ovs.PortActionNoSTP
	case ("receive"):
		return ovs.PortActionReceive
	case ("no-receive"):
		return ovs.PortActionNoReceive
	case ("no-receive-stp"):
		return ovs.PortActionReceiveSTP
	case ("forward"):
		return ovs.PortActionForward
	case ("no-forward"):
		return ovs.PortActionNoForward
	case ("flood"):
		return ovs.PortActionFlood
	case ("no-flood"):
		return ovs.PortActionNoFlood
	case ("packet-in"):
		return ovs.PortActionPacketIn
	case ("no-packet-in"):
		return ovs.PortActionNoPacketIn
	}
	return ovs.PortActionUp
}

func resourcePortCreate(d *schema.ResourceData, m interface{}) error {
	port, ok := d.Get("name").(string)
	if !ok {
		return fmt.Errorf("name must be a string")
	}

	bridge, ok := d.Get("bridge_id").(string)
	if !ok {
		return fmt.Errorf("bridge_id must be a string")
	}

	action, ok := d.Get("action").(string)
	if !ok {
		return fmt.Errorf("action must be a string")
	}

	// Creates tap device for ovs port, this is not persistent
	currentUser, err := user.Current()
	if err != nil {
		return fmt.Errorf("error getting current user: %w", err)
	}

	cmd := exec.Command("sudo", "/sbin/ip", "tuntap", "add", "dev", port, "mode", "tap", "user", currentUser.Username)
	if err := cmd.Run(); err != nil {
		log.Printf("warning: error creating tap device (may already exist): %v", err)
		// Continue even if there's an error, as the tap device might already exist
	}

	if err := c.VSwitch.AddPort(bridge, port); err != nil {
		return fmt.Errorf("error adding port to bridge: %w", err)
	}

	if err := c.OpenFlow.ModPort(bridge, port, GetPortAction(action)); err != nil {
		log.Printf("warning: error modifying port action: %v", err)
		// Continue even if ModPort fails
	}

	// Set the ID using bridge:port format to ensure Terraform can track the resource
	d.SetId(bridge + ":" + port)
	return resourcePortRead(d, m)
}

func resourcePortRead(d *schema.ResourceData, m interface{}) error {
	// Use Get directly for first read, or extract from ID for subsequent reads
	var port, bridge string
	if d.Id() == "" {
		var ok bool
		port, ok = d.Get("name").(string)
		if !ok {
			return fmt.Errorf("name must be a string")
		}
		bridge, ok = d.Get("bridge_id").(string)
		if !ok {
			return fmt.Errorf("bridge_id must be a string")
		}
	} else {
		// ID format is bridge:port
		parts := strings.Split(d.Id(), ":")
		if len(parts) != 2 {
			return fmt.Errorf("invalid ID format: %s", d.Id())
		}
		bridge = parts[0]
		port = parts[1]
	}

	// Check if port exists by getting the bridge ports and checking if our port is in the list
	ports, err := c.VSwitch.ListPorts(bridge)
	if err != nil {
		log.Printf("warning: error listing ports (bridge may not exist): %v", err)
		// If we can't list ports, the bridge might not exist
		d.SetId("")
		return nil
	}

	portExists := false
	for _, p := range ports {
		if p == port {
			portExists = true
			break
		}
	}

	if !portExists {
		// Port doesn't exist, remove from state
		d.SetId("")
		return nil
	}

	// Port exists, set attributes
	if err := d.Set("name", port); err != nil {
		return fmt.Errorf("error setting name: %w", err)
	}
	if err := d.Set("bridge_id", bridge); err != nil {
		return fmt.Errorf("error setting bridge_id: %w", err)
	}

	// Keep the action and ofversion attributes in the state if they're already there
	// This ensures we don't lose attributes after apply
	if action, ok := d.GetOk("action"); ok {
		if err := d.Set("action", action); err != nil {
			return fmt.Errorf("error setting action: %w", err)
		}
	}
	if ofversion, ok := d.GetOk("ofversion"); ok {
		if err := d.Set("ofversion", ofversion); err != nil {
			return fmt.Errorf("error setting ofversion: %w", err)
		}
	}

	return nil
}

func resourcePortUpdate(d *schema.ResourceData, m interface{}) error {
	port, ok := d.Get("name").(string)
	if !ok {
		return fmt.Errorf("name must be a string")
	}

	bridge, ok := d.Get("bridge_id").(string)
	if !ok {
		return fmt.Errorf("bridge_id must be a string")
	}

	action, ok := d.Get("action").(string)
	if !ok {
		return fmt.Errorf("action must be a string")
	}

	err := c.OpenFlow.ModPort(bridge, port, GetPortAction(action))
	if err != nil {
		return fmt.Errorf("error modifying port action: %w", err)
	}
	return nil
}

func resourcePortDelete(d *schema.ResourceData, m interface{}) error {
	port, ok := d.Get("name").(string)
	if !ok {
		return fmt.Errorf("name must be a string")
	}

	bridge, ok := d.Get("bridge_id").(string)
	if !ok {
		return fmt.Errorf("bridge_id must be a string")
	}

	// Deletes tap device for ovs port
	cmd := exec.Command("sudo", "/sbin/ip", "tuntap", "del", "dev", port, "mode", "tap")
	if err := cmd.Run(); err != nil {
		log.Printf("warning: error deleting tap device: %v", err)
		// Continue even if there's an error, as we still want to try to delete the port
	}

	if err := c.VSwitch.DeletePort(bridge, port); err != nil {
		return fmt.Errorf("error deleting port from bridge: %w", err)
	}

	return nil
}
