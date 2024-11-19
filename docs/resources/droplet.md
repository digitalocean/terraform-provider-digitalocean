---
page_title: "DigitalOcean: digitalocean_droplet"
subcategory: "Droplets"
---

# digitalocean\_droplet

Provides a DigitalOcean Droplet resource. This can be used to create,
modify, and delete Droplets. Droplets also support
[provisioning](https://www.terraform.io/docs/language/resources/provisioners/syntax.html).

## Example Usage

```hcl
# Create a new Web Droplet in the nyc2 region
resource "digitalocean_droplet" "web" {
  image   = "ubuntu-20-04-x64"
  name    = "web-1"
  region  = "nyc2"
  size    = "s-1vcpu-1gb"
  backups = true
  backup_policy {
    plan    = "weekly"
    weekday = "TUE"
    hour    = 8
  }
}
```

## Argument Reference

The following arguments are supported:

* `image` - (Required) The Droplet image ID or slug. This could be either image ID or droplet snapshot ID.
* `name` - (Required) The Droplet name.
* `region` - The region where the Droplet will be created.
* `size` - (Required) The unique slug that identifies the type of Droplet. You can find a list of available slugs on [DigitalOcean API documentation](https://docs.digitalocean.com/reference/api/api-reference/#tag/Sizes).
* `backups` - (Optional) Boolean controlling if backups are made. Defaults to
   false.
* `backup_policy` - (Optional) An object specifying the backup policy for the Droplet. If omitted and `backups` is `true`, the backup plan will default to daily.
  - `plan` - The backup plan used for the Droplet. The plan can be either `daily` or `weekly`.
  - `weekday` - The day of the week on which the backup will occur (`SUN`, `MON`, `TUE`, `WED`, `THU`, `FRI`, `SAT`).
  - `hour` - The hour of the day that the backup window will start (`0`, `4`, `8`, `12`, `16`, `20`).
* `monitoring` - (Optional) Boolean controlling whether monitoring agent is installed.
   Defaults to false. If set to `true`, you can configure monitor alert policies
   [monitor alert resource](/providers/digitalocean/digitalocean/latest/docs/resources/monitor_alert)
* `ipv6` - (Optional) Boolean controlling if IPv6 is enabled. Defaults to false.
   Once enabled for a Droplet, IPv6 can not be disabled. When enabling IPv6 on
   an existing Droplet, [additional OS-level configuration](https://docs.digitalocean.com/products/networking/ipv6/how-to/enable/#on-existing-droplets)
   is required.
* `vpc_uuid` - (Optional) The ID of the VPC where the Droplet will be located.
* `private_networking` - (Optional) **Deprecated** Boolean controlling if private networking
  is enabled. This parameter has been deprecated. Use `vpc_uuid` instead to specify a VPC network for the Droplet. If no `vpc_uuid` is provided, the Droplet will be placed in your account's default VPC for the region.
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
* `user_data` (Optional) - A string of the desired User Data provided [during Droplet creation](https://docs.digitalocean.com/products/droplets/how-to/provide-user-data/). Changing this forces a new resource to be created.
* `volume_ids` (Optional) - A list of the IDs of each [block storage volume](/providers/digitalocean/digitalocean/latest/docs/resources/volume) to be attached to the Droplet.
* `droplet_agent` (Optional) - A boolean indicating whether to install the
   DigitalOcean agent used for providing access to the Droplet web console in
   the control panel. By default, the agent is installed on new Droplets but
   installation errors (i.e. OS not supported) are ignored. To prevent it from
   being installed, set to `false`. To make installation errors fatal, explicitly
   set it to `true`.
* `graceful_shutdown` (Optional) - A boolean indicating whether the droplet
   should be gracefully shut down before it is deleted.

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
