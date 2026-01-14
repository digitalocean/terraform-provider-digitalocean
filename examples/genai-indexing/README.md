# Gradient AI Indexing APIs Example

This example demonstrates how to use the DigitalOcean Gradient AI Indexing APIs with Terraform.

## Features

This example shows how to:

- List indexing jobs for a knowledge base
- Get data sources for an indexing job  
- Retrieve detailed status of an indexing job
- Cancel a running indexing job

## Prerequisites

- DigitalOcean account with Gradient AI services enabled
- Knowledge base UUID
- Terraform >= 1.0

## Usage

1. Set your knowledge base UUID:

```bash
export TF_VAR_knowledge_base_uuid="your-knowledge-base-uuid"
```

2. Initialize and apply:

```bash
terraform init
terraform apply
```

3. To cancel a running job:

```bash
terraform apply -var="cancel_running_job=true"
```

## Configuration

The example uses these variables:

- `knowledge_base_uuid` - UUID of your knowledge base (required)
- `cancel_running_job` - Whether to cancel running jobs (default: false)

## Outputs

The example outputs:

- `indexing_jobs` - List of all indexing jobs
- `job_data_sources` - Data sources for the first job
- `job_status` - Detailed status of the first job
- `cancelled_job_status` - Status after cancellation (if performed)