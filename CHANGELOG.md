## 2.17.0 (January 14, 2022)

IMPROVEMENTS:

- `digitalocean_loadbalancer`: Fetch loadbalancer resource in datasource by ID ([#773](https://github.com/digitalocean/terraform-provider-digitalocean/pull/773)). @opeco17
- `digitalocean_vpc`: Allow updating name and description of default vpcs ([#748](https://github.com/digitalocean/terraform-provider-digitalocean/pull/748)). - @andrewsomething
- `digitalocean_app`: Support preserve_path_prefix ([#768](https://github.com/digitalocean/terraform-provider-digitalocean/pull/768)). - @andrewsomething

BUG FIXES:

- `digitalocean_droplet` - Refactor post-create polling code ([#776](https://github.com/digitalocean/terraform-provider-digitalocean/pull/776)). - @andrewsomething
- `digitalocean_floating_ip_assignment`: Properly support importing existing assignments ([#771](https://github.com/digitalocean/terraform-provider-digitalocean/pull/771)). - @andrewsomething
- `digitalocean_database_cluster`: Retry on 404s in post create polling ([#761](https://github.com/digitalocean/terraform-provider-digitalocean/pull/761)). - @opeco17

MISC:

- docs: Update records examples to use domain id over name ([#770](https://github.com/digitalocean/terraform-provider-digitalocean/pull/770)). - @andrewsomething
- docs: Add k8s as available project resource to docs ([#750](https://github.com/digitalocean/terraform-provider-digitalocean/pull/750)) - @scotchneat
- docs: Fix database cluster documentation ([#766](https://github.com/digitalocean/terraform-provider-digitalocean/pull/766)) - @colinwilson
- provider: Update to v2.10.1 of the terraform-plugin-sdk ([#760](https://github.com/digitalocean/terraform-provider-digitalocean/pull/760)). - @andrewsomething
- testing: Makefile: Only run sweep against the digitalocean package ([#759](https://github.com/digitalocean/terraform-provider-digitalocean/pull/759)). - @andrewsomething
- testing: Update domains sweeper ([#753](https://github.com/digitalocean/terraform-provider-digitalocean/pull/753)) - @scotchneat


## 2.16.0 (November 8, 2021)

IMPROVEMENTS:

- `digitalocean_loadbalancer`: Add support for size_unit ([#742](https://github.com/digitalocean/terraform-provider-digitalocean/pull/742)). - @bbassingthwaite
- `digitalocean_database_firewall`: Add attributes for Kubernetes cluster IDs to firewall rules ([#741](https://github.com/digitalocean/terraform-provider-digitalocean/pull/741)). - @tdyas

BUG FIXES:

- `digitalocean_database_user`, `digitalocean_database_replica`, `digitalocean_database_db`, `digitalocean_database_connection_pool`: Provide better error messages importing database sub-resources. ([#744](https://github.com/digitalocean/terraform-provider-digitalocean/pull/744)). - @andrewsomething

## 2.15.0 (November 1, 2021)

IMPROVEMENTS:

- `digitalocean_container_registry_docker_credentials`: Revoke OAuth token when credentials are destroyed ([#735](https://github.com/digitalocean/terraform-provider-digitalocean/pull/735)). - @andrewsomething
- `digitalocean_loadbalancer`: Support disabling automatic DNS records when using Let's Encrypt certificates ([#723](https://github.com/digitalocean/terraform-provider-digitalocean/pull/723), [#730](https://github.com/digitalocean/terraform-provider-digitalocean/pull/730)). - @andrewsomething

BUG FIXES:

- docs: Remove outdated info from README ([#733](https://github.com/digitalocean/terraform-provider-digitalocean/pull/733)). - @andrewsomething

TESTING:

- testing: Check for endpoint attr in `digitalocean_kubernetes_cluster` when HA is enabled ([#725](https://github.com/digitalocean/terraform-provider-digitalocean/pull/725)). - @danaelhe
- testing: Add acceptance test workflow ([#732](https://github.com/digitalocean/terraform-provider-digitalocean/pull/732)). - @scotchneat
- testing: Add scheduled acceptance test runs ([#734](https://github.com/digitalocean/terraform-provider-digitalocean/pull/734)). - @scotchneat
- testing: Limit acceptance test job/workflow to single run at a time ([#738](https://github.com/digitalocean/terraform-provider-digitalocean/pull/738)). - @scotchneat
- testing: Acceptance tests: Don't use deprecated size slugs ([#737](https://github.com/digitalocean/terraform-provider-digitalocean/pull/737)). - @andrewsomething
- testing: Acceptance tests: Use 2048 bit private keys in test data ([#736](https://github.com/digitalocean/terraform-provider-digitalocean/pull/736)). - @andrewsomething
- testing: Fixes an invalid droplet size in some acceptance tests ([#724](https://github.com/digitalocean/terraform-provider-digitalocean/pull/724)). - @scotchneat

## 2.14.0 (October 7, 2021)

IMPROVEMENTS:

- `digitalocean_kubernetes_cluster`: Support setting the `ha` attribute ([#718](https://github.com/digitalocean/terraform-provider-digitalocean/pull/718))

## 2.13.0 (October 7, 2021)

BUG FIXES:

- Fix tag collection in digitalocean_tags data source ([#716](https://github.com/digitalocean/terraform-provider-digitalocean/pull/716))

FEATURES:

- Add digitalocean_database_ca data source. ([#717](https://github.com/digitalocean/terraform-provider-digitalocean/pull/717))

IMPROVEMENTS:

- Give a name to the kubernetes example load balancer ([#703](https://github.com/digitalocean/terraform-provider-digitalocean/pull/703))
- Gracefully shutdown droplet before deleting  ([#719](https://github.com/digitalocean/terraform-provider-digitalocean/pull/719))

## 2.12.1 (October 1, 2021)

BUGFIXES:

* docs: Correct the example for `digitalocean_monitor_alert` ([#710](https://github.com/digitalocean/terraform-provider-digitalocean/pull/710)). Thanks to @jessedobbelaere!

## 2.12.0 (September 22, 2021)

FEATURES:

- **New Resource:** `digitalocean_monitor_alert` ([#679](https://github.com/digitalocean/terraform-provider-digitalocean/pull/679)) Thanks to @atombrella!

IMPROVEMENTS:

- `digitalocean_domain`: Expose TTL ([#702](https://github.com/digitalocean/terraform-provider-digitalocean/pull/702)). Thanks to @atombrella!
- `digitalocean_app`: Support setting CORS policies ([#699](https://github.com/digitalocean/terraform-provider-digitalocean/pull/699)).
- `digitalocean_app`: Make create timeout configurable ([#698](https://github.com/digitalocean/terraform-provider-digitalocean/pull/698)).
- `digitalocean_droplet`: Mark `private_networking` as deprecated. ([#676](https://github.com/digitalocean/terraform-provider-digitalocean/issues/676))
- docs: Provide more context for apps' `instance_size_slug` ([#701](https://github.com/digitalocean/terraform-provider-digitalocean/pull/701))
- misc: Replace d.HasChange sequences with d.HasChanges  ([#681](https://github.com/digitalocean/terraform-provider-digitalocean/pull/681)) Thanks to @atombrella!

BUGFIXES:

- `digitalocean_database_user`: Handle passwords for MongoDB ([#696](https://github.com/digitalocean/terraform-provider-digitalocean/issues/696)).
- `digitalocean_app`: Error to prevent panic if no deployment found ([#678](https://github.com/digitalocean/terraform-provider-digitalocean/issues/678)).
- `digitalocean_droplet`: Protect against panic when importing Droplet errors ([#674](https://github.com/digitalocean/terraform-provider-digitalocean/issues/674)).

## 2.11.1 (August 20, 2021)

BUG FIXES:

- `digitalocean_record`: Move port validation for SRV records out of CustomizeDiff ([#670](https://github.com/digitalocean/terraform-provider-digitalocean/issues/670)).
- `digitalocean_record`: Fix unexpected diffs for when TXT records are for the apex domain ([#664](https://github.com/digitalocean/terraform-provider-digitalocean/issues/664)).

## 2.11.0 (August 9, 2021)

IMPROVEMENTS:
- `digitalocean_droplet`: Support setting the droplet_agent attribute ([#667](https://github.com/digitalocean/terraform-provider-digitalocean/pull/667)).
- `digitalocean_database_firewall`: Allow setting app as a type ([#666](https://github.com/digitalocean/terraform-provider-digitalocean/issues/666)).
- docs: Update links to API documentation ([#661](https://github.com/digitalocean/terraform-provider-digitalocean/pull/661)).
- Simplified some loops based on go-staticcheck S1011 ([#656](https://github.com/digitalocean/terraform-provider-digitalocean/pull/656)). Thanks to @atombrella!
- Update to Context-aware API v2 functions ([#657](https://github.com/digitalocean/terraform-provider-digitalocean/pull/657)). Thanks to @atombrella!

## 2.10.1 (June 29, 2021)

BUG FIXES:

- docs: Add code fence to docs for MongoDB.

## 2.10.0 (June 29, 2021)

IMPROVEMENTS:

- `digitalocean_kubernetes_cluster`:  Add support for Kubernetes maintenance policy ([#631](https://github.com/digitalocean/terraform-provider-digitalocean/pull/631)). Thanks to @atombrella!
- `digitalocean_kubernetes_cluster`, `digitalocean_database_cluster`: Make create timeouts configurable for DBaaS and K8s clusters ([#650](https://github.com/digitalocean/terraform-provider-digitalocean/pull/650)).
- docs: Fix Firewall Resource exported attributes documentation ([#648](https://github.com/digitalocean/terraform-provider-digitalocean/pull/648)). Thanks to @jubairsaidi!
- docs: Add MongoDB to database cluster docs ([#651](https://github.com/digitalocean/terraform-provider-digitalocean/pull/651)).

BUG FIXES:

- `digitalocean_database_cluster`: Protect against setting empty tags ([#652](https://github.com/digitalocean/terraform-provider-digitalocean/pull/652)).


## 2.9.0 (May 28, 2021)

IMPROVEMENTS:

- provider: Upgrade Terraform SDK to v2.6.1 ([#642](https://github.com/digitalocean/terraform-provider-digitalocean/pull/642)).
- `digitalocean_kubernetes_cluster`: Expose URNs for Kubernetes clusters ([#626](https://github.com/digitalocean/terraform-provider-digitalocean/pull/626)).
- docs: Cover `required_providers` in the index page ([#643](https://github.com/digitalocean/terraform-provider-digitalocean/pull/643)).
- docs: Update nginx example to work ([#625](https://github.com/digitalocean/terraform-provider-digitalocean/pull/625)). Thanks to @atombrella!
- Update issue templates ([#633](https://github.com/digitalocean/terraform-provider-digitalocean/pull/633)).

BUG FIXES:

- `digitalocean_loadbalancer`: Add `certificate_name` to load balancer data source ([#641](https://github.com/digitalocean/terraform-provider-digitalocean/pull/641)).
- `digitalocean_droplet`: Changing SSH key for a Droplet should be `ForceNew` ([#640](https://github.com/digitalocean/terraform-provider-digitalocean/pull/640)).
- `digitalocean_spaces_bucket`: Provide better error on malformated import ([#629](https://github.com/digitalocean/terraform-provider-digitalocean/pull/629)).

## 2.8.0 (April 20, 2021)

FEATURES:

- **New Data Source**: `digitalocean_database_replica` ([#489](https://github.com/digitalocean/terraform-provider-digitalocean/issues/489))

IMPROVEMENTS:

- `digitalocean_spaces_bucket`, `digitalocean_spaces_bucket_object`: Validate input for region ([#618](https://github.com/digitalocean/terraform-provider-digitalocean/pull/618)).
- `digitalocean_custom_image`: Support distributing custom images to multiple regions ([#616](https://github.com/digitalocean/terraform-provider-digitalocean/pull/616)).
- `digitalocean_custom_image`: Surface a better error to users if image import fails ([#613](https://github.com/digitalocean/terraform-provider-digitalocean/pull/613)).
- `digitalocean_database_cluster`: Support MongoDB beta by handling password differences ([#614](https://github.com/digitalocean/terraform-provider-digitalocean/pull/614)).

## 2.7.0 (March 29, 2021)

IMPROVEMENTS:

* `digitalocean_kubernetes_cluster`, `digitalocean_kubernetes_node_pool`: Support for Kubernetes node pool taints ([#374](https://github.com/digitalocean/terraform-provider-digitalocean/issues/374)).
* `digitalocean_loadbalancer`: Support resizing load balancers ([#606](https://github.com/digitalocean/terraform-provider-digitalocean/issues/606)).

BUG FIXES:

* docs: Fix Kubernetes autoscaling docs for `min_nodes` ([#602](https://github.com/digitalocean/terraform-provider-digitalocean/pull/602)). Thanks to @3dinfluence!

## 2.6.0 (March 10, 2021)

NOTES:
* With the update to go 1.16 ([#597](https://github.com/digitalocean/terraform-provider-digitalocean/pull/597)),
  the provider now supports `darwin_arm64`.

FEATURES:
* `datasource_digitalocean_firewall`: Adds Firewall datasource ([#594](https://github.com/digitalocean/terraform-provider-digitalocean/pull/594))

IMPROVEMENTS:
* Run tests on pull_request not pull_request_target. ([#589](https://github.com/digitalocean/terraform-provider-digitalocean/pull/589))
* kubernetes - enable surge upgrades by default during cluster creation ([#584](https://github.com/digitalocean/terraform-provider-digitalocean/pull/584))
* Assign and remove project resources without unnecessary churn (Fixes: #585). ([#586](https://github.com/digitalocean/terraform-provider-digitalocean/pull/586))
* dbaas replica: Add missing attrbutes to docs. ([#588](https://github.com/digitalocean/terraform-provider-digitalocean/pull/588))
* Bump Kubernetes version used in documentation ([#583](https://github.com/digitalocean/terraform-provider-digitalocean/pull/583))

BUG FIXES:
* Fix broken documentation links ([#592](https://github.com/digitalocean/terraform-provider-digitalocean/pull/592))
* Fix docs and validation for expiry_seconds on registry docker credentials resource. ([#582](https://github.com/digitalocean/terraform-provider-digitalocean/pull/582))

## 2.5.1 (February 05, 2021)

BUG FIXES:

* `digitalocean_database_cluster`: Protect against panic if connection details not available. ([#577](https://github.com/digitalocean/terraform-provider-digitalocean/pull/577)).
* `digitalocean_cdn`: Handle certificate name updates. ([#579](https://github.com/digitalocean/terraform-provider-digitalocean/pull/579)).

## 2.5.0 (February 03, 2021)

NOTES:

* `digitalocean_app`: In order to support additional features, the `domains` attribute has been deprecated and will be removed in a future release. It has been replaced by a repeatable `domain` block which supports wildcard domains and specifying DigitalOcean managed zones.

IMPROVEMENTS:

* `digitalocean_app`: Deprecate domains list in favor of domain block. ([#572](https://github.com/digitalocean/terraform-provider-digitalocean/issues/572)).
* `digitalocean_app`: Add support for images as a component source ([#565](https://github.com/digitalocean/terraform-provider-digitalocean/issues/565)). Thanks to @rienafairefr and @acraven!
* `digitalocean_app`: Add support for job components ([#566](https://github.com/digitalocean/terraform-provider-digitalocean/issues/566)).  Thanks to @rienafairefr and @acraven!
* `digitalocean_app`: Add support for `internal_ports` ([#570](https://github.com/digitalocean/terraform-provider-digitalocean/issues/570)). Thanks to @rienafairefr!

BUG FIXES:

* `digitalocean_app`: Allow multiple routes for services and static sites ([#571](https://github.com/digitalocean/terraform-provider-digitalocean/issues/571)).

## 2.4.0 (January 19, 2021)

IMPROVEMENTS:

* `digitalocean_app`: Add support for global env vars ([#549](https://github.com/digitalocean/terraform-provider-digitalocean/issues/549)).
* `digitalocean_app`: Add GitLab support ([#556](https://github.com/digitalocean/terraform-provider-digitalocean/issues/556)).
* `digitalocean_app`: Support `catchall_document` for static sites ([#539](https://github.com/digitalocean/terraform-provider-digitalocean/issues/539)).
* `digitalocean_custom_image`: Support updating `description` and `distribution` ([#538](https://github.com/digitalocean/terraform-provider-digitalocean/issues/538)). Thanks to @frezbo!

BUG FIXES:

* `digitalocean_vpc`: Protect against race conditions in IP range assignment ([#552](https://github.com/digitalocean/terraform-provider-digitalocean/issues/552)).
* `digitalocean_app`: Mark env var values as sensitive ([#554](https://github.com/digitalocean/terraform-provider-digitalocean/issues/554)).

## 2.3.0 (December 03, 2020)

IMPROVEMENTS:

* provider: Build and release OpenBSD binaries ([#533](https://github.com/digitalocean/terraform-provider-digitalocean/issues/533)).
* `digitalocean_loadbalancer`: Add support for new `size` attribute ([#532](https://github.com/digitalocean/terraform-provider-digitalocean/issues/532)). Thanks to @anitgandhi!

BUG FIXES:

* `digitalocean_database_cluster`: Handle Redis version change with DiffSuppressFunc ([#534](https://github.com/digitalocean/terraform-provider-digitalocean/issues/534)).

## 2.2.0 (November 06, 2020)

FEATURES:

* **New Data Source**: `digitalocean_ssh_keys` ([#519](https://github.com/digitalocean/terraform-provider-digitalocean/pull/519)) Thanks to @stack72!
* **New Resource**: `digitalocean_custom_image` ([#517](https://github.com/digitalocean/terraform-provider-digitalocean/pull/517)) Thanks to @frezbo!

BUG FIXES:

*  `digitalocean_kubernetes_node_pool`: Validate min_nodes is at least 1 and fix example ([#525](https://github.com/digitalocean/terraform-provider-digitalocean/pull/525)).
* `digitalocean_app`: Document the database spec ([#524](https://github.com/digitalocean/terraform-provider-digitalocean/pull/524)).
* `digitalocean_container_registry`: Update docs w/ `subscription_tier_slug` ([#523](https://github.com/digitalocean/terraform-provider-digitalocean/pull/523)).

## 2.1.0 (November 06, 2020)

NOTES:

* DigitalOcean Container Registry is now in general availablity and requires a [subscription plan](https://www.digitalocean.com/docs/container-registry/#plans-and-pricing). As a result, the `digitalocean_container_registry` resource now requires setting a `subscription_tier_slug`.

IMPROVEMENTS:

* `digitalocean_container_registry`: Supports setting and updating a `subscription_tier_slug` ([#516](https://github.com/digitalocean/terraform-provider-digitalocean/pull/516)).

BUG FIXES:

* `digitalocean_app`: Suppress diff when env type is `GENERAL` ([#515](https://github.com/digitalocean/terraform-provider-digitalocean/pull/515)).

## 2.0.2 (October 28, 2020)

BUG FIXES:

* `digitalocean_spaces_bucket`: Add retry logic to ensure bucket is available before proceeding ([#510](https://github.com/digitalocean/terraform-provider-digitalocean/issues/510)).
* Docs: Fix broken link to DigitalOcean app spec ([#509](https://github.com/digitalocean/terraform-provider-digitalocean/pull/509)). Thanks to @edbedbe!

## 2.0.1 (October 22, 2020)

BUG FIXES:

* `digitalocean_cdn`, `digitalocean_app`: Fix panics introduced in move from ReadFunc to ReadContextFunc ([#505](https://github.com/digitalocean/terraform-provider-digitalocean/issues/505)).

## 2.0.0 (October 20, 2020)

NOTES:

* This release uses v2.0.3 of the Terraform Plugin SDK and now only supports Terraform v0.12 and higher.
* The `certificate_id` attribute of the `digitalocean_cdn` and `digitalocean_loadbalancer` resources has been deprecated in favor of `certificate_name`. It will become a read-only computed attrbute in a future release.

FEATURES:

* **New Data Source**: `digitalocean_records` ([#502](https://github.com/digitalocean/terraform-provider-digitalocean/pull/502)) Thanks to @tdyas!

IMPROVEMENTS:

* provider: Upgrade to v2.0.3 of the Terraform Plugin SDK ([#492](https://github.com/digitalocean/terraform-provider-digitalocean/pull/492), [#503](https://github.com/digitalocean/terraform-provider-digitalocean/pull/503)). Thanks to @tdyas!
* docs: Migrate documentation to new registry format ([#501](https://github.com/digitalocean/terraform-provider-digitalocean/pull/501)).

BUG FIXES:

* `digitalocean_certificate`, `digitalocean_cdn`, `digitalocean_loadbalancer`: Use certificate name as primary identifier instead of ID as a Let's Encrypt certificate's ID will change when it's auto-renewed ([#500](https://github.com/digitalocean/terraform-provider-digitalocean/pull/500)).

## 1.23.0 (October 13, 2020)

FEATURES:

* **New Resource**: `digitalocean_app` ([#497](https://github.com/digitalocean/terraform-provider-digitalocean/pull/497))
* **New Data Source**: `digitalocean_app` ([#497](https://github.com/digitalocean/terraform-provider-digitalocean/pull/497))
* **New Data Source**: `digitalocean_domains` ([#484](https://github.com/digitalocean/terraform-provider-digitalocean/pull/484)) Thanks to @tdyas!

IMPROVEMENTS:

* All list-style data sources now support `all` and `match_by` attributes for filter blocks ([#491](https://github.com/digitalocean/terraform-provider-digitalocean/pull/491)) and other improvements ([#481](https://github.com/digitalocean/terraform-provider-digitalocean/pull/491)). Thanks to @tdyas!
* Additional acceptance test sweepers ([#458](https://github.com/digitalocean/terraform-provider-digitalocean/pull/458)).
* Upgrade to v1.15.0 of the terraform-plugin-sdk ([#483](https://github.com/digitalocean/terraform-provider-digitalocean/pull/458)). Thanks to @tdyas!

## 1.22.1 (August 7, 2020)

BUG FIXES:

* resource/digitalocean_record: Allow for SRV records with port 0 ([#475](https://github.com/digitalocean/terraform-provider-digitalocean/issues/475)).

## 1.22.0 (July 22, 2020)

IMPROVEMENTS:

* resource/digitalocean_kubernetes_cluster: Enable auto_upgrade on Kubernetes clusters ([#237](https://github.com/digitalocean/terraform-provider-digitalocean/issues/237)). Thanks to @lfarnell!
* resource/digitalocean_kubernetes_cluster: Add support for surge upgrades ([#465](https://github.com/digitalocean/terraform-provider-digitalocean/pull/465)). Thanks to @varshavaradarajan!

BUG FIXES:

* resource/digitalocean_container_registry_docker_credentials: Update default expiry and don't ignore error ([#467](https://github.com/digitalocean/terraform-provider-digitalocean/pull/467)).

## 1.21.0 (July 20, 2020)

IMPROVEMENTS:

* resource/digitalocean_loadbalancer: Adds 'https' to list of acceptable healthcheck protocols ([#460](https://github.com/digitalocean/terraform-provider-digitalocean/pull/460)). Thanks to @gcox!
* provider: Update module and import paths for repository transfer ([#463](https://github.com/digitalocean/terraform-provider-digitalocean/pull/463)).

BUG FIXES:

* resource/digitalocean_vpc: Increase timeout on VPC deletion retry ([#455](https://github.com/digitalocean/terraform-provider-digitalocean/pull/455)).

## 1.20.0 (June 15, 2020)

FEATURES:

* **New Data Source**: `digitalocean_tags` ([#451](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/451)).

IMPROVEMENTS:

* resource/digitalocean_tag, datasource/digitalocean_tag: Export the counts of tagged resources as a computed attribute ([#451](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/451)).

BUG FIXES:

* datasource/digitalocean_droplets: Set ID in `flattenDigitalOceanDroplet` to ensure the individual Droplets have their ID exported in the list data source. [#450](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/450)).

## 1.19.0 (June 03, 2020)

FEATURES:

* **New Resources**: `digitalocean_container_registry`, `digitalocean_container_registry_docker_credentials` ([#335](https://github.com/terraform-providers/terraform-provider-digitalocean/issues/335)). Thanks to Zelgius!
* **New Data Source**: `digitalocean_container_registry` ([#335](https://github.com/terraform-providers/terraform-provider-digitalocean/issues/335)). Thanks to Zelgius!

IMPROVEMENTS:

* resource/digitalocean_database_replica: Add support for specifying a VPC for read-only replicas ([#440](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/440)).

BUG FIXES:

* resource/digitalocean_kubernetes_cluster: Add forcenew to vpc_uuid field ([#443](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/443)). Thanks to KadenLNelson!
* resource/digitalocean_kubernetes_cluster: Fail faster on cluster create error [#435](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/435)).

## 1.18.0 (May 05, 2020)

FEATURES:

* resource/digitalocean_loadbalancer, datasource/digitalocean_loadbalancer: Add support for the backend keepalive option ([#427](https://github.com/terraform-providers/terraform-provider-digitalocean/issues/427)).

BUG FIXES:

* provider: Spaces API Endpoint setting is optional ([#431](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/431)).

## 1.17.0 (April 28, 2020)

FEATURES:

* **New Data Sources**: `digitalocean_spaces_bucket` and `digitalocean_spaces_buckets` ([#416](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/416)) Thanks to @tdyas!
* **New Data Sources**: `digitalocean_spaces_bucket_object` and `digitalocean_spaces_bucket_objects` ([#423](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/423)) Thanks to @tdyas!
* **New Data Sources**: `digitalocean_droplets` ([#418](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/418)) Thanks to @tdyas!

BUG FIXES:

* resource/digitalocean_record: Fix handling of CAA records with iodef tag ([#421](https://github.com/terraform-providers/terraform-provider-digitalocean/issues/421)).
* resource/digitalocean_loadbalancer: Fix support for multiple forwarding rules ([#414](https://github.com/terraform-providers/terraform-provider-digitalocean/issues/414)).

## 1.16.0 (April 14, 2020)

FEATURES:

* **New Resource**: digitalocean_vpc ([#410](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/410))
* **New Data Source**: digitalocean_vpc ([#410](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/410))
* **New Resource**: digitalocean_spaces_bucket_object ([#408](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/408)) Thanks to @tdyas!
* resource/digitalocean_spaces_bucket: Support for bucket versioning ([#409](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/409)). Thanks to @tdyas!
* resource/digitalocean_spaces_bucket: Support for lifecycle rules ([#411](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/411)). Thanks to @tdyas!

IMPROVEMENTS:

* provider: Support overriding the Spaces API endpoint using the `spaces_endpoint` attribute or `SPACES_ENDPOINT_URL` environment variable. ([#384]( https://github.com/terraform-providers/terraform-provider-digitalocean/issues/384)). Thanks to @tdyas!

BUG FIXES:

* resource/digitalocean_volume: Revert local name validation to support volumes created before DigitalOcean API change ([#406](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/406)). Thanks to @Amygos & @BrianHicks!
* resource/digitalocean_spaces_bucket: Region attribute should be ForceNew ([#413](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/413)).

## 1.15.1 (March 19, 2020)

BUG FIXES:

* resource/digitalocean_volume: Fix validation on volume names ([#400](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/400)). Thanks to @Nevon!
* datasource/digitalocean_kubernetes_cluster: Don't error when `terraform:default-node-pool` tag is not found ([#399](https://github.com/terraform-providers/terraform-provider-digitalocean/issues/399)).
* resource/digitalocean_kubernetes_cluster: Fix local failure for Kubernetes interoperablity acceptance test when local kubeconfig file is present ([#402](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/402)).

## 1.15.0 (March 18, 2020)

FEATURES:

* **New Data Sources**: digitalocean_regions and digitalocean_region ([#380](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/380)) Thanks to @tdyas!
* **New Data Sources**: digitalocean_projects and digitalocean_project ([#391](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/391)) Thanks to @tdyas!
* **New Data Source**: digitalocean_images ([#394](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/394)) Thanks to @tdyas!
* **New Resource**: digitalocean_project_resources ([#396](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/396))
* resource/digitalocean_kubernetes_cluster, resource/digitalocean_kubernetes_node_pool: Add support for importing existing resources ([#365](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/365)). Thanks to @tdyas!

IMPROVEMENTS:

* datasource/digitalocean_droplet_snapshot, datasource/digitalocean_volume_snapshot resource/digitalocean_floating_ip, resource/digitalocean_floating_ip_assignment: Update for deprecated terraform-plugin-sdk helper/validation methods ([#376](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/376)).
* resource/digitalocean_kubernetes_cluster, resource/digitalocean_kubernetes_node_pool: Add support for Kubernetes node pool labels ([#379](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/379)). Thanks to @tdyas!
* internal/datalist: Add a generic filter/sort framework for datasources ([#385](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/385)). Thanks to @tdyas!
* resource/digitalocean_database_user: Add support for MySQL user authentication management ([#393](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/393)).

BUG FIXES:

* datasource/digitalocean_droplet: Validate only one of "id", "tag", or "name" are provided as an argument ([#375](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/375)). Thanks to @tdyas!
* resource/digitalocean_volume: Validate that volume names are lowercase and alphanumeric ([#386](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/386)). Thanks to @danrabinowitz!
* resource/digitalocean_database_cluster: "version" is now a required argument ([#382](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/382)).

## 1.14.0 (February 05, 2020)

IMPROVEMENTS:

* resource/digitalocean_kubernetes_cluster, resource/digitalocean_kubernetes_node_pool: Expose the Droplet IDs for individual nodes ([#366](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/366)). Thanks to @tdyas!
* datasource/digitalocean_droplet: Allow lookup by ID ([#366](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/366)). Thanks to @tdyas!

BUG FIXES:

* resource/digitalocean_project: Handle pagination for project resources ([#368](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/368)).

## 1.13.0 (January 27, 2020)

IMPROVEMENTS:

* resource/digitalocean_database_cluster: Add support for tags ([#353](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/353)). Thanks to @aqche!

BUG FIXES:

* provider: Mark API token as optional to support Spaces only usage and running `validate` without a token specified ([#356](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/356)).

## 1.12.0 (December 19, 2019)

FEATURES:

* **New Data Source**: `digitalocean_kubernetes_versions` ([#341](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/341))
* **New Resource:** : `digitalocean_database_firewall` ([#340](https://github.com/terraform-providers/terraform-provider-digitalocean/issues/340))

IMPROVEMENTS:

* resource/digitalocean_volume, datasource/digitalocean_volume: Add support for tags ([#336](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/336)). Thanks to @aqche!
* resource/digitalocean_volume_snapshot, datasource/digitalocean_volume_snapshot: Add support for tags ([#339](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/339)). Thanks to @aqche!
* resource/digitalocean_database_cluster: Add support for Redis eviction policies ([#342](https://github.com/terraform-providers/terraform-provider-digitalocean/issues/342)).
* resource/digitalocean_database_cluster: Add support for configuring SQL mode ([#347](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/347)).
* resource/digitalocean_certificate: Don't store full certificate data in state file ([#156](https://github.com/terraform-providers/terraform-provider-digitalocean/issues/156)).

## 1.11.0 (November 13, 2019)

FEATURES:

* **New Resource:** `digitalocean_database_connection_pool` ([#225](https://github.com/terraform-providers/terraform-provider-digitalocean/issues/225)) Thanks to @Photonios!

IMPROVEMENTS:

* resource/digitalocean_kubernetes_cluster: Add support for upgrading cluster versions ([#333](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/333)). Thanks to @aqche!

## 1.10.0 (October 31, 2019)

FEATURES:

* **New Resource:** `digitalocean_database_user` ([#328](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/328)) Thanks to @Permagate!
* **New Resource:** `digitalocean_database_db` ([#327](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/327)) Thanks to @Permagate!
* **New Data Source:** `digitalocean_sizes` ([#325](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/325)) hanks to @Permagate!
* **New Data Source:** `digitalocean_account` ([#324](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/324)) Thanks to @Permagate!

IMPROVEMENTS:

* resource/digitalocean_kubernetes_node_pool: Add support for node pool auto-scaling ([#307](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/307)). Thanks to @snormore!
* resource/digitalocean_space_bucket: Add support for configuring CORS rules ([#254](https://github.com/terraform-providers/terraform-provider-digitalocean/issues/254)). Thanks to @pohzipohzi!
* Migrate to using the Terraform Plugin SDK ([#316](https://github.com/terraform-providers/terraform-provider-digitalocean/issues/316)). Thanks to @stack72!
* website: Update all documentation to use v0.12 syntax ([#314](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/314)).

## 1.9.1 (October 09, 2019)

BUG FIXES:

* resource/digitalocean_kubernetes_cluster: Ensure `raw_config` is a valid kubeconfig file ([#315](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/315)).

## 1.9.0 (October 08, 2019)

IMPROVEMENTS:

* resource/digitalocean_database_cluster: Expose private connection details ([#299](https://github.com/terraform-providers/terraform-provider-digitalocean/issues/299)). Thanks to @clickyotomy!
* resource/digitalocean_database_cluster: Mark `uri`, `private_uri`, `password` as sensitive ([#298](https://github.com/terraform-providers/terraform-provider-digitalocean/issues/298)). Thanks to @clickyotomy!
* resource/digitalocean_database_replica: Expose private connection details ([#302](https://github.com/terraform-providers/terraform-provider-digitalocean/issues/302)).
* resource/digitalocean_kubernetes_cluster: Expose new `token` attribute for use in Kubernetes authentication ([#309](https://github.com/terraform-providers/terraform-provider-digitalocean/issues/309)). Thanks to @snormore!
* resource/digitalocean_kubernetes_cluster: Only fetch new Kubernetes credentials when expired ([#311](https://github.com/terraform-providers/terraform-provider-digitalocean/issues/311)). Thanks to @snormore!

## 1.8.0 (September 30, 2019)

IMPROVEMENTS:

* **New Data Source:** `digitalocean_database_replica` ([#224](https://github.com/terraform-providers/terraform-provider-digitalocean/issues/224)) Thanks to @Zyqsempai!
* resource/digitalocean_database_cluster: Add support for tags ([#253](https://github.com/terraform-providers/terraform-provider-digitalocean/issues/253)). Thanks to @Zyqsempai!
* resource/digitalocean_kubernetes_cluster: Mark the kube_config field as sensitive ([#289](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/289)). Thanks to @RoboticCheese!
* provider: Remove usage of `github.com/hashicorp/terraform/config` package ([#291](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/291)). Thanks to @appilon!
* datasource/digitalocean_droplet: Allow lookup by tag ([#290](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/290)). Thanks to @danramteke!
* provider: Remove usage of deprecated `terraform.VersionString()` ([#295](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/295)).

## 1.7.0 (August 27, 2019)

IMPROVEMENTS:

* resource/digitalocean_droplet: Expose created_at attribute ([#277](https://github.com/terraform-providers/terraform-provider-digitalocean/issues/277)). Thanks to @petems!
* datasource/digitalocean_droplet: Expose created_at attribute ([#277](https://github.com/terraform-providers/terraform-provider-digitalocean/issues/277)). Thanks to @petems!

BUG FIXES:

* resource/digitalocean_database_cluster: `version` should not be required ([#288](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/288)).
* resource/digitalocean_droplet: Verify Droplet is destroyed before removing from state ([#283](https://github.com/terraform-providers/terraform-provider-digitalocean/issues/283)).

## 1.6.0 (August 05, 2019)

IMPROVEMENTS:

* provider: Add the ability to set `DIGITALOCEAN_ACCESS_TOKEN` ([#260](https://github.com/terraform-providers/terraform-provider-digitalocean/issues/260)). Thanks to @stack72!
* resource/digitalocean_droplet: Support enabling and disabling backups on Droplets ([#266](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/267)).
* resource/digitalocean_database_cluster: Allow Databases to have URNs for use with Projects ([#270](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/270)). Thanks to @stack72!

BUG FIXES:

* Consistently protect against nil response in error handling ([#272](https://github.com/terraform-providers/terraform-provider-digitalocean/pull/273)).

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
