---
page_title: "DigitalOcean: digitalocean_microvm_ingress"
subcategory: "MicroVMs"
---

# digitalocean_microvm_ingress

Provides an additional ingress URL for a [`digitalocean_microvm`](microvm).
Ingress URLs are minted by the platform, so the configuration cannot declare a
hostname; it only chooses the guest `port`. Adding this resource adds a URL to
the MicroVM in place, and removing it removes the URL. Neither operation
replaces the MicroVM.

Use the `ingress_enabled` argument on `digitalocean_microvm` to control the
*default* URL. Use this resource for any URL beyond the default.

-> **Note:** MicroVMs are in public preview.

## Example Usage

```hcl
resource "digitalocean_microvm" "web" {
  name      = "example-microvm"
  region    = "nyc3"
  http_port = 8080
  ports     = [8080, 9090]

  size {
    cpu    = 2
    memory = 4096
  }

  source {
    oci_ref = "docker.io/library/nginx:latest"
  }
}

# An extra URL that routes to port 9090.
resource "digitalocean_microvm_ingress" "admin" {
  microvm_id = digitalocean_microvm.web.id
  port       = 9090
}

output "admin_url" {
  value = digitalocean_microvm_ingress.admin.hostname
}
```

### MicroVM without a default URL

```hcl
resource "digitalocean_microvm" "api" {
  name            = "example-api"
  region          = "nyc3"
  ingress_enabled = false

  size {
    cpu    = 2
    memory = 4096
  }

  source {
    oci_ref = "docker.io/library/nginx:latest"
  }
}

resource "digitalocean_microvm_ingress" "api" {
  microvm_id = digitalocean_microvm.api.id
}
```

## Argument Reference

The following arguments are supported:

* `microvm_id` - (Required) The ID of the MicroVM to add the URL to. Changing this forces a new URL.
* `port` - (Optional) The guest port the URL routes to. Must be one of the MicroVM's `ports`. Defaults to the MicroVM's `http_port` when omitted. Changing this forces a new URL.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:

* `id` - The stable identifier of the URL.
* `url_id` - The stable identifier of the URL (same as `id`).
* `hostname` - The assigned hostname, without a scheme.
* `default` - Whether this is the MicroVM's system default URL.
* `status` - The URL status: `PENDING` or `ACTIVE`.

## Import

A MicroVM ingress URL can be imported using the MicroVM ID and the URL ID,
separated by a comma, e.g.

```
terraform import digitalocean_microvm_ingress.admin 506f78a4-e098-11e5-ad9f-000f53306ae1,b1c2d3e4-f5a6-7890-abcd-ef1234567890
```
