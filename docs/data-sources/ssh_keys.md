---
page_title: "DigitalOcean: digitalocean_ssh_keys"
---

# digitalocean_ssh_keys

Get information on SSH Keys for use in other resources.

This data source is useful if the SSH Keys in question are not managed by Terraform or you need to
utilize any of the SSH Keys' data.

Note: You can use the [`digitalocean_ssh_key`](droplet) data source to obtain metadata
about a single SSH Key if you already know the unique `name` to retrieve.

## Example Usage

For example to find all SSH Keys:

```hcl
data "digitalocean_ssh_keys" "keys" {
  sort {
    key       = "name"
    direction = "asc"
  }
}
```

## Argument Reference

* `filter` - (Optional) Filter the results.
  The `filter` block is documented below.

* `sort` - (Optional) Sort the results.
  The `sort` block is documented below.
 
`filter` supports the following arguments:

* `key` - (Required) Filter the SSH Keys by this key. This may be one of `name`, `public_key`, or `fingerprint`.

`sort` supports the following arguments:

* `key` - (Required) Sort the SSH Keys by this key. This may be one of `name`, `public_key`, or `fingerprint`.

* `direction` - (Required) The sort direction. This may be either `asc` or `desc`.

## Attributes Reference

* `ssh_keys` - A list of SSH Keys. Each SSH Key has the following attributes:  

  * `id` - The ID of the ssh key.
  * `name`: The name of the ssh key.
  * `public_key`: The public key of the ssh key.
  * `fingerprint`: The fingerprint of the public key of the ssh key.
