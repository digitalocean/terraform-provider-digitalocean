---
page_title: "DigitalOcean: digitalocean_reserved_ip"
---

# digitalocean_reserved_ip

Get information on a reserved IP. This data source provides the region and Droplet id
as configured on your DigitalOcean account. This is useful if the reserved IP
in question is not managed by Terraform or you need to find the Droplet the IP is
attached to.

An error is triggered if the provided reserved IP does not exist.

## Example Usage

Get the reserved IP:

```hcl
variable "public_ip" {}

data "digitalocean_reserved_ip" "example" {
  ip_address = var.public_ip
}

output "fip_output" {
  value = data.digitalocean_reserved_ip.example.droplet_id
}
```

## Argument Reference

The following arguments are supported:

* `ip_address` - (Required) The allocated IP address of the specific reserved IP to retrieve.

## Attributes Reference

The following attributes are exported:

* `region`: The region that the reserved IP is reserved to.
* `urn`: The uniform resource name of the reserved IP.
* `droplet_id`: The Droplet id that the reserved IP has been assigned to.
