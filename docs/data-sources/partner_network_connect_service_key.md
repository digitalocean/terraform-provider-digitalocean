---
page_title: "DigitalOcean: digitalocean_partner_attachment_service_key"
subcategory: "Networking"
---

# digitalocean_partner_attachment_service_key

Retrieve the Service Key needed to create a Megaport Virtual Cross Connect (VXC) connecting the Partner Attachment with a Megaport Cloud Router (MCR). 

The Service Key is only available once the Partner Attachment is created and the Service Key is not part of the `digitalocean_partner_attachment` attributes. The `digitalocean_partner_attachment_service_key` data source is needed to retrieve the Service Key.

## Example Usage

```hcl
resource "digitalocean_partner_attachment" "example" {
  name                         = "example-partner-attachment"
  connection_bandwidth_in_mbps = 1000
  region                       = "nyc"
  naas_provider                = "MEGAPORT"
  vpc_ids                      = ["0bcef6a5-0000-0000-0000-000000000000"]
}

data "digitalocean_partner_attachment_service_key" "example" {
  attachment_id = digitalocean_partner_attachment.example.id
}

output "service_key" {
  value = data.digitalocean_partner_attachment_service_key.example.value
}
```

## Argument Reference

The following arguments are supported:

* `attachment_id` - (Required) The id of the Partner Attachment.

## Attributes Reference

* `value` - The value of the Service Key used with Megaport.
* `state` - The state of the Partner Attachment.
* `created_at` - The date and time of when the Partner Attachment was created.
