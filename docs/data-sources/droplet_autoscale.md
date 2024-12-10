---
page_title: "DigitalOcean: digitalocean_droplet_autoscale"
subcategory: "Droplets"
---

# digitalocean\_droplet\_autoscale

Get information on a Droplet Autoscale pool for use with other managed resources. This datasource provides all the
Droplet Autoscale pool properties as configured on the DigitalOcean account. This is useful if the Droplet Autoscale 
pool in question is not managed by Terraform, or any of the relevant data would need to referenced in other managed 
resources.

## Example Usage

Get the Droplet Autoscale pool by name:

```hcl
data "digitalocean_droplet_autoscale" "my-imported-autoscale-pool" {
  name = digitalocean_droplet_autoscale.my-existing-autoscale-pool.name
}
```

Get the Droplet Autoscale pool by ID:

```hcl
data "digitalocean_droplet_autoscale" "my-imported-autoscale-pool" {
  id = digitalocean_droplet_autoscale.my-existing-autoscale-pool.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of Droplet Autoscale pool.
* `id` - (Optional) The ID of Droplet Autoscale pool.

## Attributes Reference

See the [Droplet Autoscale Resource](../resources/droplet_autoscale.md) for details on the
returned attributes - they are identical.
