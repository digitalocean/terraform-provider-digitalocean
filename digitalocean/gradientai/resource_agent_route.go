package gradientai

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceDigitalOceanAgentRoute() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanAgentRouteCreate,
		ReadContext:   resourceDigitalOceanAgentRouteRead,
		UpdateContext: resourceDigitalOceanAgentRouteUpdate,
		DeleteContext: resourceDigitalOceanAgentRouteDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"parent_agent_uuid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The UUID of the parent agent.",
			},
			"child_agent_uuid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The UUID of the child agent.",
			},
			"if_case": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "if-case condition for the route.",
			},
			"route_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A name for the route.",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of the linkage",
			},
			"rollback": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceDigitalOceanAgentRouteCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	route := &godo.AgentRouteCreateRequest{
		ParentAgentUuid: d.Get("parent_agent_uuid").(string),
		ChildAgentUuid:  d.Get("child_agent_uuid").(string),
		IfCase:          d.Get("if_case").(string),
		RouteName:       d.Get("route_name").(string),
	}

	createdRoute, _, err := client.GradientAI.AddAgentRoute(
		ctx,
		d.Get("parent_agent_uuid").(string),
		d.Get("child_agent_uuid").(string),
		route,
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating agent route: %w", err))
	}

	d.SetId(createdRoute.UUID)

	return resourceDigitalOceanAgentRouteRead(ctx, d, meta)
}

func resourceDigitalOceanAgentRouteRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	routeUUID := d.Id()
	if routeUUID == "" {
		return diag.Errorf("route UUID is required")
	}

	route, _, err := client.GradientAI.GetAgent(ctx, routeUUID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error reading agent route: %w", err))
	}

	d.Set("parent_agent_uuid", route.ParentAgents)
	d.Set("child_agent_uuid", route.ChildAgents)
	d.Set("if_case", route.IfCase)
	d.Set("route_name", route.RouteName)
	d.Set("uuid", route.RouteUuid)

	return nil
}

func resourceDigitalOceanAgentRouteUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	parentAgentUUID := d.Get("parent_agent_uuid").(string)
	childAgentUUID := d.Get("child_agent_uuid").(string)

	updateRequest := &godo.AgentRouteUpdateRequest{}
	if d.HasChange("if_case") {
		updateRequest.IfCase = d.Get("if_case").(string)
	}
	if d.HasChange("route_name") {
		updateRequest.RouteName = d.Get("route_name").(string)
	}

	_, _, err := client.GradientAI.UpdateAgentRoute(ctx, parentAgentUUID, childAgentUUID, updateRequest)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating agent route: %w", err))
	}

	return resourceDigitalOceanAgentRouteRead(ctx, d, meta)
}

func resourceDigitalOceanAgentRouteDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	routeUUID := d.Id()
	if routeUUID == "" {
		return diag.Errorf("route UUID is required for deletion")
	}

	parentAgentUUID := d.Get("parent_agent_uuid").(string)
	childAgentUUID := d.Get("child_agent_uuid").(string)

	_, resp, err := client.GradientAI.DeleteAgentRoute(ctx, parentAgentUUID, childAgentUUID)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error deleting agent route: %w", err))
	}

	d.SetId("")
	return nil
}
