---
page_title: "DigitalOcean: digitalocean_droplet_autoscale"
subcategory: "Droplets"
---

# digitalocean\_droplet\_autoscale

Provides a DigitalOcean Droplet Autoscale resource. This can be used to create, modify, 
read and delete Droplet Autoscale pools.

## Example Usage

```hcl
resource "digitalocean_ssh_key" "my-ssh-key" {
  name       = "terraform-example"
  public_key = file("/Users/terraform/.ssh/id_rsa.pub")
}

resource "digitalocean_tag" "my-tag" {
  name = "terraform-example"
}

resource "digitalocean_droplet_autoscale" "my-autoscale-pool" {
  name = "terraform-example"

  config {
    min_instances             = 10
    max_instances             = 50
    target_cpu_utilization    = 0.5
    target_memory_utilization = 0.5
    cooldown_minutes          = 5
  }

  droplet_template {
    size               = "c-2"
    region             = "nyc3"
    image              = "ubuntu-24-04-x64"
    tags               = [digitalocean_tag.my-tag.id]
    ssh_keys           = [digitalocean_ssh_key.my-ssh-key.id]
    with_droplet_agent = true
    ipv6               = true
    user_data          = "\n#cloud-config\nruncmd:\n- apt-get update\n- apt-get install -y stress-ng\n"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Droplet Autoscale pool.
* `config` - (Required) The configuration parameters for Droplet Autoscale pool, the supported arguments are 
documented below.
* `droplet_template` - (Required) The droplet template parameters for Droplet Autoscale pool, the supported arguments 
are documented below.

`config` supports the following:

* `min_instances` - The minimum number of instances to maintain in the Droplet Autoscale pool.
* `max_instances` - The maximum number of instances to maintain in the Droplet Autoscale pool.
* `target_cpu_utilization` - The target average CPU load (in range `[0, 1]`) to maintain in the Droplet Autoscale pool. 
* `target_memory_utilization` - The target average Memory load (in range `[0, 1]`) to maintain in the Droplet Autoscale 
pool.
* `cooldown_minutes` - The cooldown duration between scaling events for the Droplet Autoscale pool.
* `target_number_instances` - The static number of instances to maintain in the pool Droplet Autoscale pool. This
argument cannot be used with any other config options.

`droplet_template` supports the following:

* `size` - (Required) Size slug of the Droplet Autoscale pool underlying resource(s).
* `region` - (Required) Region slug of the Droplet Autoscale pool underlying resource(s).
* `image` - (Required) Image slug of the Droplet Autoscale pool underlying resource(s).
* `tags` - List of tags to add to the Droplet Autoscale pool underlying resource(s).
* `ssh_keys` - (Required) SSH fingerprints to add to the Droplet Autoscale pool underlying resource(s).
* `vpc_uuid` - VPC UUID to create the Droplet Autoscale pool underlying resource(s). If not provided, this is inferred
from the specified `region` (default VPC).
* `with_droplet_agent` - Boolean flag to enable metric agent on the Droplet Autoscale pool underlying resource(s). The
metric agent enables collecting resource utilization metrics, which allows making resource based scaling decisions.
* `project_id` - Project UUID to create the Droplet Autoscale pool underlying resource(s).
* `ipv6` - Boolean flag to enable IPv6 networking on the Droplet Autoscale pool underlying resource(s).
* `user_data` - Custom user data that can be added to the Droplet Autoscale pool underlying resource(s). This can be a 
cloud init script that user may configure to setup their application workload.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Droplet Autoscale pool.
* `current_utilization` - The current average resource utilization of the Droplet Autoscale pool, this attribute further
embeds `memory` and `cpu` attributes to respectively report utilization data.
* `status` - Droplet Autoscale pool health status; this reflects if the pool is currently healthy and ready to accept
traffic, or in an error state and needs user intervention.
* `created_at` - Created at timestamp for the Droplet Autoscale pool.
* `updated_at` - Updated at timestamp for the Droplet Autoscale pool.

## Import

Droplet Autoscale pools can be imported using the Droplet `id`, e.g.

```
terraform import digitalocean_droplet_autoscale.my-autoscale-pool 38e66834-d741-47ec-88e7-c70cbdcz0445
```
