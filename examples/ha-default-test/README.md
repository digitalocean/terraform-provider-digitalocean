# Using Your Locally Modified Terraform Provider

## Step 1: Build the provider

From the provider repo root:

```bash
cd /Users/pyadagiri/dev/digitalocean/terraform-provider-digitalocean
make build
```

This installs the binary to `~/go/bin/terraform-provider-digitalocean`.

## Step 2: Configure Terraform to use your local build

Add this to `~/.terraformrc` (create the file if it doesn't exist):

```hcl
provider_installation {
  dev_overrides {
    "digitalocean/digitalocean" = "/Users/pyadagiri/go/bin"
  }
  direct {}
}
```

Terraform will look for the provider binary in that directory. You'll see a warning when using dev overrides—that's expected.

## Step 3: Set your API token

```bash
export DIGITALOCEAN_TOKEN="your_token_here"
# or: export DIGITALOCEAN_ACCESS_TOKEN="your_token_here"
```

Or add `token = "..."` to the provider block in main.tf.

## Step 4: Run Terraform

```bash
cd /Users/pyadagiri/dev/digitalocean/terraform-provider-digitalocean/examples/ha-default-test
export TF_LOG=TRACE # enable debug logs
terraform init
terraform plan
terraform apply
```

## Step 5: Check the result

```bash
terraform output ha
# Should show true for DOKS 1.36+
```

## Step 6: Clean up

```bash
terraform destroy
```

## Rebuilding after changes

After modifying the provider code:

```bash
cd /Users/pyadagiri/dev/digitalocean/terraform-provider-digitalocean
make build
```

No need to run `terraform init` again—Terraform uses the binary from the dev override path.