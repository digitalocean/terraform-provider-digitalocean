# DigitalOcean Terraform Provider

[![Terraform Registry](https://img.shields.io/badge/terraform-registry-blue.svg)](https://registry.terraform.io/providers/digitalocean/digitalocean/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/digitalocean/terraform-provider-digitalocean)](https://goreportcard.com/report/github.com/digitalocean/terraform-provider-digitalocean)

The official Terraform provider for DigitalOcean.

📖 **Documentation**: https://registry.terraform.io/providers/digitalocean/digitalocean/latest/docs

## Table of Contents

- [Requirements](#requirements)
- [Building The Provider](#building-the-provider)
- [Using the Provider](#using-the-provider)
- [Developing the Provider](#developing-the-provider)

## Requirements
------------

-	[Terraform](https://developer.hashicorp.com/terraform/install) 0.10+
-	[Go](https://go.dev/doc/install) 1.14+ (to build the provider plugin)

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

See [CONTRIBUTING.md](./CONTRIBUTING.md) for information about contributing to this project.
