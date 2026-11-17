package functions

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceDigitalOceanFunctionsNamespace() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDigitalOceanFunctionsNamespaceRead,
		Schema: map[string]*schema.Schema{
			"namespace_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"namespace_id", "label"},
				ValidateFunc: validation.NoZeroValues,
				Description:  "The ID of the namespace, e.g. fn-b90faf52-2b42-49c2-9792-75edfbb6f397.",
			},
			"label": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"namespace_id", "label"},
				ValidateFunc: validation.NoZeroValues,
				Description:  "The namespace's label.",
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"api_host": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDigitalOceanFunctionsNamespaceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := functionsService(meta)

	var ns *godo.FunctionsNamespace
	if id, ok := d.GetOk("namespace_id"); ok {
		found, _, err := svc.GetNamespace(ctx, id.(string))
		if err != nil {
			return diag.Errorf("Error retrieving Functions namespace %s: %s", id, err)
		}
		ns = found
	} else {
		label := d.Get("label").(string)
		found, err := findNamespaceByLabel(ctx, svc, label)
		if err != nil {
			return diag.FromErr(err)
		}
		ns = found
	}

	d.SetId(ns.Namespace)
	if err := d.Set("namespace_id", ns.Namespace); err != nil {
		return diag.FromErr(err)
	}
	if err := setFunctionsNamespaceAttributes(d, ns); err != nil {
		return diag.Errorf("Error setting Functions namespace attributes: %s", err)
	}
	return nil
}

// findNamespaceByLabel returns the namespace with the given label. The list
// endpoint is not paginated and returns every namespace on the account.
func findNamespaceByLabel(ctx context.Context, svc godo.FunctionsService, label string) (*godo.FunctionsNamespace, error) {
	namespaces, _, err := svc.ListNamespaces(ctx)
	if err != nil {
		return nil, fmt.Errorf("Error listing Functions namespaces: %s", err)
	}

	var matches []godo.FunctionsNamespace
	for _, ns := range namespaces {
		if ns.Label == label {
			matches = append(matches, ns)
		}
	}

	switch len(matches) {
	case 0:
		return nil, fmt.Errorf("no Functions namespace found with label %q", label)
	case 1:
		// The list response does not include the access key; fetch the full record.
		ns, _, err := svc.GetNamespace(ctx, matches[0].Namespace)
		if err != nil {
			return nil, fmt.Errorf("Error retrieving Functions namespace %s: %s", matches[0].Namespace, err)
		}
		return ns, nil
	default:
		return nil, fmt.Errorf("found %d Functions namespaces with label %q; use namespace_id instead", len(matches), label)
	}
}
