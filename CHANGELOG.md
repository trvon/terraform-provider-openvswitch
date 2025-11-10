# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Comprehensive input validation for OpenFlow versions and port actions
- Unit tests for helper functions with 100% coverage
- `.golangci.yml` configuration with 20+ linters enabled
- Security scanning with `govulncheck` in CI pipeline
- Race detection in CI tests
- Coverage reporting with codecov
- Matrix testing for Terraform 1.6/1.10 and OpenTofu 1.6/1.8
- Schema descriptions for all resource fields
- `ForceNew` flags for immutable resource fields
- Lint, security, and integration test jobs to GitHub Actions

### Changed
- **BREAKING**: Updated minimum Go version from 1.18 to 1.22
- Improved error handling with proper error wrapping (`%w`)
- All `d.Set()` calls now check for errors
- All `d.Get()` type assertions now validated
- Enhanced error messages with better context
- Updated CI to use Go 1.22 across all jobs
- Streamlined README.md with clearer structure
- Improved DEVELOPMENT.md with comprehensive workflows

### Fixed
- Fixed typos: "recieve" → "receive" in port action handling
- Fixed unchecked errors in all resource CRUD operations
- Fixed error wrapping to use `%w` instead of `%s` (Go 1.13+ compatibility)
- Fixed unchecked type assertions throughout codebase
- Fixed formatting issues identified by `gofmt`

### Security
- Added `govulncheck` scanning to CI pipeline
- Enabled `gosec` linter for security issue detection
- Added race condition detection with `-race` flag

## [0.0.1] - Previous Release

### Added
- Initial provider implementation
- `openvswitch_bridge` resource
- `openvswitch_port` resource
- Basic acceptance tests
- Travis CI configuration (deprecated)

---

## Compatibility

### Terraform & OpenTofu Support

This provider is compatible with:
- **Terraform**: 1.6.0+ (tested up to 1.10.5)
- **OpenTofu**: 1.6.0+ (tested up to 1.8.10)

The same binary works for both tools with no configuration changes needed.

### Deprecations

- Travis CI support removed in favor of GitHub Actions
- Legacy Terraform SDK v0.12 (planned migration to terraform-plugin-sdk/v2)

---

## Migration Guide

### From 0.0.1 to Unreleased

No breaking changes to resource schemas or provider configuration. Existing Terraform/OpenTofu configurations will continue to work without modification.

**Note**: Building from source now requires Go 1.22 or later (previously Go 1.18).

---

## Contributing

See [CONTRIBUTING](CONTRIBUTING.md) guidelines for information on how to contribute to this project.

For detailed development instructions, see [DEVELOPMENT.md](DEVELOPMENT.md).
