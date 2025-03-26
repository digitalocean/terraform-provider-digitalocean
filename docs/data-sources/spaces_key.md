---
page_title: "DigitalOcean: digitalocean_spaces_key"
subcategory: "Spaces Object Storage"
---

# digitalocean_spaces_key

Get information on a Spaces key for use in other resources. This is useful if the Spaces key in question
is not managed by Terraform or you need to utilize any of the key's data.

## Example Usage

Get the key by access key ID:

```hcl
data "digitalocean_spaces_key" "example" {
  access_key = "ACCESS_KEY_ID"
}

output "key_grants" {
  value = data.digitalocean_spaces_key.example.grant
}
```

## Argument Reference

The following arguments must be provided:

* `access_key` - (Required) The Access Key ID of the Spaces key.

## Attributes Reference

The following attributes are exported:

* `name` - The name of the Spaces key
* `grant` - The list of grants associated with the Spaces key.
* `access_key` - The access key ID of the Spaces key
* `created_at` - The creation time of the Spaces key
