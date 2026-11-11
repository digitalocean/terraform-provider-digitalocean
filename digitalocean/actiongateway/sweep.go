package actiongateway

import (
	"context"
	"log"
	"strings"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("digitalocean_action_gateway_connection", &resource.Sweeper{
		Name: "digitalocean_action_gateway_connection",
		F:    sweepActionGatewayConnections,
	})

	resource.AddTestSweepers("digitalocean_action_gateway_toolbelt", &resource.Sweeper{
		Name: "digitalocean_action_gateway_toolbelt",
		F:    sweepActionGatewayToolbelts,
	})

	resource.AddTestSweepers("digitalocean_action_gateway_mcp_server", &resource.Sweeper{
		Name: "digitalocean_action_gateway_mcp_server",
		F:    sweepActionGatewayMCPServers,
	})

	resource.AddTestSweepers("digitalocean_action_gateway_endpoint", &resource.Sweeper{
		Name: "digitalocean_action_gateway_endpoint",
		F:    sweepActionGatewayEndpoints,
	})
}

func sweepActionGatewayConnections(region string) error {
	meta, err := sweep.SharedConfigForRegion(region)
	if err != nil {
		return err
	}
	client := meta.(*config.CombinedConfig).GodoClient()

	opt := &godo.ActionGatewayConnectionListOptions{UserID: sweep.TestNamePrefix, Page: 1, PerPage: 100}
	var stale []string
	for {
		page, _, err := client.ActionGateway.Connections.List(context.Background(), opt)
		if err != nil {
			return err
		}
		for _, conn := range page.Connections {
			if strings.HasPrefix(conn.UserID, sweep.TestNamePrefix) && conn.Status != connectionStatusRevoked {
				stale = append(stale, conn.ID)
			}
		}
		if page.Pagination == nil || len(page.Connections) == 0 || opt.Page*opt.PerPage >= page.Pagination.Total {
			break
		}
		opt.Page++
	}

	for _, id := range stale {
		log.Printf("[DEBUG] Revoking Action Gateway connection %s", id)
		if _, _, err := client.ActionGateway.Connections.Delete(context.Background(), id); err != nil {
			return err
		}
	}
	return nil
}

func sweepActionGatewayToolbelts(region string) error {
	meta, err := sweep.SharedConfigForRegion(region)
	if err != nil {
		return err
	}
	client := meta.(*config.CombinedConfig).GodoClient()

	opt := &godo.ActionGatewayToolbeltsListOptions{Page: 1, PerPage: 100}
	var stale []string
	for {
		page, _, err := client.ActionGateway.Toolbelts.List(context.Background(), opt)
		if err != nil {
			return err
		}
		for _, tb := range page.Toolbelts {
			if strings.HasPrefix(tb.Name, sweep.TestNamePrefix) {
				stale = append(stale, tb.Name)
			}
		}
		if page.Pagination == nil || len(page.Toolbelts) == 0 || opt.Page*opt.PerPage >= page.Pagination.Total {
			break
		}
		opt.Page++
	}

	for _, name := range stale {
		log.Printf("[DEBUG] Deleting Action Gateway toolbelt %s", name)
		if _, _, err := client.ActionGateway.Toolbelts.Delete(context.Background(), name); err != nil {
			return err
		}
	}
	return nil
}

// sweepActionGatewayEndpoints deletes the acceptance suite's endpoints. A
// version has no independent lifetime, so deleting an endpoint deletes every
// live version under it.
func sweepActionGatewayEndpoints(region string) error {
	meta, err := sweep.SharedConfigForRegion(region)
	if err != nil {
		return err
	}
	client := meta.(*config.CombinedConfig).GodoClient()

	opt := &godo.ActionGatewayEndpointListOptions{PageSize: 100}
	var stale []string
	for {
		page, _, err := client.ActionGateway.Endpoints.List(context.Background(), opt)
		if err != nil {
			return err
		}
		for _, endpoint := range page.Endpoints {
			if strings.HasPrefix(endpoint.Name, sweep.TestNamePrefix) {
				stale = append(stale, endpoint.ID)
			}
		}
		if page.NextPageToken == "" || len(page.Endpoints) == 0 {
			break
		}
		opt.PageToken = page.NextPageToken
	}

	for _, id := range stale {
		log.Printf("[DEBUG] Deleting Action Gateway endpoint %s", id)
		if _, err := client.ActionGateway.Endpoints.Delete(context.Background(), id); err != nil {
			return err
		}
	}
	return nil
}

func sweepActionGatewayMCPServers(region string) error {
	meta, err := sweep.SharedConfigForRegion(region)
	if err != nil {
		return err
	}
	client := meta.(*config.CombinedConfig).GodoClient()

	servers, _, err := client.ActionGateway.MCPServers.List(context.Background())
	if err != nil {
		return err
	}
	for _, s := range servers.MCPServers {
		if !strings.HasPrefix(s.ServerRef, sweep.TestNamePrefix) {
			continue
		}
		log.Printf("[DEBUG] Deleting Action Gateway MCP server %s", s.ServerRef)
		if _, _, err := client.ActionGateway.MCPServers.Delete(context.Background(), s.ServerRef); err != nil {
			return err
		}
	}
	return nil
}
