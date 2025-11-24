---
page_title: "DigitalOcean: digitalocean_byoip_prefix"
subcategory: "Networking"
---

# digitalocean_byoip_prefix

Provides a DigitalOcean BYOIP (Bring Your Own IP) prefix resource. This can be used to
create, modify, and delete BYOIP prefixes.

BYOIP prefixes allow you to bring your own IP address space to DigitalOcean. You can
use this feature to maintain your IP reputation or meet specific compliance requirements.

Note: By default, newly provisioned BYOIP prefixes are not advertised to the internet. After the initial `terraform apply`, BYOIP provisioning request is initiated and DigitalOcean provisions the prefix, the prefix status changes to Active. At this point, you can initiate advertising prefix to the internet by setting field `advertised = true` and apply the configuration to make your prefix fully usable and accessible from the internet. 

## Example Usage

```hcl
# Create a new BYOIP prefix
resource "digitalocean_byoip_prefix" "example" {
  prefix     = "192.0.2.0/24"
  signature  = var.prefix_signature
  region     = "nyc3"
  advertised = false
}
```

## Argument Reference

The following arguments are supported:

* `prefix` - (Required) The CIDR notation of the prefix (e.g., "192.0.2.0/24").
* `signature` - (Required) The cryptographic signature proving ownership of the prefix.
  This is required during creation but can be omitted in subsequent updates.
* `region` - (Required) The DigitalOcean region where the prefix will be deployed.
* `advertised` - (Optional) A boolean indicating whether the prefix should be advertised.
  Defaults to `false`.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:

* `id`: The UUID of the BYOIP prefix.
* `uuid`: The UUID of the BYOIP prefix.
* `status`: The current status of the BYOIP prefix (e.g., "verified", "pending", "failed").
* `failure_reason`: The reason for failure if the status is "failed".

## Import

BYOIP prefixes can be imported using the prefix `uuid`, e.g.

```
terraform import digitalocean_byoip_prefix.example 506f78a4-e098-11e5-ad9f-000f53306ae1
```
