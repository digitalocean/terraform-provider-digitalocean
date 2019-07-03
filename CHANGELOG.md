## 1.5.0 (July 03, 2019)

FEATURES:

* **New Data Source:** `digitalocean_database_cluster` ([#251](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/251)). Thanks to @stack72!

BUG FIXES:

* resource/digitalocean_droplet: DigitalOcean doesn't support IPv6 private networking. Mark the `ipv6_address_private` attribute in the schema as removed ([#181](https://github.com/terraform-providers/terraform-provider-digitalocean/issues/181)). Thanks to @stack72!
* resource/digitalocean_kubernetes_cluster: Do not filter out node pool tags also applied to the cluster ([#184](https://github.com/terraform-providers/terraform-provider-digitalocean/issues/184)).

## 1.4.0 (May 29, 2019)

IMPROVEMENTS:

* resource/digitalocean_droplet: Make importing more robust and test error case when importing non-existent resource ([#231](https://github.com/terraform-providers/terraform-provider-digitalocean/issues/231)).
* resource/digitalocean_spaces_bucket: Simplify importing and test error case when importing non-existent resource ([#231](https://github.com/terraform-providers/terraform-provider-digitalocean/issues/231)).
* resource/digitalocean_volume: Remove need for custom import function ([#231](https://github.com/terraform-providers/terraform-provider-digitalocean/issues/231)).

BUG FIXES:

* resource/digitalocean_record: Simplify importing and provide better error messaging when attempting to import non-existent resource ([#232](https://github.com/terraform-providers/terraform-provider-digitalocean/issues/232)).
* resource/digitalocean_kubernetes_cluster: Fix access to kube_config attributes under Terraform 0.12 by using TypeList rather than TypeSet in the schema ([#239](https://github.com/terraform-providers/terraform-provider-digitalocean/issues/239)).

## 1.3.0 (May 09, 2019)

IMPROVEMENTS:

* Terraform SDK upgrade with compatibility for Terraform v0.12.

## 1.2.0 (April 23, 2019)

FEATURES:

* **New Resource:** `digitalocean_spaces_bucket` ([#42](https://github.com/terraform-providers/terraform-provider-digitalocean/issues/42)). Thanks to @slapula!
* **New Resource:** `digitalocean_database_cluster` ([#198](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/198)) Thanks to @slapula!
* **New Resource:** `digitalocean_cdn` ([#204](https://github.com/terraform-providers/terraform-provider-digitalocean/issues/204))
* **New Resource:** `digitalocean_project` ([#207](https://github.com/terraform-providers/terraform-provider-digitalocean/issues/207))

IMPROVEMENTS:

* provider: The DigitalOcean API URL can now be overridden using `api_endpoint` attribute ([#84](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/84)).
  Thanks to @protochron!
* provider: Refactor logic for logging API requests/responses ([#190](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/190)). Thanks to @radeksimko!
* resource/digitalocean_loadbalancer: Add support for enabling PROXY Protocol ([#199](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/199)).
* docs/digitalocean_firewall: Update syntax to be compatible with Terraform 0.12-beta ([#201](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/201)).
* resource/digitalocean_droplet: Expose uniform resource name (URN) attribute for use with Projects resource ([#215](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/215)).
* resource/digitalocean_loadbalancer: Expose uniform resource name (URN) attribute for use with Projects resource ([#214](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/214)).
* resource/digitalocean_domain: Expose uniform resource name (URN) attribute for use with Projects resource ([#213](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/213)).
* resource/digitalocean_volume: Expose uniform resource name (URN) attribute for use with Projects resource ([#212](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/212)).
* resource/digitalocean_floating_ip: Expose uniform resource name (URN) attribute for use with Projects resource ([#211](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/211)).
* resource/digitalocean_spaces_bucket: Expose uniform resource name (URN) attribute for use with Projects resource ([#210](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/210)).

BUG FIXES:

* resource/digitalocean_certificate: Fix issue when using computed values for custom certificates ([#163](https://github.com/terraform-providers/terraform-provider-digitalocean/issues/163)).
* resource/digitalocean_droplet: Prevent unexpected rebuilds for Droplets created using image slugs when the backing image is updated ([#152](https://github.com/terraform-providers/terraform-provider-digitalocean/issues/152)).

NOTES:

* This provider is now built and tested using Go 1.11.x ([#178](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/178)).
* Dependencies for this provider are now managed using Go Modules ([#187](https://github.com/terraform-providers/terraform-provider-digitalocean/issues/187)).

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
