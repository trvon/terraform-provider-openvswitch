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
				Type:     schema.TypeString,
				Required: true,
			},

			"bridge_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"action": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "up",
			},
			"ofversion": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "OpenFlow13",
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
	case ("recieve"):
		return ovs.PortActionReceive
	case ("no-recieve"):
		return ovs.PortActionNoReceive
	case ("no-recieve-stp"):
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
	port := d.Get("name").(string)
	bridge := d.Get("bridge_id").(string)
	action := d.Get("action").(string)

	// Creates tap device for ovs port, this is not persistent
	user, _ := user.Current()
	cmd := exec.Command("sudo", "/sbin/ip", "tuntap", "add", "dev", port, "mode", "tap", "user", user.Username)
	if err := cmd.Run(); err != nil {
		log.Print(err)
		// Continue even if there's an error, as the tap device might already exist
	}

	if err := c.VSwitch.AddPort(bridge, port); err != nil {
		return err
	}

	if err := c.OpenFlow.ModPort(bridge, port, GetPortAction(action)); err != nil {
		log.Print(err)
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
		port = d.Get("name").(string)
		bridge = d.Get("bridge_id").(string)
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
		log.Print(err)
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
	d.Set("name", port)
	d.Set("bridge_id", bridge)

	// Keep the action and ofversion attributes in the state if they're already there
	// This ensures we don't lose attributes after apply
	if action, ok := d.GetOk("action"); ok {
		d.Set("action", action)
	}
	if ofversion, ok := d.GetOk("ofversion"); ok {
		d.Set("ofversion", ofversion)
	}

	return nil
}

func resourcePortUpdate(d *schema.ResourceData, m interface{}) error {
	port := d.Get("name").(string)
	bridge := d.Get("bridge_id").(string)
	action := d.Get("action").(string)
	err := c.OpenFlow.ModPort(bridge, port, GetPortAction(action))
	if err != nil {
		log.Print(err)
	}
	return err
}

func resourcePortDelete(d *schema.ResourceData, m interface{}) error {
	port := d.Get("name").(string)
	bridge := d.Get("bridge_id").(string)

	// Deletes tap device for ovs port
	cmd := exec.Command("sudo", "/sbin/ip", "tuntap", "del", "dev", port, "mode", "tap")
	if err := cmd.Run(); err != nil {
		log.Print(err)
		// Continue even if there's an error, as we still want to try to delete the port
	}

	if err := c.VSwitch.DeletePort(bridge, port); err != nil {
		log.Print(err)
		return err
	}

	return nil
}
