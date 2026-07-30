---
page_title: "DigitalOcean: digitalocean_microdroplet"
subcategory: "MicroDroplets"
---

# digitalocean_microdroplet

Fetches a single [DigitalOcean MicroDroplet](https://docs.digitalocean.com/products/microdroplets/) by `id` or by `name`.

## Example Usage

```hcl
data "digitalocean_microdroplet" "byid" {
  id = "506f78a4-e098-11e5-ad9f-000f53306ae1"
}

data "digitalocean_microdroplet" "byname" {
  name = "my-microdroplet"
}
```

## Argument Reference

Exactly one of the following must be provided:

* `id` - (Optional) The MicroDroplet UUID.
* `name` - (Optional) The MicroDroplet name. Errors if zero or more than one match.

## Attributes Reference

The following attributes are exported. See the [resource docs](../resources/microdroplet.md) for descriptions:

* `region`, `size`, `image`, `networking`, `vpc_uuid`, `http_port`, `http_protocol`, `environment`, `auto_pause`, `auto_resume`, `endpoint`, `current_state`, `created_at`, `urn`.

~> **Note:** `tags` are intentionally not exported. The MicroDroplets API does not return the tags attached at create time, so the data source cannot round-trip them. See the [resource docs](../resources/microdroplet.md) for details.
