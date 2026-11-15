package actiongateway

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// providerCredentialWriteOnlyArguments are the resource arguments the API
// never returns, so the data source cannot publish them.
var providerCredentialWriteOnlyArguments = []string{"api_key", "client_id", "client_secret", "credential_parameters"}

// DataSourceDigitalOceanActionGatewayProviderCredential returns a data source
// that looks up an Action Gateway provider credential by `id`, or by
// `tool_provider` and `display_name`.
func DataSourceDigitalOceanActionGatewayProviderCredential() *schema.Resource {
	recordSchema := computedOnly(providerCredentialResourceSchema())
	for _, key := range providerCredentialWriteOnlyArguments {
		delete(recordSchema, key)
	}

	recordSchema["id"] = &schema.Schema{
		Type:          schema.TypeString,
		Optional:      true,
		Computed:      true,
		Description:   "ID of the provider credential",
		ValidateFunc:  validation.NoZeroValues,
		ConflictsWith: []string{"tool_provider", "display_name"},
		AtLeastOneOf:  []string{"id", "tool_provider"},
	}
	recordSchema["tool_provider"] = &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Computed:     true,
		Description:  "Slug of the credential's tool provider",
		ValidateFunc: validation.NoZeroValues,
		RequiredWith: []string{"display_name"},
	}
	recordSchema["display_name"] = &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Computed:     true,
		Description:  "Label of the credential, matched ignoring case",
		ValidateFunc: validation.NoZeroValues,
		RequiredWith: []string{"tool_provider"},
	}

	return &schema.Resource{
		ReadContext: dataSourceDigitalOceanActionGatewayProviderCredentialRead,
		Schema:      recordSchema,
	}
}

func dataSourceDigitalOceanActionGatewayProviderCredentialRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	var cred *providerCredential
	if id := d.Get("id").(string); id != "" {
		found, _, err := getProviderCredential(ctx, client, id)
		if err != nil {
			return diag.Errorf("Error retrieving Action Gateway provider credential (%s): %s", id, err)
		}
		if found.Credential == nil {
			return diag.Errorf("Error retrieving Action Gateway provider credential (%s): response contained no credential", id)
		}
		cred = found.Credential
	} else {
		var err error
		cred, err = findProviderCredential(ctx, client, d.Get("tool_provider").(string), d.Get("display_name").(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(cred.ID)
	for k, v := range flattenProviderCredential(cred) {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("error setting %s: %s", k, err)
		}
	}
	return nil
}

// findProviderCredential resolves a tool provider and a display name to
// exactly one credential visible to the caller. Display names are unique for a
// provider ignoring case, which is how they are matched here.
func findProviderCredential(ctx context.Context, client *godo.Client, provider, displayName string) (*providerCredential, error) {
	creds, _, err := listProviderCredentials(ctx, client, provider)
	if err != nil {
		return nil, fmt.Errorf("Error listing Action Gateway provider credentials: %s", err)
	}

	var matches []providerCredential
	for _, cred := range creds {
		if strings.EqualFold(cred.Provider, provider) && strings.EqualFold(cred.DisplayName, displayName) {
			matches = append(matches, cred)
		}
	}

	switch len(matches) {
	case 0:
		return nil, fmt.Errorf("no Action Gateway provider credential named %q found for tool provider %s", displayName, provider)
	case 1:
		return &matches[0], nil
	default:
		ids := make([]string, 0, len(matches))
		for _, cred := range matches {
			ids = append(ids, cred.ID)
		}
		sort.Strings(ids)
		return nil, fmt.Errorf("too many Action Gateway provider credentials named %q found for tool provider %s (%d, expected 1): %s; set id to choose one",
			displayName, provider, len(matches), strings.Join(ids, ", "))
	}
}
