---
page_title: "DigitalOcean: digitalocean_droplet"
---

# digitalocean\_droplet

Provides a DigitalOcean Droplet resource. This can be used to create,
modify, and delete Droplets. Droplets also support
[provisioning](https://www.terraform.io/docs/language/resources/provisioners/syntax.html).

## Example Usage

```hcl
# Create a new Web Droplet in the nyc2 region
resource "digitalocean_droplet" "web" {
  image  = "ubuntu-18-04-x64"
  name   = "web-1"
  region = "nyc2"
  size   = "s-1vcpu-1gb"
}
```

## Argument Reference

The following arguments are supported:

* `image` - (Required) The Droplet image ID or slug.
* `name` - (Required) The Droplet name.
* `region` - (Required) The region to start in.
* `size` - (Required) The unique slug that indentifies the type of Droplet. You can find a list of available slugs on [DigitalOcean API documentation](https://docs.digitalocean.com/reference/api/api-reference/#tag/Sizes).
* `backups` - (Optional) Boolean controlling if backups are made. Defaults to
   false.
* `monitoring` - (Optional) Boolean controlling whether monitoring agent is installed.
   Defaults to false.
* `ipv6` - (Optional) Boolean controlling if IPv6 is enabled. Defaults to false.
* `vpc_uuid` - (Optional) The ID of the VPC where the Droplet will be located.
* `private_networking` - (Optional)  Boolean controlling if private networking
  is enabled. When VPC is enabled on an account, this will provision the
  Droplet inside of your account's default VPC for the region. Use the
  `vpc_uuid` attribute to specify a different VPC.
* `ssh_keys` - (Optional) A list of SSH key IDs or fingerprints to enable in
   the format `[12345, 123456]`. To retrieve this info, use the
   [DigitalOcean API](https://docs.digitalocean.com/reference/api/api-reference/#tag/SSH-Keys)
   or CLI (`doctl compute ssh-key list`). Once a Droplet is created keys can not
   be added or removed via this provider. Modifying this field will prompt you
   to destroy and recreate the Droplet.
* `resize_disk` - (Optional) Boolean controlling whether to increase the disk
   size when resizing a Droplet. It defaults to `true`. When set to `false`,
   only the Droplet's RAM and CPU will be resized. **Increasing a Droplet's disk
   size is a permanent change**. Increasing only RAM and CPU is reversible.
* `tags` - (Optional) A list of the tags to be applied to this Droplet.
* `user_data` (Optional) - A string of the desired User Data for the Droplet.
* `volume_ids` (Optional) - A list of the IDs of each [block storage volume](/providers/digitalocean/digitalocean/latest/docs/resources/volume) to be attached to the Droplet.

~> **NOTE:** If you use `volume_ids` on a Droplet, Terraform will assume management over the full set volumes for the instance, and treat additional volumes as a drift. For this reason, `volume_ids` must not be mixed with external `digitalocean_volume_attachment` resources for a given instance.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Droplet
* `urn` - The uniform resource name of the Droplet
* `name`- The name of the Droplet
* `region` - The region of the Droplet
* `image` - The image of the Droplet
* `ipv6` - Is IPv6 enabled
* `ipv6_address` - The IPv6 address
* `ipv4_address` - The IPv4 address
* `ipv4_address_private` - The private networking IPv4 address
* `locked` - Is the Droplet locked
* `private_networking` - Is private networking enabled
* `price_hourly` - Droplet hourly price
* `price_monthly` - Droplet monthly price
* `size` - The instance size
* `disk` - The size of the instance's disk in GB
* `vcpus` - The number of the instance's virtual CPUs
* `status` - The status of the Droplet
* `tags` - The tags associated with the Droplet
* `volume_ids` - A list of the attached block storage volumes

## Import

Droplets can be imported using the Droplet `id`, e.g.

```
terraform import digitalocean_droplet.mydroplet 100823
```
