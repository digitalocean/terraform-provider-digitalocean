DigitalOcean Terraform Provider
==================

- Documentation: https://registry.terraform.io/providers/digitalocean/digitalocean/latest/docs

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.10.x
-	[Go](https://golang.org/doc/install) 1.14 (to build the provider plugin)

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/digitalocean/terraform-provider-digitalocean`

```sh
$ mkdir -p $GOPATH/src/github.com/digitalocean; cd $GOPATH/src/github.com/digitalocean
$ git clone git@github.com:digitalocean/terraform-provider-digitalocean
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/digitalocean/terraform-provider-digitalocean
$ make build
```

Using the provider
----------------------

See the [DigitalOcean Provider documentation](https://registry.terraform.io/providers/digitalocean/digitalocean/latest/docs) to get started using the DigitalOcean provider.

Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.11+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-digitalocean
...
```

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```

In order to run a specific acceptance test, use the `TESTARGS` environment variable. For example, the following command will run `TestAccDigitalOceanDomain_Basic` acceptance test only:

```sh
$ make testacc TESTARGS='-run=TestAccDigitalOceanDomain_Basic'
```

In order to check changes you made locally to the provider, you can use the binary you just compiled by adding the following
to your `~/.terraformrc` file. This is valid for Terraform 0.14+. Please see
[Terraform's documentation](https://www.terraform.io/docs/cli/config/config-file.html#development-overrides-for-provider-developers)
for more details.

```
provider_installation {

  # Use /home/developer/go/bin as an overridden package directory
  # for the digitalocean/digitalocean provider. This disables the version and checksum
  # verifications for this provider and forces Terraform to look for the
  # digitalocean provider plugin in the given directory.
  dev_overrides {
    "digitalocean/digitalocean" = "/home/developer/go/bin"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```

For information about writing acceptance tests, see the main Terraform [contributing guide](https://github.com/hashicorp/terraform/blob/master/.github/CONTRIBUTING.md#writing-acceptance-tests).

Releasing the Provider
----------------------

This repository contains a GitHub Action configured to automatically build and
publish assets for release when a tag is pushed that matches the pattern `v*`
(ie. `v0.1.0`).

A [Gorelaser](https://goreleaser.com/) configuration is provided that produces
build artifacts matching the [layout required](https://www.terraform.io/docs/registry/providers/publishing.html#manually-preparing-a-release)
to publish the provider in the Terraform Registry.

Releases will appear as drafts. Once marked as published on the GitHub Releases page,
they will become available via the Terraform Registry.
