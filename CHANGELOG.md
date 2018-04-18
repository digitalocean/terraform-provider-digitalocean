## 0.1.4 (Unreleased)

IMPROVEMENTS:

* resource/digitalocean_record: Manage CAA domain records ([#48](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/48))

## 0.1.3 (December 18, 2017)

IMPROVEMENTS:

* provider: Report Terraform version via User-Agent ([#43](https://github.com/terraform-providers/terraform-provider-digitalocean/issues/43))
* resource/digitalocean_droplet: Add `monitoring` field ([#38](https://github.com/terraform-providers/terraform-provider-digitalocean/issues/38))

BUG FIXES:

* resource/digitalocean_droplet: Avoid crash on conditional volumes ([#40](https://github.com/terraform-providers/terraform-provider-digitalocean/issues/40))
* resource/digitalocean_droplet: Make sure we've got a proper IP address from DO ([#29](https://github.com/terraform-providers/terraform-provider-digitalocean/issues/29))
* resource/digitalocean_firewall: Correctly handle `destination_tags` ([#36](https://github.com/terraform-providers/terraform-provider-digitalocean/issues/36))
* resource/digitalocean_firewall: Suppress diff for 'all' port range ([#41](https://github.com/terraform-providers/terraform-provider-digitalocean/issues/41))

## 0.1.2 (July 31, 2017)

BUG FIXES:

* resource/digitalocean_droplet: Detaching the disks before deleting a droplet ([#22](https://github.com/terraform-providers/terraform-provider-digitalocean/issues/22))

## 0.1.1 (June 21, 2017)

NOTES:

Bumping the provider version to get around provider caching issues - still same functionality

## 0.1.0 (June 19, 2017)

FEATURES:

* **New Resource:** `digitalocean_firewall` ([#1](https://github.com/terraform-providers/terraform-provider-digitalocean/issues/1))
