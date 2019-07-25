---
layout: "digitalocean"
page_title: "DigitalOcean: digitalocean_loadbalancer"
sidebar_current: "docs-do-datasource-loadbalancer"
description: |-
  Get information on a loadbalancer.
---

# digitalocean_loadbalancer

Get information on a load balancer for use in other resources. This data source
provides all of the load balancers properties as configured on your DigitalOcean
account. This is useful if the load balancer in question is not managed by
Terraform or you need to utilize any of the load balancers data.

An error is triggered if the provided load balancer name does not exist.

## Example Usage

Get the load balancer:

```hcl
data "digitalocean_loadbalancer" "example" {
  name = "app"
}

output "lb_output" {
  value = data.digitalocean_loadbalancer.example.ip
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of load balancer.
* `urn` - The uniform resource name for the Load Balancer

## Attributes Reference

See the [Load Balancer Resource](/docs/providers/do/r/loadbalancer.html) for details on the
returned attributes - they are identical.