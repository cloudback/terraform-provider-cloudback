terraform {
  required_providers {
    cloudback = {
      source = "cloudback/cloudback"
    }
  }
}

provider "cloudback" {
  endpoint = "https://app.cloudback.it"
  api_key  = "your-api-key"
}

