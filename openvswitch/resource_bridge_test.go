package openvswitch

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

// Skip tests if ovs-vsctl is not available
func skipIfOvsNotInstalled(t *testing.T) {
	_, err := exec.LookPath("ovs-vsctl")
	if err != nil {
		t.Skip("ovs-vsctl not found, skipping test")
	}
}

// Check if sudo is available to run the tests
func skipIfNoSudo(t *testing.T) {
	if os.Getuid() != 0 {
		_, err := exec.LookPath("sudo")
		if err != nil {
			t.Skip("sudo not found and not running as root, skipping test")
		}
	}
}

func TestAccBridge_basic(t *testing.T) {
	skipIfOvsNotInstalled(t)
	skipIfNoSudo(t)

	var bridgeName = "testbridge"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBridgeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBridgeConfig(bridgeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBridgeExists("openvswitch_bridge.test"),
					resource.TestCheckResourceAttr("openvswitch_bridge.test", "name", bridgeName),
					resource.TestCheckResourceAttr("openvswitch_bridge.test", "ofversion", "OpenFlow13"),
				),
			},
		},
	})
}

func testAccCheckBridgeDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openvswitch_bridge" {
			continue
		}

		bridgeName := rs.Primary.Attributes["name"]
		cmd := exec.Command("ovs-vsctl", "br-exists", bridgeName)
		err := cmd.Run()
		if err == nil {
			return fmt.Errorf("Bridge %s still exists", bridgeName)
		}
	}

	return nil
}

func testAccCheckBridgeExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		bridgeName := rs.Primary.Attributes["name"]
		cmd := exec.Command("ovs-vsctl", "br-exists", bridgeName)
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("Error checking bridge %s: %s", bridgeName, err)
		}

		return nil
	}
}

func testAccBridgeConfig(bridgeName string) string {
	return fmt.Sprintf(`
resource "openvswitch_bridge" "test" {
  name = "%s"
  ofversion = "OpenFlow13"
}
`, bridgeName)
}

func testAccPreCheck(t *testing.T) {
	// Could add environment checks here if needed
}

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"openvswitch": testAccProvider,
	}
}
