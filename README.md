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

See [CONTRIBUTING.md](./CONTRIBUTING.md) for information about contributing to this project.
