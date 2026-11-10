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

// DataSourceDigitalOceanActionGatewayEndpoint returns a data source that looks
// up an Action Gateway endpoint by `id`, or by `name` with an optional
// `actor_id` filter, and publishes the endpoint with its head version.
func DataSourceDigitalOceanActionGatewayEndpoint() *schema.Resource {
	recordSchema := computedOnly(endpointResourceSchema())

	recordSchema["id"] = &schema.Schema{
		Type:          schema.TypeString,
		Optional:      true,
		Computed:      true,
		Description:   "ID of the endpoint",
		ValidateFunc:  validation.NoZeroValues,
		ConflictsWith: []string{"name"},
		AtLeastOneOf:  []string{"id", "name"},
	}
	recordSchema["name"] = &schema.Schema{
		Type:          schema.TypeString,
		Optional:      true,
		Computed:      true,
		Description:   "Name of the endpoint",
		ValidateFunc:  validation.NoZeroValues,
		ConflictsWith: []string{"id"},
		AtLeastOneOf:  []string{"id", "name"},
	}
	recordSchema["actor_id"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		Description: "Actor ID to filter the name lookup by. Endpoints with a different actor are ignored.",
	}

	return &schema.Resource{
		ReadContext: dataSourceDigitalOceanActionGatewayEndpointRead,
		Schema:      recordSchema,
	}
}

func dataSourceDigitalOceanActionGatewayEndpointRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	id := d.Get("id").(string)
	if id == "" {
		var err error
		id, err = findEndpointID(ctx, client, d.Get("name").(string), d.Get("actor_id").(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	found, _, err := client.ActionGateway.Endpoints.Get(ctx, id)
	if err != nil {
		return diag.Errorf("Error retrieving Action Gateway endpoint (%s): %s", id, err)
	}
	if found.Endpoint == nil {
		return diag.Errorf("Error retrieving Action Gateway endpoint (%s): response contained no endpoint", id)
	}

	values := flattenEndpoint(found.Endpoint)
	values["id"] = found.Endpoint.ID
	values["actor_id"] = found.Endpoint.ActorID
	for k, v := range values {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("error setting %s: %s", k, err)
		}
	}
	d.SetId(found.Endpoint.ID)
	return nil
}

// findEndpointID resolves a name, and an optional actor filter, to exactly one
// endpoint. Several matches are an error naming the candidates, because an
// automatically created default endpoint is exempt from the team's name
// uniqueness and can share a name with a managed one.
func findEndpointID(ctx context.Context, client *godo.Client, name, actorID string) (string, error) {
	opt := &godo.ActionGatewayEndpointListOptions{Name: name, PageSize: 100}
	if actorID != "" {
		opt.ActorID = actorID
	}

	var matches []godo.ActionGatewayEndpoint
	for {
		page, _, err := client.ActionGateway.Endpoints.List(ctx, opt)
		if err != nil {
			return "", fmt.Errorf("Error listing Action Gateway endpoints: %s", err)
		}
		for _, endpoint := range page.Endpoints {
			if endpoint.Name != name {
				continue
			}
			if actorID != "" && endpoint.ActorID != actorID {
				continue
			}
			matches = append(matches, endpoint)
		}
		if page.NextPageToken == "" || len(page.Endpoints) == 0 {
			break
		}
		opt.PageToken = page.NextPageToken
	}

	switch len(matches) {
	case 0:
		if actorID != "" {
			return "", fmt.Errorf("no Action Gateway endpoint named %s found for actor ID %s", name, actorID)
		}
		return "", fmt.Errorf("no Action Gateway endpoint named %s found", name)
	case 1:
		return matches[0].ID, nil
	default:
		ids := make([]string, 0, len(matches))
		for _, endpoint := range matches {
			ids = append(ids, endpoint.ID)
		}
		sort.Strings(ids)
		return "", fmt.Errorf("too many Action Gateway endpoints named %s found (%d, expected 1): %s; set id to choose one",
			name, len(matches), strings.Join(ids, ", "))
	}
}
