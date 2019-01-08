---
layout: "digitalocean"
page_title: "Provider: DigitalOcean"
sidebar_current: "docs-do-index"
description: |-
  The DigitalOcean (DO) provider is used to interact with the resources supported by DigitalOcean. The provider needs to be configured with the proper credentials before it can be used.
---

# DigitalOcean Provider

The DigitalOcean (DO) provider is used to interact with the
resources supported by DigitalOcean. The provider needs to be configured
with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
# Set the variable value in *.tfvars file
# or using -var="do_token=..." CLI option
variable "do_token" {}

# Configure the DigitalOcean Provider
provider "digitalocean" {
  token = "${var.do_token}"
}

# Create a web server
resource "digitalocean_droplet" "web" {
  # ...
}
```

## Argument Reference

The following arguments are supported:

* `token` - (Required) This is the DO API token. This can also be specified
  with the `DIGITALOCEAN_TOKEN` shell environment variable.
* `spaces_access_id` - (Optional) The access key ID used for Spaces API
  operations (Defaults to the value of the `SPACES_ACCESS_KEY_ID` environment
  variable).
* `spaces_secret_key` - (Optional) The secret access key used for Spaces API
  operations (Defaults to the value of the `SPACES_SECRET_ACCESS_KEY`
  environment variable).
