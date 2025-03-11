# Terraform/OpenTofu provider for Open vSwitch
[![release](https://github.com/trvon/terraform-provider-openvswitch/actions/workflows/release.yml/badge.svg)](https://github.com/trvon/terraform-provider-openvswitch/actions/workflows/release.yml)

This provider manages local Open vSwitch bridges and ports. Compatible with both Terraform and OpenTofu.

## Sample usage

From [examples/sample-bridge](./examples/sample-bridge/):

```
provider "openvswitch" {}

resource "openvswitch_bridge" "sample_bridge" {
  name = "testbr0"
  # Optional Parameters
  # OpenFlow10, OpenFlow11, OpenFlow12, OpenFlow14, OpenFlow15
  ofversion = "OpenFlow13"
}

resource "openvswitch_port" "sample_port" {
  count     = 2
  name      = "p${count.index}"
  ofversion = "OpenFlow13"
  bridge_id = openvswitch_bridge.sample_bridge.name
  # Optional Field
  action    = "up"
}
```

## Important notes
- The ip, ovs-vsctl, ovs-ofctl commands all require sudo or root access
- Error handling is currently broken

## Installation from source

Requirements:

* Go 1.18.x or later
* GNU Make
* Terraform v0.12 or later
* OpenTofu v1.6.0 or later

Clone this repo, and then do the following:

```
$ make build
```

### Development Setup

See the [Development Guide](./DEVELOPMENT.md) for detailed instructions on how to set up your development environment and use the provider locally before it's published to the Terraform Registry.

This includes:
- Using development overrides
- Setting up CLI configuration
- Linking the provider for local testing

## Running Tests

To run the unit tests:

```
$ go test ./...
```

The acceptance tests require:
- OpenVSwitch installed
- Root privileges or sudo

To run the acceptance tests:

```
$ sudo -E TF_ACC=1 go test ./openvswitch -v
```

Or use the Makefile target:

```
$ sudo -E make testacc
```

## Continuous Integration

This project uses GitHub Actions for CI/CD:

1. **Tests Workflow (`tests.yml`)**: Runs on every push and pull request to master
   - Unit tests in standard environment
   - Acceptance tests in OpenVSwitch container
   - Integration tests with both Terraform and OpenTofu

2. **Release Workflow (`release.yml`)**: Triggered manually with a version parameter
   - First runs all tests from the tests workflow
   - Creates GitHub release and publishes provider to Terraform Registry

The test workflow uses a custom container with OpenVSwitch installed to ensure tests have the necessary environment.
