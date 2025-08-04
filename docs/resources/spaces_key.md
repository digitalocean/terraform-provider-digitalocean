---
page_title: "DigitalOcean: digitalocean_spaces_key"
subcategory: "Spaces Object Storage"
---

# digitalocean\_spaces\_key

Provides a key resource for Spaces, DigitalOcean's object storage product.

The [Spaces API](https://docs.digitalocean.com/reference/api/spaces-api/) was
designed to be interoperable with Amazon's AWS S3 API. This allows users to
interact with the service while using the tools they already know. Spaces
mirrors S3's authentication framework and requests to Spaces require a key pair
similar to Amazon's Access ID and Secret Key.

As a Spaces owner, you limit others’ access to your buckets using Spaces access 
keys. Access keys can provide several levels of permissions to create, destroy,
read, and write to specific associated buckets. However, access keys only limit 
access to certain commands using the S3 API or CLI, not the control panel or 
other DigitalOcean resources.

## Example Usage

### Create a New Key

```hcl
resource "digitalocean_spaces_key" "foobar" {
  name = "foobar"
}
```

### Create a New Key with Grants

```hcl
resource "digitalocean_spaces_key" "foobar" {
  name = "foobar"
  grant {
    bucket     = "my-bucket"
    permission = "read"
  }
}
```

### Create a New Key with full access

```hcl
resource "digitalocean_spaces_key" "foobar" {
  name = "foobar"
  grant {
    bucket     = ""
    permission = "fullaccess"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the key
* `grant` - (Optional) A grant for the key (documented below).

The `grant` object supports the following:

* `bucket` - (Required) Name of the bucket associated with this grant. In case of a `fullaccess` permission, this value should be an empty string.
* `permission` - (Required) Permission associated with this grant. Values can be `read`, `readwrite`, `fullaccess`.

## Attributes Reference

The following attributes are exported:

* `name` - The name of the key
* `grant` - The list of grants associated with the key
* `access_key` - The access key ID of the key
* `secret_key` - The access key secret of the key
* `created_at` - The creation time of the key
