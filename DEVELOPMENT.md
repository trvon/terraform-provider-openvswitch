# Development Guide

This guide explains how to set up your development environment for the OpenVSwitch Terraform provider.

## Building the Provider

1. Clone the repository
2. Build the provider using `make build`
```
make build
```

This will generate the provider binary in the `bin/` directory.

## Testing the Provider Locally

### Method 1: Using Development Overrides

Create a `.terraformrc` file (or `.tofurc` for OpenTofu) in your home directory with the following content:

```hcl
provider_installation {
  dev_overrides {
    "trevon/openvswitch" = "/absolute/path/to/your/terraform-provider-openvswitch/bin"
  }
  direct {}
}
```

Replace `/absolute/path/to/your` with your actual project path.

### Method 2: Using the Terraform CLI Configuration File

For specific Terraform configurations, you can create a `terraform.tfrc` file in the same directory as your configuration:

```hcl
provider_installation {
  dev_overrides {
    "trevon/openvswitch" = "/absolute/path/to/your/terraform-provider-openvswitch/bin"
  }
  direct {}
}
```

Then export the path to this file:

```bash
export TF_CLI_CONFIG_FILE="$(pwd)/terraform.tfrc"
```

### Method 3: Using the Terraform Provider Registry Protocol (TPRP)

For advanced development cycles, you may want to use TPRP which supports features like automatic binary reloading.

1. Install Go, if not already installed.
2. Set up your project with provider caching:

```bash
mkdir -p ~/.terraform.d/plugins/trevon.dev/trevon/openvswitch/0.1.0/linux_amd64
ln -s $(pwd)/bin/terraform-provider-openvswitch ~/.terraform.d/plugins/trevon.dev/trevon/openvswitch/0.1.0/linux_amd64/
```

## Running Tests

```bash
# Run unit tests
make test

# Run acceptance tests
make testacc
```

Note: Acceptance tests require an actual OpenVSwitch instance and will make real changes.