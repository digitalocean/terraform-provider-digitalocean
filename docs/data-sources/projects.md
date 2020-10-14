---
page_title: "DigitalOcean: digitalocean_projects"
---

# digitalocean_projects

Retrieve information about all DigitalOcean projects associated with an account, with
the ability to filter and sort the results. If no filters are specified, all projects
will be returned.

Note: You can use the [`digitalocean_project`](project) data source to
obtain metadata about a single project if you already know the `id` to retrieve or the unique
`name` of the project.

## Example Usage

Use the `filter` block with a `key` string and `values` list to filter projects.

For example to find all staging environment projects:

```hcl
data "digitalocean_projects" "staging" {
  filter {
    key    = "environment"
    values = ["Staging"]
  }
}
```

You can filter on multiple fields and sort the results as well:

```hcl
data "digitalocean_projects" "non-default-production" {
  filter {
    key    = "environment"
    values = ["Production"]
  }
  filter {
    key    = "is_default"
    values = ["false"]
  }
  sort {
    key       = "name"
    direction = "asc"
  }
}
```

## Argument Reference

* `filter` - (Optional) Filter the results.
  The `filter` block is documented below.

* `sort` - (Optional) Sort the results.
  The `sort` block is documented below.

`filter` supports the following arguments:

* `key` - (Required) Filter the projects by this key. This may be one of `name`,
  `purpose`, `description`, `environment`, or `is_default`.
  
* `values` - (Required) A list of values to match against the `key` field. Only retrieves projects
  where the `key` field takes on one or more of the values provided here.

* `match_by` - (Optional) One of `exact` (default), `re`, or `substring`. For string-typed fields, specify `re` to
  match by using the `values` as regular expressions, or specify `substring` to match by treating the `values` as
  substrings to find within the string field.
  
* `all` - (Optional) Set to `true` to require that a field match all of the `values` instead of just one or more of
  them. This is useful when matching against multi-valued fields such as lists or sets where you want to ensure
  that all of the `values` are present in the list or set.

`sort` supports the following arguments:

* `key` - (Required) Sort the projects by this key. This may be one of `name`,
  `purpose`, `description`, or `environment`.
* `direction` - (Required) The sort direction. This may be either `asc` or `desc`.

## Attributes Reference

* `projects` - A set of projects satisfying any `filter` and `sort` criteria. Each project has
  the following attributes:  
  - `id` - The ID of the project
  - `name` - The name of the project
  - `description` - The description of the project
  - `purpose` -  The purpose of the project (Default: "Web Application")
  - `environment` - The environment of the project's resources. The possible values are: `Development`, `Staging`, `Production`.
  - `resources` - A set of uniform resource names (URNs) for the resources associated with the project
  - `owner_uuid` - The unique universal identifier of the project owner
  - `owner_id` - The ID of the project owner
  - `created_at` - The date and time when the project was created, (ISO8601)
  - `updated_at` - The date and time when the project was last updated, (ISO8601)
