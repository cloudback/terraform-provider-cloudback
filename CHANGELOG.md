# Cloudback Terraform Provider Changelog

## 1.0.4 (2025-11-22)

- Update backup definition resource: Add AzureDevOps to platform examples and adjust schema for subject_type and subject_name fields.

## 1.0.3 (2025-07-15)

- Security updates for dependencies

## 1.0.2 (2025-05-11)

- Bump golang.org/x/crypto to v0.35.0 in /tools
- Bump golang.org/x/net to v0.38.0

## 1.0.1 (2025-03-09)

- Integration tests were fixed

## 1.0.0 (2025-03-09)

### Features

- Introduced an initial stable release of the Cloudback Terraform Provider.

### New Resource

- **`cloudback_backup_definition`**: Defines and manages backup configurations for Cloudback, including schedule, retention, and other backup parameters.

### Additional Notes

- If you encounter any issues or have suggestions for improvements, please [open an issue](https://github.com/cloudback/terraform-provider-cloudback/issues/new) on our repository.
