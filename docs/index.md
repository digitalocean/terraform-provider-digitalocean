---
page_title: "Provider: DigitalOcean"
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
  token = var.do_token
}

# Create a web server
resource "digitalocean_droplet" "web" {
  # ...
}
```

## Module Example

The provider is also not maintained by hashicorp, so you will need to explicitly define the required provider in every module so that terraform knows where to find the provider and version. Provider Inheritance does not work with Digital Ocean like it does for hashicorp maintained providers like AWS, Azure, and GCP.

Albeit against DRY programming principles, you will need to repeat this block at the top of _each of your modules_ because the inheritance from the parent terraform file that calls the module will not trickle down properly.

```hcl
#Set the required providers block to define the proper routing to digitalocean.
terraform {
  required_providers {
    digitalocean = {
      source = "digitalocean/digitalocean"
      version = ">= 2.8.0"
    }
  }
}
#Just like in the example above, pass your `DIGITALOCEAN_TOKEN` to the module
provider "digitalocean" {
  token = var.do_token
}
```

For Example, if you had a tree of modules, you would need to include this provider code in the main.tf of each module. `Dedicated-Server`, `Network`, and `Security-Groups`
```
├── main.tf
├── modules
│   ├── dedicated-server
│   │   ├── main.tf
│   │   ├── output.tf
│   │   └── vars.tf
│   ├── network
│   │   ├── main.tf
│   │   ├── output.tf
│   │   └── vars.tf
│   ├── security-groups
│   │   ├── main.tf
│   │   ├── output.tf
│   │   └── vars.tf
```

For the Network Module, we call it in main.tf like so:

```hcl
# BUILD NETWORK
module "network" {
    source  = "./modules/network"
    vpcname = "examplevpc"
    ip_range = "10.10.10.0/24"
    region  = var.region
    token   = var.token
}
```

Passing the parameters to the module, but making sure we explicitly define the provider again so the module knows where to look for the resource since it defaults to looking into hashicorp maintained providers.

```hcl
#NETWORK MODULE: main.tf
terraform {
  required_providers {
    digitalocean = {
      source = "digitalocean/digitalocean"
      version = ">= 2.8.0"
    }
  }
}
provider "digitalocean" {
  token = var.token
}

#RANGE 10.10.10.x VPC NETWORK
resource "digitalocean_vpc" "examplevpc" {
  name     = var.vpcname
  region   = var.region
  ip_range = var.ip_range
}

```

## Argument Reference

The following arguments are supported:

* `token` - (Required) This is the DO API token. Alternatively, this can also be specified
  using environment variables ordered by precedence:
  * `DIGITALOCEAN_TOKEN`
  * `DIGITALOCEAN_ACCESS_TOKEN`
* `spaces_access_id` - (Optional) The access key ID used for Spaces API
  operations (Defaults to the value of the `SPACES_ACCESS_KEY_ID` environment
  variable).
* `spaces_secret_key` - (Optional) The secret access key used for Spaces API
  operations (Defaults to the value of the `SPACES_SECRET_ACCESS_KEY`
  environment variable).
* `api_endpoint` - (Optional) This can be used to override the base URL for
  DigitalOcean API requests (Defaults to the value of the `DIGITALOCEAN_API_URL`
  environment variable or `https://api.digitalocean.com` if unset).
* `spaces_endpoint` - (Optional) This can be used to override the endpoint URL
  used for DigitalOcean Spaces requests. (It defaults to the value of the
  `SPACES_ENDPOINT_URL` environment variable or `https://{{.Region}}.digitaloceanspaces.com`
  if unset.) The provider will replace `{{.Region}}` (via Go's templating engine) with the slug
  of the applicable Spaces region. 
