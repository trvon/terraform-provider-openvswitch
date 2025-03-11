package openvswitch

import (
	"log"
	"os/exec"
	"os/user"

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
	err := cmd.Run()
	log.Print(err)
	err = c.VSwitch.AddPort(bridge, port)
	_ = c.OpenFlow.ModPort(bridge, port, GetPortAction(action))
	log.Print(err)
	return err
}

func resourcePortRead(d *schema.ResourceData, m interface{}) error {
	port := d.Get("name").(string)
	bridge := d.Get("bridge_id").(string)
	_, err := c.OpenFlow.DumpPort(bridge, port)
	log.Print(err)
	return err
}

func resourcePortUpdate(d *schema.ResourceData, m interface{}) error {
	port := d.Get("name").(string)
	bridge := d.Get("bridge_id").(string)
	action := d.Get("action").(string)
	err := c.OpenFlow.ModPort(bridge, port, GetPortAction(action))
	log.Print(err)
	return nil
}

func resourcePortDelete(d *schema.ResourceData, m interface{}) error {
	port := d.Get("name").(string)
	bridge := d.Get("bridge_id").(string)

	// Creates tap device for ovs port, this is not persistent
	cmd := exec.Command("sudo", "/sbin/ip", "tuntap", "del", "dev", port, "mode", "tap")
	err := cmd.Run()
	log.Print(err)
	err = c.VSwitch.DeletePort(bridge, port)
	log.Print(err)
	return err
}
