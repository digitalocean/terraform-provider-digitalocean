---
page_title: "DigitalOcean: digitalocean_loadbalancer"
---

# digitalocean_loadbalancer

Get information on a load balancer for use in other resources. This data source
provides all of the load balancers properties as configured on your DigitalOcean
account. This is useful if the load balancer in question is not managed by
Terraform or you need to utilize any of the load balancers data.

An error is triggered if the provided load balancer name does not exist.

## Example Usage

Get the load balancer by name:

```hcl
data "digitalocean_loadbalancer" "example" {
  name = "app"
}

output "lb_output" {
  value = data.digitalocean_loadbalancer.example.ip
}
```

Get the load balancer by ID:

```hcl
data "digitalocean_loadbalancer" "example" {
  id = "loadbalancer_id"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of load balancer.
* `id` - (Optional) The ID of load balancer.
* `urn` - The uniform resource name for the Load Balancer

## Attributes Reference

See the [Load Balancer Resource](/providers/digitalocean/digitalocean/latest/docs/resources/loadbalancer) for details on the
returned attributes - they are identical.
