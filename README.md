# Terraform OpenVSwitch Provider

[![CI](https://github.com/trvon/terraform-provider-openvswitch/actions/workflows/main.yml/badge.svg)](https://github.com/trvon/terraform-provider-openvswitch/actions/workflows/main.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/trvon/terraform-provider-openvswitch)](https://goreportcard.com/report/github.com/trvon/terraform-provider-openvswitch)

A Terraform provider for managing local Open vSwitch bridges and ports.

## Features

- Manage OVS bridges with OpenFlow protocol configuration
- Create and configure ports on bridges
- Support for tap devices and port actions
- Input validation for OpenFlow versions and port actions
- Works with Terraform 1.6+ and OpenTofu 1.6+

## Requirements

- [Go](https://golang.org/doc/install) 1.22 or later
- [Open vSwitch](https://www.openvswitch.org/) installed and running
- Root/sudo access (required for `ovs-vsctl`, `ovs-ofctl`, and `ip` commands)

## Provider Availability

This provider is not currently published to the Terraform Registry or the OpenTofu Registry. `terraform init` / `tofu init` cannot download it by `source` alone.

## Quick Start

```hcl
terraform {
  required_providers {
    openvswitch = {
      source  = "trvon/openvswitch"
      version = "0.0.1"
    }
  }
}

provider "openvswitch" {}

resource "openvswitch_bridge" "br0" {
  name      = "testbr0"
  ofversion = "OpenFlow13"  # Optional: OpenFlow10-15
}

resource "openvswitch_port" "port" {
  name      = "port0"
  bridge_id = openvswitch_bridge.br0.name
  ofversion = "OpenFlow13"
  action    = "up"  # Optional: up, down, flood, etc.
}
```

## Local Installation

Build and create a repo-local filesystem mirror:

```bash
make build
./scripts/make-provider-mirror.sh 0.0.1
```

This writes `.provider-mirror/` (ignored by git). The examples are configured to use it via `TF_CLI_CONFIG_FILE=./terraform.tfrc`.

## Resources

### `openvswitch_bridge`

Creates and manages an Open vSwitch bridge.

**Arguments:**
- `name` (Required) - Bridge name
- `ofversion` (Optional) - OpenFlow version: `OpenFlow10`, `OpenFlow11`, `OpenFlow12`, `OpenFlow13` (default), `OpenFlow14`, or `OpenFlow15`

### `openvswitch_port`

Creates and manages a port on an OVS bridge.

**Arguments:**
- `name` (Required) - Port name
- `bridge_id` (Required) - Name of the bridge to attach to
- `ofversion` (Optional) - OpenFlow version (default: `OpenFlow13`)
- `action` (Optional) - Port action: `up` (default), `down`, `stp`, `no-stp`, `receive`, `no-receive`, `no-receive-stp`, `forward`, `no-forward`, `flood`, `no-flood`, `packet-in`, or `no-packet-in`

## Installation

### From Source

```bash
git clone https://github.com/trvon/terraform-provider-openvswitch.git
cd terraform-provider-openvswitch
make build
```

Binary will be created at `bin/terraform-provider-openvswitch`.

### Development Setup

For local development and testing, see [DEVELOPMENT.md](./DEVELOPMENT.md).

## Testing

### Unit Tests

```bash
go test ./...
```

### Acceptance Tests

Requires Open vSwitch and root access:

```bash
sudo -E TF_ACC=1 go test ./openvswitch -v
```

Or using Make:

```bash
sudo -E make testacc
```

### Linting

```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linters
golangci-lint run ./...
```

## Examples

See the [examples](./examples/) directory for complete working examples:
- [sample-bridge](./examples/sample-bridge/) - Terraform example
- [opentofu-sample](./examples/opentofu-sample/) - OpenTofu example

Note: until this provider is published in the registries:

```bash
make build
./scripts/make-provider-mirror.sh 0.0.1

cd examples/sample-bridge
TF_CLI_CONFIG_FILE=./terraform.tfrc terraform init

cd ../opentofu-sample
TF_CLI_CONFIG_FILE=./terraform.tfrc tofu init
```

## Contributing

Contributions welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Ensure tests pass: `make build && go test ./...`
5. Ensure linting passes: `golangci-lint run ./...`
6. Submit a pull request

## License

Apache License 2.0 - see [LICENSE](LICENSE) for details.
