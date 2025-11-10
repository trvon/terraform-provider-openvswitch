# Development Guide

This guide covers local development and testing of the OpenVSwitch Terraform provider.

## Prerequisites

- Go 1.22 or later
- Open vSwitch installed (`ovs-vsctl --version`)
- Sudo/root access
- Terraform 1.6+ or OpenTofu 1.6+

## Quick Setup

```bash
# Clone the repository
git clone https://github.com/trvon/terraform-provider-openvswitch.git
cd terraform-provider-openvswitch

# Install dependencies and build
make build

# Run tests
go test ./...
```

## Local Development Workflow

### 1. Development Overrides (Recommended)

Create a `~/.terraformrc` (or `~/.tofurc` for OpenTofu):

```hcl
provider_installation {
  dev_overrides {
    "registry.terraform.io/trvon/openvswitch" = "/path/to/terraform-provider-openvswitch/bin"
  }
  direct {}
}
```

Then build and test:

```bash
make build
cd examples/sample-bridge
terraform plan  # Uses your local binary
```

### 2. Manual Installation

```bash
# Build
make build

# Install to local plugins directory
mkdir -p ~/.terraform.d/plugins/local/trvon/openvswitch/1.0.0/darwin_arm64
cp bin/terraform-provider-openvswitch ~/.terraform.d/plugins/local/trvon/openvswitch/1.0.0/darwin_arm64/

# Use in Terraform config
terraform {
  required_providers {
    openvswitch = {
      source = "local/trvon/openvswitch"
      version = "1.0.0"
    }
  }
}
```

Adjust the architecture folder (`darwin_arm64`, `linux_amd64`, etc.) for your system.

## Testing

### Unit Tests

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific test
go test -v -run TestGetPortAction ./openvswitch
```

### Acceptance Tests

Requires Open vSwitch and root access:

```bash
# Run acceptance tests
sudo -E TF_ACC=1 go test ./openvswitch -v

# Run specific acceptance test
sudo -E TF_ACC=1 go test ./openvswitch -v -run TestAccBridge_basic
```

### Race Detection

```bash
go test -race ./...
```

### Coverage

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Linting

### Install Tools

```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Install govulncheck
go install golang.org/x/vuln/cmd/govulncheck@latest
```

### Run Linters

```bash
# Run all configured linters
golangci-lint run ./...

# Run go vet
go vet ./...

# Check formatting
gofmt -s -l .

# Security scan
govulncheck ./...
```

### Fix Issues

```bash
# Auto-fix formatting
gofmt -s -w .

# Auto-fix some lint issues
golangci-lint run --fix ./...
```

## Build Targets

```bash
make build        # Build provider binary
make test         # Run unit tests
make testacc      # Run acceptance tests (requires sudo)
make fmt          # Format code
make fmtcheck     # Check formatting
make vet          # Run go vet
make lint         # Run linters (requires golangci-lint)
```

## Debugging

### Enable Terraform Debug Logging

```bash
export TF_LOG=DEBUG
export TF_LOG_PATH=./terraform.log
terraform apply
```

### Provider Debug Logging

The provider uses Go's standard `log` package. Logs are output to stderr.

### Debugging with Delve

```bash
# Build with debug symbols
go build -gcflags="all=-N -l" -o bin/terraform-provider-openvswitch

# Run with delve
dlv exec bin/terraform-provider-openvswitch
```

## Code Organization

```
.
├── main.go                          # Provider entry point
├── openvswitch/                     # Provider implementation
│   ├── provider.go                  # Provider definition
│   ├── resource_bridge.go           # Bridge resource
│   ├── resource_bridge_test.go      # Bridge tests
│   ├── resource_port.go             # Port resource
│   ├── resource_port_test.go        # Port tests
│   └── resource_port_helpers_test.go # Unit tests
├── examples/                        # Usage examples
├── .golangci.yml                    # Linter configuration
└── .github/workflows/main.yml       # CI/CD pipeline
```

## Making Changes

### Before Submitting a PR

1. **Write Tests**: Add unit or acceptance tests for new features
2. **Run Tests**: Ensure all tests pass
   ```bash
   make build
   go test ./...
   sudo -E TF_ACC=1 go test ./openvswitch -v
   ```
3. **Run Linters**: Fix all linting issues
   ```bash
   golangci-lint run ./...
   gofmt -s -w .
   ```
4. **Update Documentation**: Update README.md if needed

### Code Style

- Follow standard Go conventions
- Run `gofmt -s` on all files
- Use meaningful variable and function names
- Add comments for exported functions
- Check errors explicitly (no `_` for errors)
- Use error wrapping with `fmt.Errorf(..., %w, err)`

## CI/CD Pipeline

GitHub Actions runs automatically on push and PR:

1. **Lint Job**: golangci-lint, go vet, gofmt check
2. **Security Job**: govulncheck, race detector
3. **Unit Tests**: Standard Go tests with coverage
4. **Acceptance Tests**: In OVS container with sudo
5. **Integration Tests**: Matrix with Terraform/OpenTofu versions

All jobs must pass before merging.

## Common Issues

### "ovs-vsctl: command not found"

Install Open vSwitch:

```bash
# macOS
brew install openvswitch
brew services start openvswitch

# Ubuntu/Debian
sudo apt-get install openvswitch-switch

# RHEL/CentOS
sudo yum install openvswitch
```

### Permission Denied

OVS commands require root:

```bash
# Add your user to sudoers for passwordless sudo (development only)
echo "$USER ALL=(ALL) NOPASSWD: ALL" | sudo tee /etc/sudoers.d/$USER
```

### Tests Skipped

If you see "ovs-vsctl not found, skipping test", install OVS and ensure it's in PATH.

## Resources

- [Terraform Plugin Development](https://developer.hashicorp.com/terraform/plugin)
- [terraform-plugin-sdk](https://github.com/hashicorp/terraform-plugin-sdk)
- [Open vSwitch Documentation](https://docs.openvswitch.org/)
- [Go Testing](https://golang.org/pkg/testing/)
