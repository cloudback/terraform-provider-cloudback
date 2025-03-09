# Cloudback Terraform Provider

- Website: https://registry.terraform.io/providers/cloudback/cloudback

## Maintainers

This provider plugin is maintained by the Cloudback team at [Cloudback](https://cloudback.it/).

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0.x
- [Go](https://golang.org/doc/install) >= 1.22 (to build the provider plugin)

## Usage

```hcl
terraform {
  required_providers {
    cloudback = {
      source = "cloudback/cloudback"
    }
  }
}

provider "cloudback" {
  endpoint = "https://app.cloudback.it"
  api_key = "your-api-key"
}

resource "cloudback_backup_definition" "example" {
  platform = "GitHub"                   # Currently only GitHub is supported
  account = "your-github-account"       # The GitHub account that owns the repository
  repository = "your-github-repository" # The repository to backup
  settings = {
    enabled = true              # Enable the scheduled automated backup
    schedule = "Daily at 6 am"  # The schedule for the automated backup, see the Cloudback Dashboard for available options
    storage = "Your S3 bucket"  # The storage name to use for the backup, see the Cloudback Dashboard for available options
    retention = "Last 30 days"  # The retention policy for the backup, see the Cloudback Dashboard for available options
  }
}
```

## Building The Provider

Clone repository to: `$GOPATH/src/github.com/cloudback/terraform-provider-cloudback`

```sh
$ mkdir -p $GOPATH/src/github.com/cloudback; cd $GOPATH/src/github.com/cloudback
$ git clone git@github.com:cloudback/terraform-provider-cloudback
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/cloudback/terraform-provider-cloudback
$ make build
```

## Using the provider

To use a released provider in your Terraform environment, run
[`terraform init`](https://www.terraform.io/docs/commands/init.html) and
Terraform will automatically install the provider. To specify a particular
provider version when installing released providers, see the [Terraform
documentation on provider versioning](https://www.terraform.io/docs/configuration/providers.html#version-provider-versions).

Refer to the section below for instructions on how to to use a custom-built
version of the provider.

For either installation method, documentation about the provider specific
configuration options can be found on the
[provider's website](https://www.terraform.io/docs/providers/cloudback/).

## Developing the Provider

If you wish to work on the provider, you'll first need
[Go](http://www.golang.org) installed on your machine (version 1.20+ is
*required*). You'll also need to correctly setup a
[GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding
`$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put
the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-cloudback
...
```

To use the provider binary in a local configuration, create a file called
`.terraformrc` in your home directory and specify a [development
override][tf_docs_dev_overrides] for the `cloudback` provider.

```hcl
provider_installation {
  dev_overrides {
    "cloudback/cloudback" = "<ABSOLUTE PATH TO YOUR GOPATH>/bin/"
  }
}
```

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```
