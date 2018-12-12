## 1.2.0 (Unreleased)
## 1.1.0 (December 12, 2018)

FEATURES:

* **New Data Source:** `digitalocean_droplet_snapshot` ([#161](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/161)). Thanks to @thefossedog!
* **New Data Source:** `digitalocean_kubernetes_cluster` ([#169](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/169)) Thanks to @nicholasjackson!
* **New Resource** `digitalocean_droplet_snapshot` ([#161](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/161)). Thanks to @thefossedog!
* **New Resource:** `digitalocean_kubernetes_cluster` ([#169](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/169)) Thanks to @nicholasjackson!
* **New Resource:** `digitalocean_kubernetes_node_pool` ([#169](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/169)) Thanks to @nicholasjackson!

## 1.0.2 (October 05, 2018)

BUG FIXES:

* resource/digitalocean_certificate: Suppress diff for DNS names on custom certificate resources ([#146](https://github.com/terraform-providers/terraform-provider-digitalocean/issues/146)).
* resource/digitalocean_floating_ip_assignment: Ensure resource works with the `create_before_destroy` lifecycle rule ([#147](https://github.com/terraform-providers/terraform-provider-digitalocean/issues/147)).

## 1.0.1 (October 02, 2018)

BUG FIXES:

* resource/digitalocean_droplet: Ensure the image ID is set to state when importing a Droplet ([#144](https://github.com/terraform-providers/terraform-provider-digitalocean/issues/144))

## 1.0.0 (September 27, 2018)

FEATURES:

*  **New Resource** `digitalocean_floating_ip_assignment` ([#115](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/115)) Thanks to @justinbarrick!
*  **New Resource** `digitalocean_volume_attachment` ([#130](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/130))

* **New Datasource:** `digitalocean_domain` ([#63](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/63)) Thanks to @slapula!
* **New Datasource:** `digitalocean_record` ([#64](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/64)) Thanks to @slapula!
* **New Datasource:** `digitalocean_certificate` ([#137](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/137))
* **New Datasource:** `digitalocean_droplet` ([#137](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/137))
* **New Datasource:** `digitalocean_floating_ip` ([#137](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/137))
* **New Datasource:** `digitalocean_loadbalancer` ([#137](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/137))
* **New Datasource:** `digitalocean_ssh_key` ([#137](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/137))
* **New Datasource:** `digitalocean_tag` ([#137](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/137))
* **New Datasource:** `digitalocean_volume` ([#137](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/137))
* **New Datasource:** `digitalocean_volume_snapshot` ([#139](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/139))

IMPROVEMENTS:

* resource/digitalocean_record: Manage CAA domain records ([#48](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/48)). Thanks to @jaymecd!
* resource/digitalocean_certificate: Existing resources are now importable ([#37](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/37)). Thanks to @jonnydford!
* resource/digitalocean_loadbalancer: Existing resources are now importable ([#37](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/37)). Thanks to @jonnydford!
* resource/digitalocean_record: Existing resources are now importable ([#71](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/71)). Thanks to @slapula!
* resource/digitalocean_tag: Validate tag name when creating ([#80](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/80)). Thanks to @inkel!
* resource/digitalocean_volume: Add support for filesystem type for volume resources ([#111](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/111)). Thanks to @pgrzesik!
* resource/digitalocean_certificate: Added support for LetsEncrypt issued certificates. ([#129](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/129))
* resource/digitalocean_volume: Added support for volume resizing. ([#125](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/125))
* resource/digitalocean_volume: Added support for creating a volume from a volume snapshot. ([#139](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/139))
* resource/digitalocean_droplet: Updated the state representation of the `user_data` property to a hashed one. ([#128](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/128))

BUG FIXES:

* resource/digitalocean_floating_ip: Gracefully handle missing Floating IPs. ([#55](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/55)). Thanks to @aknuds1!
* resource/digitalocean_floating_ip: Gracefully handle unassigning Floating IPs when the Droplet has already been destroyed. ([#57](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/57)). Thanks to @aknuds1!
* resource/digitalocean_droplet: When IPv6 and/or private networking are not enabled, default their addresses to "". ([#97](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/97))
* resource/digitalocean_droplet: Don't panic when enabling private_networking or ipv6 on an existing Droplet.  ([#94](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/94))
* resource/digitalocean_droplet: Set resize_disk to the default value on import. ([#95](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/95))
* resource/digitalocean_record: Fixed the function for generating a records FQDN to match DO API behaviour. ([#51](https://github.com/terraform-providers/terraform-provider-digitalocean/issues/51))
* resource/digitalocean_firewall: Refactored the firewall resource for better stability of diffs and updates. ([#133](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/133))
* resource/digitalocean_record_test: Enable setting a DNS' record weight to 0 (via updated godo SDK). ([#132](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/132))


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
