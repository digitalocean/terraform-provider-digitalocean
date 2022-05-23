---
page_title: "DigitalOcean: digitalocean_floating_ip"
---

# digitalocean_floating_ip

~> **Deprecated:** DigitalOcean Floating IPs have been reserved IPs. This data source will be removed in a future release. Please use `digitalocean_reserved_ip` instead.

Get information on a floating ip. This data source provides the region and Droplet id
as configured on your DigitalOcean account. This is useful if the floating IP
in question is not managed by Terraform or you need to find the Droplet the IP is
attached to.

An error is triggered if the provided floating IP does not exist.

## Example Usage

Get the floating IP:

```hcl
variable "public_ip" {}

data "digitalocean_floating_ip" "example" {
  ip_address = var.public_ip
}

output "fip_output" {
  value = data.digitalocean_floating_ip.example.droplet_id
}
```

## Argument Reference

The following arguments are supported:

* `ip_address` - (Required) The allocated IP address of the specific floating IP to retrieve.

## Attributes Reference

The following attributes are exported:

* `region`: The region that the floating IP is reserved to.
* `urn`: The uniform resource name of the floating IP.
* `droplet_id`: The Droplet id that the floating IP has been assigned to.
