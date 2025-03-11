package openvswitch

import (
	"fmt"
	"os/exec"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccPort_basic(t *testing.T) {
	skipIfOvsNotInstalled(t)
	skipIfNoSudo(t)

	var bridgeName = "testbridge"
	var portName = "testport"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPortConfig(bridgeName, portName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPortExists("openvswitch_port.test"),
					resource.TestCheckResourceAttr("openvswitch_port.test", "name", portName),
					resource.TestCheckResourceAttr("openvswitch_port.test", "ofversion", "OpenFlow13"),
					resource.TestCheckResourceAttr("openvswitch_port.test", "action", "up"),
				),
			},
		},
	})
}

func testAccCheckPortDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openvswitch_port" {
			continue
		}

		portName := rs.Primary.Attributes["name"]
		bridgeName := rs.Primary.Attributes["bridge_id"]

		// Check if port still exists on the bridge
		cmd := exec.Command("ovs-vsctl", "list-ports", bridgeName)
		output, err := cmd.Output()
		if err == nil && string(output) != "" {
			return fmt.Errorf("Port %s still exists on bridge %s", portName, bridgeName)
		}

		// Check if tap device still exists
		cmd = exec.Command("ip", "link", "show", portName)
		err = cmd.Run()
		if err == nil {
			return fmt.Errorf("Tap device %s still exists", portName)
		}
	}

	return nil
}

func testAccCheckPortExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		portName := rs.Primary.Attributes["name"]
		bridgeName := rs.Primary.Attributes["bridge_id"]

		// Check if port exists on the bridge
		cmd := exec.Command("ovs-vsctl", "port-to-br", portName)
		output, err := cmd.Output()
		if err != nil || string(output) != bridgeName+"\n" {
			return fmt.Errorf("Port %s doesn't exist on bridge %s", portName, bridgeName)
		}

		return nil
	}
}

func testAccPortConfig(bridgeName, portName string) string {
	return fmt.Sprintf(`
resource "openvswitch_bridge" "test" {
  name = "%s"
  ofversion = "OpenFlow13"
}

resource "openvswitch_port" "test" {
  name = "%s"
  bridge_id = openvswitch_bridge.test.name
  ofversion = "OpenFlow13"
  action = "up"
}
`, bridgeName, portName)
}
