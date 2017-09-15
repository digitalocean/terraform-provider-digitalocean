## 0.1.3 (Unreleased)

IMPROVEMENTS:

* resource/digitalocean_droplet: Add `monitoring` field [GH-38]

BUG FIXES:

* resource/digitalocean_droplet: Make sure we've got a proper IP address from DO [GH-29]
* resource/digitalocean_firewall: Correctly handle `destination_tags` [GH-36]
* resource/digitalocean_firewall: Suppress diff for 'all' port range [GH-41]

## 0.1.2 (July 31, 2017)

BUG FIXES:

* resource/digitalocean_droplet: Detaching the disks before deleting a droplet ([#22](https://github.com/terraform-providers/terraform-provider-digitalocean/issues/22))

## 0.1.1 (June 21, 2017)

NOTES:

Bumping the provider version to get around provider caching issues - still same functionality

## 0.1.0 (June 19, 2017)

FEATURES:

* **New Resource:** `digitalocean_firewall` ([#1](https://github.com/terraform-providers/terraform-provider-digitalocean/issues/1))
