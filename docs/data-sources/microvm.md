---
page_title: "DigitalOcean: digitalocean_microvm"
subcategory: "MicroVMs"
---

# digitalocean_microvm

Get information on a MicroVM for use in other resources. This is useful if the
MicroVM is not managed by Terraform or you need to use any of its data.

**Note:** This data source returns a single MicroVM. When looking up by `name`,
an error is triggered if more than one MicroVM has that name.

-> **Note:** MicroVMs are in public preview.

## Example Usage

Get the MicroVM by name:

```hcl
data "digitalocean_microvm" "example" {
  name = "example-microvm"
}

output "hostname" {
  value = data.digitalocean_microvm.example.urls[0].hostname
}
```

Get the MicroVM by ID:

```hcl
data "digitalocean_microvm" "example" {
  id = "506f78a4-e098-11e5-ad9f-000f53306ae1"
}
```

## Argument Reference

One of the following arguments must be provided:

* `id` - (Optional) The ID of the MicroVM.
* `name` - (Optional) The name of the MicroVM.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the MicroVM.
* `name` - The name of the MicroVM.
* `urn` - The uniform resource name (URN) for the MicroVM.
* `region` - The slug of the region the MicroVM is in.
* `size` - The compute size of the MicroVM:
  - `cpu` - The number of vCPUs.
  - `memory` - The amount of memory in MiB.
  - `disk` - The size of the attached disk in GB.
* `source` - What the MicroVM runs:
  - `oci_ref` - The OCI reference of the workload container image.
  - `checkpoint_id` - The ID of the checkpoint the MicroVM was restored from.
* `networking` - The networking mode: `public` or `vpc`.
* `http_protocol` - The HTTP protocol the workload serves: `http` or `http2`.
* `ports` - The guest ports open for ingress.
* `auto_pause` - Auto-pause configuration:
  - `enabled` - Whether auto-pause is enabled.
  - `idle_timeout` - How long the MicroVM must be idle before it pauses.
* `auto_resume` - Whether the MicroVM resumes automatically on incoming HTTP traffic.
* `tags` - The tags applied to the MicroVM.
* `current_state` - The observed lifecycle state of the MicroVM, such as `running` or `paused`.
* `failure_reason` - A description of the failure when `current_state` is `failed`.
* `urls` - A list of ingress URLs for the MicroVM:
  - `hostname` - The hostname, without a scheme.
  - `port` - The guest port the URL forwards to.
  - `default` - Whether this is the default URL.
  - `status` - The URL status: `PENDING` or `ACTIVE`.
* `created_at` - The date and time the MicroVM was created.
