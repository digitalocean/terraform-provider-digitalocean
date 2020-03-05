---
layout: "digitalocean"
page_title: "DigitalOcean: digitalocean_project"
sidebar_current: "docs-do-datasource-project"
description: |-
  Get information on a DigitalOcean project.
---

# digitalocean_project

Get information on a single DigitalOcean project. If neither the `id` nor `name` attributes is provided,
then this datasource returns the default project.

## Example Usage

```hcl
data "digitalocean_project" "default" {
} 

data "digitalocean_project" "staging" {
  name = "My Staging Project"
}
```

## Argument Reference

* `id` - (Optional) The id of the project to retrieve
* `name` - (Optional) The name of the project to retrieve. The datasource will raise an error if more than
  one project has the provided name.

## Attributes Reference

* `description` - the description of the project
* `purpose` -  the purpose of the project, (Default: "Web Application")
* `environment` - the environment of the project's resources. The possible values are: `Development`, `Staging`, `Production`.
* `resources` - a list of uniform resource names (URNs) for the resources associated with the project
* `owner_uuid` - the unique universal identifier of the project owner.
* `owner_id` - the id of the project owner.
* `created_at` - the date and time when the project was created, (ISO8601)
* `updated_at` - the date and time when the project was last updated, (ISO8601)
