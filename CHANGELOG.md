# 2.34.1

BUG FIXES:

- `digitalocean_cdn`: handle 'needs-cloudflare-cert' case in read func (#1095). - @andrewsomething
- `digitalocean_database_cluster`: ignore seconds in maintenance_window.hour (#1094). - @andrewsomething
- build(deps): bump golang.org/x/crypto from 0.14.0 to 0.17.0 (#1096). - @dependabot[bot]

# 2.34.0

IMPROVEMENTS:

- `digitalocean_database_user`: Support updating ACL settings (#1090). - @dweinshenker

BUG FIXES:

- `digitalocean_cdn`: Add Support for "needs-cloudflare-cert" (#1089). - @danaelhe
- `digitalocean_spaces_bucket`: blr1 is a supported region (#1085). - @andrewsomething
- `digitalocean_database_kafka_topic`: Kafka topic + user ACL management doc fixes (#1082). - @dweinshenker

# 2.33.0

IMPROVEMENTS:

- #1073 - @T-jegou - Add `digitalocean_database_connection_pool` datasource

BUG FIXES:

- #1078 - @nemcikjan - fix: added missing option to set port on health_check
- #1076 - @dweinshenker - Remove unclean_leader_election_enable for kafka topic configuration
- #1080 - @danaelhe - Apps: Reference Port in expandAppHealthCheck and flattenAppHealthCheck
- #1074 - @T-jegou - Fixing Case Handling for Volume Resource


# 2.32.0

IMPROVEMENTS:

- `digitalocean_app`: Support `features` in App spec (#1066). - @T-jegou
- `digitalocean_database_user`: Add support for Kafka Topic User ACL management (#1056). - @dweinshenker
- `digitalocean_kubernetes_cluster`: Support enabling HA post-create (#1058). - @andrewsomething

BUG FIXES:

- `digitalocean_loadbalancer`: ignore 404 on delete (#1067). - @andrewsomething
- `digitalocean_database_mysql_config`: Use GetOkExists for bools (#1063). - @andrewsomething
- `digitalocean_kubernetes_cluster`: Handle error from GetCredentials and protect against panic (#1064). - @andrewsomething

MISC:

- `provider`: Bump godo to v1.105.1 (#1071). - @andrewsomething
- `provider`: bump google.golang.org/grpc from 1.53.0 to 1.56.3 (#1057). - @dependabot[bot]

# 2.31.0

FEATURES:

- **New Resource:** `digitalocean_database_kafka_topic` (#1052) - @dweinshenker
- **New Resource:** `digitalocean_database_mysql_config` (#1051) - @kallydev
- **New Resource:** `digitalocean_database_redis_config` (#1037) - @andrewsomething

IMPROVEMENTS:

- `digitalocean_database_cluster`: Add support for Scalable Storage (#1054). - @dweinshenker
- `digitalocean_app`: Add support for ingress for component routing, rewrites, and redirects (#1053). - @andrewsomething
- `digitalocean_loadbalancer`: Add support type param (#1023). - @asaha2

BUG FIXES:

- `digitalocean_loadbalancer`: no region field needed for global lb type (#1046). - @apinonformoso
- `digitalocean_loadbalancer`: Parse nil region for global lb (#1043). - @asaha2
- `digitalocean_app`: Rework deployment logic (#1048). - @davidsbond
- `digitalocean_spaces_bucket`: set force_destroy false on import (#1041). - @andrewsomething

MISC:

- `build(deps)`: bump golang.org/x/net from 0.14.0 to 0.17.0 (#1050). - @dependabot[bot]
- `docs`: Clarify Database Docs for Referencing DB Replicas (#1045). - @danaelhe
- `testing`: Use terrafmt on docs directory (#1036). - @andrewsomething 
- `docs`: Update Droplet example (#1035). - @danaelhe

# 2.30.0

FEATURES:

- **New Resource:** `digitalocean_spaces_bucket_cors_configuration` (#1021) - @danaelhe 

IMPROVEMENTS:

- `provider`: Enable retries for requests that fail with a 429 or 500-level error by default (#1016). - @danaelhe

BUG FIXES:

- `digitalocean_database_user`: Prevent creating multiple users for the same cluster in parallel (#1027). - @andrewsomething
- `digitalocean_database_user`:  Remove unneeded GET request post-create (#1028). - @andrewsomething

MISC:

- `docs`: Make it clear that volume name has to start with a letter (#1024). - @ahasna
- `docs`: Update Postgres version in example (#1014). - @danaelhe
- `provider`: Bump Go version to v1.21.0 (#1025). - @andrewsomething 
- `provider`: Update godo to v1.102.1 (#1020). - @danaelhe
- `provider`: Update godo dependency to v1.102.0 (#1018). - @danaelhe
- `provider`: Update godo dependency to v1.101.0 (#1017.) - @danaelhe

# 2.29.0

FEATURES:

- **New Data Source:** `digitalocean_database_user` (#989). - @lezgomatt

IMPROVEMENTS:

- `digitalocean_kubernetes_cluster`: Add destroy_all_associated_resources option (#1007). - @andrewsomething

BUG FIXES:

- `digitalocean_spaces_bucket`: Update `retryOnAwsCode` to five minutes (#999). - @danaelhe 

MISC:

- `docs`: Note how to get `id` for record import (#1004) - @nimboya
- `provider`: Bump Go version to 1.20.x (#975). - @andrewsomething
- `testing`: Update Postgres versions in acceptance tests (#1002). - @andrewsomething
- `provider`: build(deps): bump google.golang.org/grpc from 1.51.0 to 1.53.0 (#1003). - @dependabot[bot]

# 2.28.1

BUG FIXES:

- `digitalocean_database_cluster`: Fix custom create timeouts (#987). - @andrewsomething
- `digitalocean_droplet`: Prevent inconsistent plan when enabling IPv6 (#982). - @andrewsomething
- `digitalocean_custom_image`: use custom create timeout (#985). - @andrewsomething

# 2.28.0

IMPROVEMENTS:

- `provider`: Add godo's rate limiter configuration & retryable http client (#967). - @DanielHLelis
- `digitalocean_kubernetes_cluster`: Support container registry integration (#963). - @mohsenSy
- `digitalocean_database_replica`: Add support for resizing replicas (#977). - @andrewsomething
- `digitalocean_database_cluster`: Add backup-restore functionality to db create (#970). - @danaelhe

BUG FIXES:

- `digitalocean_record`: Handle pagination in data source (#979). - @andrewsomething
- `digitalocean_kubernetes_cluster`: Require importing additional node pools manually (#976). - @andrewsomething
- `digitalocean_database_replica`: Add uuid to data source schema (#969). - @andrewsomething

MISC:

- `docs`: Fix inconsistencies in `digitalocean_uptime_alert` documentation #972 - @nicwortel
- `docs`: Use correct links in uptime docs. #973 - @andrewsomething
- `provider`: Update Terraform SDK to v2.26.1. #975 - @andrewsomething

# 2.27.1

BUG FIXES:

- `digitalocean_database_replica`: Set UUID on read to resolve import issue (#964). - @andrewsomething

MISC:

- dependencies: bump golang.org/x/net (#957). - @dependabot
- dependencies: bump golang.org/x/crypto (#960). - @dependabot

# 2.27.0

IMPROVEMENTS:

- `digitalocean_database_cluster`: Support project assignment (#955). - @andrewsomething

BUG FIXES:

- `digitalocean_custom_image`: use correct pending statuses for custom images (#931). - @rsmitty

DOCS:

-  `digitalocean_app`: Fix typo in resource digitalocean_app (#961). - @tobiasehlert

MISC:

- `provider`: Package reorganization (#927). - @andrewsomething
- `testing`: Use comment trigger to run acceptance tests for PRs. (#929). - @andrewsomething
- `testing`: Fix formatting of Terraform configs in tests and enforce in PRs using terrafmt (#932). - @andrewsomething
- `testing`: droplet: Fix acceptance testing (#949). - @andrewsomething
- `testing`: certificates: Add retry on delete (#951). - @andrewsomething
- `testing`: cdn: Add test sweeper and retry with backoff (#947). - @andrewsomething
- `testing`: Add sweeper and use consistent naming for all Spaces buckets in tests (#945). - @andrewsomething
- `testing`: Add sweeper for uptime and monitoring alerts (#944). - @andrewsomething
- `testing`: Add sweeper for projects and add retry for project deletion (#943). - @andrewsomething
- `testing`: Add sweeper for VPCs (#942). - @andrewsomething
- `testing`: Add sweeper for custom images and fix acceptance tests (#941). - @andrewsomething
- `testing`: Use consistent naming for all volumes created in tests (#939). - @andrewsomething
- `testing`: Use consistent naming for all snapshots created in tests (#938). - @andrewsomething
- `testing`: Use consistent naming for all load balancers created in tests (#937). - @andrewsomething
- `testing`: Use consistent naming for all firewalls created in tests (#935). - @andrewsomething
- `testing`: Add sweeper for SSH keys (#940). - @andrewsomething
- `testing`: Use consistent naming for all certs created in tests (#934). - @andrewsomething
- `testing`: Use consistent naming for all Droplets created in tests (#933). - @andrewsomething
- `testing`: Remove unused const to fix linting (#930). - @andrewsomething
- `testing`: Fix flaky database acceptance tests (#953). - @andrewsomething
- Remove .go-version and add to .gitignore (#958). - @ChiefMateStarbuck

# 2.26.0

IMPROVEMENTS:

- `database replica`: Expose Database Replica ID (#921) - @danaelhe
- `uptime`: Add Uptime Checks and Alerts Support (#919) - @danaelhe
- `databases`: Support upgrading the database version (#918) - @scotchneat
- `loadbalancers`: Add firewall support for Loadbalancers (#911) - @jrolheiser
- `loadbalancers`: Loadbalancers support http alerts metrics (#903) - @StephenVarela

MISC:

- `docs`: `routes` documentation in `app.md` matches `app_spec.go` (#915) - @olaven
- `testing`: Find previous K8s release dynamically. (#916) - @andrewsomething
- `docs`: Fix typo in README (#920) - @mbardelmeijer
- `docs`: Add releasing notes & missing changelog entries (#922) - @scotchneat

# 2.25.2

IMPROVEMENTS:

- `database_replica`: add retry on db replica create (#907) - @DMW2151

# 2.25.1

IMPROVEMENTS:

- `monitoring`: Support HTTP idle timeout & Project ID (#897) - @StephenVarela

# 2.24.0

IMPROVEMENTS:

- `spaces`: add endpoint attribute to bucket (#886)- @selborsolrac
- `monitor_alert_resource`: Update Monitor Alert resource with new DBAAS public alert types (#893) - @dweinshenker
- `spaces`: Add new DC to spaces (#899) - @mandalae
- `loadbalancers`: load balancers: add HTTP/3 as an entry protocol (#895) - @anitgandhi

MISC:

- `docs`: Fix reference in documentation of project_resources (#890) - @Lavode

# 2.23.0 (September 27, 2022)

IMPROVEMENTS:

- `digitalocean_app`: Support deploy on push to DOCR ([#883](https://github.com/digitalocean/terraform-provider-digitalocean/pull/883)). - @andrewsomething
- `digitalocean_droplet`: Region is no longer a required value ([#879](https://github.com/digitalocean/terraform-provider-digitalocean/pull/879)). - @andrewsomething

BUG FIXES:

- `digitalocean_record`: Add SOA as possible record type ([#882](https://github.com/digitalocean/terraform-provider-digitalocean/pull/882)). - @Nosmoht

MISC:

- Upgrade to Go 1.19  ([#884](https://github.com/digitalocean/terraform-provider-digitalocean/pull/884)). - @andrewsomething

# 2.22.3 (September 12, 2022)

BUG FIXES:

- `digitalocean_droplet`: Fix configurable timeouts for Droplet creates ([#867](https://github.com/digitalocean/terraform-provider-digitalocean/pull/867)). - @andrewsomething

# 2.22.2 (August 31, 2022)

IMPROVEMENTS:

- `digitalocean_database_connection_pool`: make user optional in db connection pool, update acc tests ([#868](https://github.com/digitalocean/terraform-provider-digitalocean/pull/868)) - @DMW2151

MISC:

- `digitalocean_database_cluster`: Suppress diffs on forced Redis version upgrades ([#873](https://github.com/digitalocean/terraform-provider-digitalocean/pull/873)) - @scotchneat
- `docs`: fix app spec link([#871](https://github.com/digitalocean/terraform-provider-digitalocean/pull/871)) - @jkpe

# 2.22.1 (August 16, 2022)

BUG FIXES:

- `digitalocean_app`: Limit the number of deployments listed when polling ([#865](https://github.com/digitalocean/terraform-provider-digitalocean/pull/865)). - @andrewsomething

MISC:

- release workflow: switch to `crazy-max/ghaction-import-gpg@v5.0.0` ([#863](https://github.com/digitalocean/terraform-provider-digitalocean/pull/863)). - @andrewsomething

## 2.22.0 (August 15, 2022)

IMPROVEMENTS:

- `digitalocean_project`: Make `is_default` configurable ([#860](https://github.com/digitalocean/terraform-provider-digitalocean/pull/860)). - @danaelhe
- `digitalocean_droplet`: Configurable timeouts for Droplet create operations ([#839](https://github.com/digitalocean/terraform-provider-digitalocean/pull/839)). - @ngharrington
- `digitalocean_app`: Add computed URN attribute ([#854](https://github.com/digitalocean/terraform-provider-digitalocean/pull/854)). - @andrewsomething

BUG FIXES:

- `digitalocean_app`: Only warn on read if there is no active deployment ([#843](https://github.com/digitalocean/terraform-provider-digitalocean/pull/843)). - @andrewsomething
- `digitalocean_database_firewall`: Remove firewall rule from state if missing ([#840](https://github.com/digitalocean/terraform-provider-digitalocean/pull/840)). - @liamjcooper

MISC:

- chore: Fix incorrect heading in bug template ([#859](https://github.com/digitalocean/terraform-provider-digitalocean/pull/859)). - @artis3n
- testing: Use supported OS for Droplet image slug. ([#855](https://github.com/digitalocean/terraform-provider-digitalocean/pull/855)). - @andrewsomething
- testing: Introduce golangci-lint in GitHub workflows. ([#755](https://github.com/digitalocean/terraform-provider-digitalocean/pull/755)). - @atombrella
- docs: Add more examples on using Droplet snapshots ([#846](https://github.com/digitalocean/terraform-provider-digitalocean/pull/846)). - @mkjmdski

## 2.21.0 (June 16, 2022)

FEATURES:

- **New Resource:** `digitalocean_reserved_ip`  ([#830](https://github.com/digitalocean/terraform-provider-digitalocean/pull/830)). - @andrewsomething
- **New Resource:** `digitalocean_reserved_ip_assignment` ([#830](https://github.com/digitalocean/terraform-provider-digitalocean/pull/830)). - @andrewsomething
- **New Data Source:** `digitalocean_reserved_ip` ([#830](https://github.com/digitalocean/terraform-provider-digitalocean/pull/830)). - @andrewsomething

MISC:

- examples: Change k8s example to use ingress v1 ([#831](https://github.com/digitalocean/terraform-provider-digitalocean/pull/837)). - @jacobgreenleaf

## 2.20.0 (May 25, 2022)

IMPROVEMENTS:

- `digitalocean_app`: Support functions components ([#831](https://github.com/digitalocean/terraform-provider-digitalocean/pull/831)). - @andrewsomething
- `digitalocean_monitor_alert`: Support load balancer alert types ([#822](https://github.com/digitalocean/terraform-provider-digitalocean/pull/822)). - @andrewsomething
- `digitalocean_loadbalancer`: support udp as a target and entry protocol ([#789](https://github.com/digitalocean/terraform-provider-digitalocean/pull/789)). - @dikshant

BUG FIXES:

- `digitalocean_kubernetes_cluster`: Always perform upgrade check ([#823](https://github.com/digitalocean/terraform-provider-digitalocean/issues/823)). - @macno

MISC:

- docs: Document values attribute in ssh_keys data source ([#832](https://github.com/digitalocean/terraform-provider-digitalocean/pull/832)). - @andrewsomething
- docs: Note limitations on importing MongoDB users ([#821](https://github.com/digitalocean/terraform-provider-digitalocean/pull/821)). - @andrewsomething
- docs: Fix k8's node_pool.tags description ([#816](https://github.com/digitalocean/terraform-provider-digitalocean/pull/816)). - @danaelhe
- testing: Fix k8s versions in acceptance tests ([#826](https://github.com/digitalocean/terraform-provider-digitalocean/pull/826)). - @gizero
- provider: Build with go 1.18 ([#813](https://github.com/digitalocean/terraform-provider-digitalocean/pull/813)). - @ChiefMateStarbuck

## 2.19.0 (March 28, 2022)

IMPROVEMENTS:

- `digitalocean_container_registry`: Support providing custom region ([#804](https://github.com/digitalocean/terraform-provider-digitalocean/pull/804)). - @andrewsomething

## 2.18.0 (March 8, 2022)

FEATURES:

- **New Resource:** `digitalocean_spaces_bucket_policy` ([#800](https://github.com/digitalocean/terraform-provider-digitalocean/pull/800)) - @pavelkovar

IMPROVEMENTS:

- `digitalocean_app`: Implement support for App Platform `log_destinations` ([#798](https://github.com/digitalocean/terraform-provider-digitalocean/pull/798)). - @jbrunton
- `digitalocean_app`: Implement support for configuring alert policies ([#797](https://github.com/digitalocean/terraform-provider-digitalocean/pull/797)). - @andrewsomething

BUG FIXES:

- `digitalocean_project`: Environment is optional, don't set default ([#788](https://github.com/digitalocean/terraform-provider-digitalocean/pull/788)). - @andrewsomething
- `digitalocean_droplet`: Handle optional boolean `droplet_agent` ([#785](https://github.com/digitalocean/terraform-provider-digitalocean/pull/785)). - @Kidsan

MISC:

- docs: Add default ttl to `digitalocean_record` ([#791](https://github.com/digitalocean/terraform-provider-digitalocean/pull/791)). - @unixlab

## 2.17.1 (January 28, 2022)

IMPROVEMENTS:

- `digitalocean_app`: Allow using `MONGODB` as a database engine option ([#783](https://github.com/digitalocean/terraform-provider-digitalocean/pull/783)). - @cnunciato
- doc: Update docs for `digitalocean_monitor_alert` to highlight it is currently Droplet-only ([#780](https://github.com/digitalocean/terraform-provider-digitalocean/pull/780)). - @andrewsomething

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
