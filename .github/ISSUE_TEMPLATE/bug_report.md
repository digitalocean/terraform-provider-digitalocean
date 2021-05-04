---
name: Bug report
about: Create a report to help us improve
title: ''
labels: bug
assignees: ''

---

Thank you for opening an issue. Please note that we try to keep the issue
tracker reserved for bug reports and feature requests. For general usage
questions, please see:
https://github.com/digitalocean/terraform-provider-digitalocean/discussions

**NOTE: Before submitting a bug**

There are cases where the provider receives HTTP Service Error (500 level HTTP
statuses) responses from the DigitalOcean API. There are some cases where the
provider might handle these and retry. If the problem persists, it's best to
contact [DigitalOcean support](https://cloudsupport.digitalocean.com/) 

# Bug Report

Include as much of the following details with your bug report:

> everything above this can be ommitted
---

## Describe the bug
A clear and concise description of what the bug is.

### Affected Resource(s)
Please list the resources, for example:
- digitalocean_droplet
- digitalocean_kubernetes_cluster

If this issue appears to affect multiple resources, it may be an issue with
Terraform's core, so please mention this.

### Expected Behavior
What should have happened?

### Actual Behavior
What actually happened?

### Steps to Reproduce
Please list the steps required to reproduce the issue, for example:
1. `terraform apply`

**Terraform Configuration Files**
```hcl
# Copy-paste your Terraform configurations here - for large Terraform configs,
# please use a service like Dropbox and share a link to the ZIP file. For
# security, you can also encrypt the files using our GPG public key.
```

**Expected behavior**
Run `terraform -v` to show the version. If you are not running the latest
version of Terraform, please upgrade because your issue may have already been
fixed.

**Debug Output**
Please provide a link to a GitHub Gist containing the complete debug output:
https://www.terraform.io/docs/internals/debugging.html. Please do NOT paste the
debug output in the issue; just paste a link to the Gist.

**Panic Output**
If Terraform produced a panic, please provide a link to a GitHub Gist
containing the output of the `crash.log`.

## Additional context
Add any other context about the problem here.

**Important Factoids**
Droplets use custom images or kernels.

**References**
Include links to other GitHub issues (open or closed) or Pull Requests that
relate to this issue.

