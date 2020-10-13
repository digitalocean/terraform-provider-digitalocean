---
page_title: "DigitalOcean: digitalocean_spaces_bucket"
---

# digitalocean_spaces_bucket

Get information on a Spaces bucket for use in other resources. This is useful if the Spaces bucket in question
is not managed by Terraform or you need to utilize any of the bucket's data.

## Example Usage

Get the bucket by name:

```hcl
data "digitalocean_spaces_bucket" "example" {
  name = "my-spaces-bucket"
  region = "nyc3"
}

output "bucket_domain_name" {
  value = data.digitalocean_spaces_bucket.example.bucket_domain_name
}
```

## Argument Reference

The following arguments must be provided:

* `name` - (Required) The name of the Spaces bucket.
* `region` - (Required) The slug of the region where the bucket is stored.

## Attributes Reference

The following attributes are exported:

* `name` - The name of the Spaces bucket
* `region` - The slug of the region where the bucket is stored.
* `urn` - The uniform resource name of the bucket
* `bucket_domain_name` - The FQDN of the bucket (e.g. bucket-name.nyc3.digitaloceanspaces.com)
