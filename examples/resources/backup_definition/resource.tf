# Example using the new subject_type and subject_name fields (recommended)
resource "cloudback_backup_definition" "example" {
  platform     = "GitHub"                 # Platform name (GitHub, AzureDevOps, etc.)
  account      = "your-github-account"    # The account that owns the subject
  subject_type = "Repository"             # Subject type (Repository, Project, etc.)
  subject_name = "your-github-repository" # The name of the subject to backup
  settings = {
    enabled   = true             # Enable the scheduled automated backup
    schedule  = "Daily at 6 am"  # The schedule for the automated backup, see the Cloudback Dashboard for available options
    storage   = "Your S3 bucket" # The storage name to use for the backup, see the Cloudback Dashboard for available options
    retention = "Last 30 days"   # The retention policy for the backup, see the Cloudback Dashboard for available options
  }
}

# Legacy example using repository field (still supported for backward compatibility)
resource "cloudback_backup_definition" "example_legacy" {
  platform   = "GitHub"                 # Currently only GitHub is supported
  account    = "your-github-account"    # The GitHub account that owns the repository
  repository = "your-github-repository" # The repository to backup (deprecated: use subject_type and subject_name)
  settings = {
    enabled   = true             # Enable the scheduled automated backup
    schedule  = "Daily at 6 am"  # The schedule for the automated backup, see the Cloudback Dashboard for available options
    storage   = "Your S3 bucket" # The storage name to use for the backup, see the Cloudback Dashboard for available options
    retention = "Last 30 days"   # The retention policy for the backup, see the Cloudback Dashboard for available options
  }
}