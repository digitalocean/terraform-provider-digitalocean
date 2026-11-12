# AGENTS.md

Guidance for coding agents working in this repository.

## What this is

The DigitalOcean Terraform provider (`digitalocean/digitalocean` on the Terraform Registry). It is written in Go on `terraform-plugin-sdk/v2` and calls the DigitalOcean API through [godo](https://github.com/digitalocean/godo).

- User docs: https://docs.digitalocean.com/reference/terraform/ and https://registry.terraform.io/providers/digitalocean/digitalocean/latest/docs
- Developer guide: [CONTRIBUTING.md](./CONTRIBUTING.md)

## Layout

- `main.go`: plugin entry point.
- `digitalocean/provider.go`: provider schema plus `ResourcesMap` and `DataSourcesMap`. Every resource and data source is registered here.
- `digitalocean/<service>/`: one package per service (`droplet`, `domain`, `kubernetes`, `vpc`, ...).
- `digitalocean/acceptance/`: shared acceptance test helpers (`TestAccPreCheck`, `TestAccProviderFactories`, `RandomTestName`).
- `digitalocean/config/`: provider configuration and API clients.
- `digitalocean/sweep/`: sweeper entry point; each service package registers its sweepers in `sweep.go`.
- `docs/resources/`, `docs/data-sources/`: Registry documentation, one Markdown file per resource or data source.
- `examples/`: sample Terraform configurations.
- `vendor/`: vendored Go modules. Update with `make vendor` (or `make upgrade_godo` for godo).

## Commands

- `make build`: runs the gofmt check and `go install`.
- `make test`: unit tests. Safe to run; acceptance tests skip without `TF_ACC`.
- `make vet`, `make lint` (golangci-lint), `make fmtcheck`, `make fmt`.
- `make terrafmt-check` / `make terrafmt`: format HCL embedded in tests and docs.

CI (`.github/workflows/test.yml`) runs `make vet`, `make lint`, `make terrafmt-check`, and `make test`.

## Acceptance tests

Acceptance tests (`TestAcc*`) create real, billable resources in a DigitalOcean account. They run only when `TF_ACC=1` and `DIGITALOCEAN_TOKEN` are set, usually through `make testacc`. Do not run `make testacc`, set `TF_ACC`, or run `make sweep` (which destroys infrastructure) unless the user explicitly asks. When asked, scope the run to one test or package:

```sh
make testacc TESTARGS='-run=TestAccDigitalOceanDomain_Basic' PKG_NAME=digitalocean/domain
```

On pull requests, maintainers trigger acceptance tests with a `/testacc pkg=... sha=...` comment.

## Conventions

Follow `digitalocean/domain/` as a reference:

- Resources live in `resource_<name>.go` and export `ResourceDigitalOcean<Name>()`; data sources live in `datasource_<name>.go` and export `DataSourceDigitalOcean<Name>()`.
- CRUD functions use the context-aware signatures and return `diag.Diagnostics`. Get the godo client with `meta.(*config.CombinedConfig).GodoClient()`.
- Tests sit beside the code: `resource_<name>_test.go`, `datasource_<name>_test.go`, `import_<name>_test.go`. Name test resources with `acceptance.RandomTestName()` so they get the `tf-acc-test-` prefix that sweepers clean up.
- New service packages must be imported in `digitalocean/sweep/sweep_test.go` so their sweepers run.
- When you add or change a resource or data source, update its page in `docs/resources/` or `docs/data-sources/`. The file name drops the `digitalocean_` prefix, and the front matter sets `page_title` and `subcategory`.
- Run `make fmt` and `make terrafmt` before committing.
