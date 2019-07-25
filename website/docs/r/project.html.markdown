---
layout: "digitalocean"
page_title: "DigitalOcean: digitalocean_project"
sidebar_current: "docs-do-resource-project"
description: |-
  Provides a DigitalOcean Project resource.
---

# digitalocean\_project

Provides a DigitalOcean Project resource.

Projects allow you to organize your resources into groups that fit the way you work.
You can group resources (like Droplets, Spaces, Load Balancers, domains, and Floating IPs)
in ways that align with the applications you host on DigitalOcean.

The following resources can be associated with a project:

* Droplet
* Load Balancer
* Domain
* Volume
* Floating IP
* Spaces Bucket

**Note:** A Terrafrom managed project cannot be set as a default project.

## Example Usage

The following example demonstrates the creation of an empty project:

```hcl
resource "digitalocean_project" "playground" {
  name        = "playground"
  description = "A project to represent development resources."
  purpose     = "Web Application"
  environment = "Development"
}
```

The following example demonstrates the creation of a project with a Droplet resource:

```hcl
resource "digitalocean_droplet" "foobar" {
  name   = "example"
  size   = "512mb"
  image  = "centos-7-x64"
  region = "nyc3"
}

resource "digitalocean_project" "playground" {
  name        = "playground"
  description = "A project to represent development resources."
  purpose     = "Web Application"
  environment = "Development"
  resources   = [digitalocean_droplet.foobar.urn]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Project
* `description` - (Optional) the description of the project
* `purpose` - (Optional) the purpose of the project, (Default: "Web Application")
* `environment` - (Optional) the environment of the project's resources. The possible values are: `Development`, `Staging`, `Production`)
* `resources` - a list of uniform resource names (URNs) for the resources associated with the project

## Attributes Reference

The following attributes are exported:

* `id` - The id of the project
* `owner_uuid` - the unique universal identifier of the project owner.
* `owner_id` - the id of the project owner.
* `created_at` - the date and time when the project was created, (ISO8601)
* `updated_at` - the date and time when the project was last updated, (ISO8601)
