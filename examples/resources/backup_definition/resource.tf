resource "cloudback_backup_definition" "example" {
  platform   = "GitHub"                 # Currently only GitHub is supported
  account    = "your-github-account"    # The GitHub account that owns the repository
  repository = "your-github-repository" # The repository to backup
  settings = {
    enabled   = true             # Enable the scheduled automated backup
    schedule  = "Daily at 6 am"  # The schedule for the automated backup, see the Cloudback Dashboard for available options
    storage   = "Your S3 bucket" # The storage name to use for the backup, see the Cloudback Dashboard for available options
    retention = "Last 30 days"   # The retention policy for the backup, see the Cloudback Dashboard for available options
  }
}