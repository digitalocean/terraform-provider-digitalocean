package actiongateway

import (
	"context"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// DataSourceDigitalOceanActionGatewayConnection returns a data source that
// looks up an Action Gateway connection by `id`, or by `tool_provider` and
// `user_id`.
func DataSourceDigitalOceanActionGatewayConnection() *schema.Resource {
	recordSchema := computedOnly(connectionResourceSchema())

	recordSchema["id"] = &schema.Schema{
		Type:          schema.TypeString,
		Optional:      true,
		Computed:      true,
		Description:   "ID of the connection",
		ValidateFunc:  validation.NoZeroValues,
		ConflictsWith: []string{"tool_provider", "user_id"},
		AtLeastOneOf:  []string{"id", "tool_provider"},
	}
	recordSchema["tool_provider"] = &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Computed:     true,
		Description:  "Slug of the connection's tool provider",
		ValidateFunc: validation.NoZeroValues,
		RequiredWith: []string{"user_id"},
	}
	recordSchema["user_id"] = &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Computed:     true,
		Description:  "Actor ID the connection acts for",
		ValidateFunc: validateUserID,
		RequiredWith: []string{"tool_provider"},
	}

	return &schema.Resource{
		ReadContext: dataSourceDigitalOceanActionGatewayConnectionRead,
		Schema:      recordSchema,
	}
}

func dataSourceDigitalOceanActionGatewayConnectionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	id := d.Get("id").(string)
	if id == "" {
		provider := d.Get("tool_provider").(string)
		userID := d.Get("user_id").(string)
		matches, err := findConnections(ctx, client, provider, userID)
		if err != nil {
			return diag.Errorf("Error listing Action Gateway connections: %s", err)
		}
		switch len(matches) {
		case 0:
			return diag.Errorf("no Action Gateway connection found for tool provider %s and user ID %s", provider, userID)
		case 1:
			id = matches[0].ID
		default:
			return diag.Errorf("too many Action Gateway connections found for tool provider %s and user ID %s (found %d, expected 1)", provider, userID, len(matches))
		}
	}

	found, _, err := client.ActionGateway.Connections.Get(ctx, id)
	if err != nil {
		return diag.Errorf("Error retrieving Action Gateway connection: %s", err)
	}
	if found.Connection == nil {
		return diag.Errorf("Error retrieving Action Gateway connection (%s): response contained no connection", id)
	}

	d.SetId(found.Connection.ID)
	for k, v := range flattenConnection(found.Connection, found.Authorization) {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("error setting %s: %s", k, err)
		}
	}
	return nil
}

// findConnections lists the visible connections for a provider whose user ID
// is exactly userID. The API's user_id filter matches substrings.
func findConnections(ctx context.Context, client *godo.Client, provider, userID string) ([]godo.ActionGatewayConnection, error) {
	opt := &godo.ActionGatewayConnectionListOptions{
		Provider: provider,
		UserID:   userID,
		Page:     1,
		PerPage:  100,
	}

	var out []godo.ActionGatewayConnection
	for {
		page, _, err := client.ActionGateway.Connections.List(ctx, opt)
		if err != nil {
			return nil, err
		}
		for _, conn := range page.Connections {
			if conn.Provider == provider && conn.UserID == userID {
				out = append(out, conn)
			}
		}
		if page.Pagination == nil || len(page.Connections) == 0 || opt.Page*opt.PerPage >= page.Pagination.Total {
			break
		}
		opt.Page++
	}
	return out, nil
}
