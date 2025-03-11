terraform {
  required_providers {
    openvswitch = {
      source  = "trevon/openvswitch"
      version = "0.1.0"
    }
  }
}

provider "openvswitch" {}

resource "openvswitch_bridge" "sample_bridge" {
  name = "testbr0"
  ofversion = "OpenFlow13" # OpenFlow10, OpenFlow11, OpenFlow12, OpenFlow14, OpenFlow15
}

resource "openvswitch_port" "sample_port" {
  count     = 2
  name      = "p${count.index}"
  ofversion = "OpenFlow13" # OpenFlow10, OpenFlow11, OpenFlow12, OpenFlow14, OpenFlow15
  bridge_id = openvswitch_bridge.sample_bridge.name
  # optional
  action    = "up"
}