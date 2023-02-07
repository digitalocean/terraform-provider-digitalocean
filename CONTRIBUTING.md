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

All acceptance tests for a specific package can be run by setting the `PKG_NAME` environment variable. For example:

```sh
$ make testacc PKG_NAME=digitalocean/account
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

To release the provider:

1. Use
   [github-changelog-generator](https://github.com/digitalocean/github-changelog-generator)
   to list the changes since the last release and decide what kind of release
   you are doing (bugfix, feature or breaking).
1. Create a new commit that only contains updates to
   [CHANGELOG.md](CHANGELOG.md) listing the respective changes for the new
   version. Godo follows [semver](https://www.semver.org/) versioning
   semantics.
1. Once the CHANGELOG.md commit is merged, create a new tag on that commit with
   the new version that will be released (be sure to pull the latest from
   git).

   ```bash
   git tag -m "release $new_version" -a "$new_version"
   ```

1. Push the tag:

   ```bash
   git push "$origin" tag "$new_version"
   ```

This repository contains a GitHub Action configured to automatically build and
publish assets for release when a tag is pushed that matches the pattern `v*`
(ie. `v0.1.0`).

A [Goreleaser](https://goreleaser.com/) configuration is provided that produces
build artifacts matching the [layout required](https://www.terraform.io/docs/registry/providers/publishing.html#manually-preparing-a-release)
to publish the provider in the Terraform Registry.

Releases will appear as drafts. Once marked as published on the GitHub Releases page,
they will become available via the Terraform Registry.

Reviewing Pull Requests
-----------------------

Acceptance tests use the production API and create resources that incur costs.
Running the full suite of acceptance tests can also take quite a long time to run.
In order to prevent misuse, the acceptance tests for a pull request must be manually
triggered by a reviewer with write access to the repository.

To trigger a run of the acceptance tests for a PR, you may use the `/testacc` in a
comment. The `pkg` and `sha` arguments are required. This allows us to limit the
packages being tested for quick feedback and protect against timing attacks.
For example, to run the acceptance tests for the `droplet` package, the command
may look like:

    /testacc pkg=digitalocean/droplet sha=d358bd2418b4e30d7bdf2b98b4c151e357814c63

To run the entire suite, use `pkg=digitalocean`.

If multiple packages are to be tested, each command must be posted as a separate
comment. Only the first line of the comment is evaluated.

We leverage the [`peter-evans/slash-command-dispatch`](https://github.com/peter-evans/slash-command-dispatch)
GitHub Action for the command processing.
