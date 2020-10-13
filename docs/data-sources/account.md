---
page_title: "DigitalOcean: digitalocean_account"
---

# digitalocean_account

Get information on your DigitalOcean account.

## Example Usage

Get the account:

```hcl
data "digitalocean_account" "example" {
}
```

## Attributes Reference

The following attributes are exported:

* `droplet_limit`: The total number of droplets current user or team may have active at one time.
* `floating_ip_limit`: The total number of floating IPs the current user or team may have.
* `email`: The email address used by the current user to register for DigitalOcean.
* `uuid`: The unique universal identifier for the current user.
* `email_verified`: If true, the user has verified their account via email. False otherwise.
* `status`: This value is one of "active", "warning" or "locked".
* `status_message`: A human-readable message giving more details about the status of the account.
