# Global Load Balancer: two-phase custom certificate rotation

This example shows the **safe immediate process** for rotating a **custom** TLS certificate on a **Global** load balancer when a single `terraform apply` can hit API errors deleting the old certificate while it is still referenced.

It is a **structural template** only: you must supply real PEM material and merge the `domains` change into your **existing** Global Load Balancer configuration.

## Prerequisites

- A working **`digitalocean_loadbalancer`** with `type = "GLOBAL"` and `domains` using a custom certificate today.
- A **new** certificate `name` that does not collide with any existing certificate name in the account.

## Phase 1 — switch GLB to new cert (keep old cert in Terraform)

1. Copy `main.tf` and fill variables (or use `.tfvars`).
2. Ensure **both** `digitalocean_certificate.old` and `digitalocean_certificate.new` exist in configuration.
3. In your **existing** `digitalocean_loadbalancer` (GLOBAL), set `domains` to use the new cert, for example:

```hcl
  domains {
    name               = "your.hostname.example"
    is_managed         = false
    certificate_name = digitalocean_certificate.new.name
  }
```

Or use the value from `terraform output attach_this_certificate_name_to_glb_domains` if wiring across modules.

4. Run `terraform plan` — you should see **create** (new cert) and **update** (GLB), **no destroy** of the old cert.
5. Run `terraform apply`.

## Phase 2 — delete old cert only

1. Remove the `digitalocean_certificate.old` resource block (and any references) from your configuration.
2. Run `terraform plan` — you should see **destroy** for the old cert only.
3. Run `terraform apply`.

## If you start from a single certificate resource

Use `terraform state mv` to split one resource into `old` / `new` addresses, or add a second resource and import — pair with support on the first migration if needed. See `docs/design/glb-certificate-terraform-rotation.md` section **Immediate safe process**.
