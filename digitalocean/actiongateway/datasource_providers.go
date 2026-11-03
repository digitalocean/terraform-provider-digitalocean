package actiongateway

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/digitalocean/terraform-provider-digitalocean/internal/datalist"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceDigitalOceanActionGatewayProviders returns a data source that
// lists the Action Gateway tool providers and the parameters a connection to
// each needs.
func DataSourceDigitalOceanActionGatewayProviders() *schema.Resource {
	return datalist.NewResource(&datalist.ResourceConfig{
		RecordSchema:        toolProviderSchema(),
		ResultAttributeName: "providers",
		GetRecords:          getActionGatewayProviders,
		FlattenRecord:       flattenActionGatewayProvider,
	})
}

func toolProviderSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "Provider slug, used as tool_provider on connections",
		},
		"display_name": {
			Type: schema.TypeString,
		},
		"description": {
			Type: schema.TypeString,
		},
		"auth_types": {
			Type:        schema.TypeList,
			Description: "Credential kinds the provider accepts",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"scopes": {
			Type:        schema.TypeList,
			Description: "OAuth scopes the provider offers",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"connection_parameters": {
			Type:        schema.TypeList,
			Description: "Parameters a connection to the provider accepts",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"key": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"label": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"description": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"input_kind": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"required": {
						Type:     schema.TypeBool,
						Computed: true,
					},
					"allowed_values": {
						Type:     schema.TypeList,
						Computed: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},
				},
			},
		},
	}
}

func getActionGatewayProviders(meta interface{}, _ map[string]interface{}) ([]interface{}, error) {
	client := meta.(*config.CombinedConfig).GodoClient()

	resp, _, err := client.ActionGateway.Tools.ListProviders(context.Background())
	if err != nil {
		return nil, fmt.Errorf("Error listing Action Gateway tool providers: %s", err)
	}

	records := make([]interface{}, 0, len(resp.Providers))
	for _, p := range resp.Providers {
		records = append(records, p)
	}
	return records, nil
}

func flattenActionGatewayProvider(raw, _ interface{}, _ map[string]interface{}) (map[string]interface{}, error) {
	p := raw.(godo.ActionGatewayProvider)

	authTypes := p.AuthTypes
	if len(authTypes) == 0 && p.AuthType != "" {
		authTypes = []string{p.AuthType}
	}

	params := make([]interface{}, 0, len(p.ConnectionParameters))
	for _, param := range p.ConnectionParameters {
		params = append(params, map[string]interface{}{
			"key":            param.Key,
			"label":          param.Label,
			"description":    param.Description,
			"input_kind":     param.InputKind,
			"required":       param.Required,
			"allowed_values": param.AllowedValues,
		})
	}

	return map[string]interface{}{
		"name":                  p.Name,
		"display_name":          p.DisplayName,
		"description":           p.Description,
		"auth_types":            authTypes,
		"scopes":                p.Scopes,
		"connection_parameters": params,
	}, nil
}
