package actiongateway

import (
	"context"
	"fmt"
	"log"
	"slices"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const connectionStatusRevoked = "revoked"

// connectionRequestArguments are the arguments Read keeps as configured: they
// cannot change after creation, and the API may widen scopes or normalize
// parameter values. They are filled from the API when the connection is
// imported, and scopes also when the configuration omits them.
var connectionRequestArguments = []string{"scopes", "connection_parameters"}

// ResourceDigitalOceanActionGatewayConnection returns the
// digitalocean_action_gateway_connection resource. The API has no update
// operation, so every argument forces a new connection.
func ResourceDigitalOceanActionGatewayConnection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanActionGatewayConnectionCreate,
		ReadContext:   resourceDigitalOceanActionGatewayConnectionRead,
		DeleteContext: resourceDigitalOceanActionGatewayConnectionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDigitalOceanActionGatewayConnectionImport,
		},
		Schema: connectionResourceSchema(),
	}
}

func connectionResourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"tool_provider": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			Description:  "Slug of the tool provider to connect, such as github or jira",
			ValidateFunc: validation.NoZeroValues,
		},
		"user_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			Description:  "Actor ID the connection acts for. Sessions whose actor ID matches use this connection.",
			ValidateFunc: validateUserID,
		},
		"scopes": {
			Type:        schema.TypeSet,
			Optional:    true,
			Computed:    true,
			ForceNew:    true,
			Description: "OAuth scopes to request. Defaults to every scope the provider offers.",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"connection_parameters": {
			Type:        schema.TypeMap,
			Optional:    true,
			ForceNew:    true,
			Description: "Values for the provider's connection parameters, such as site_url",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"credential_id": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			Description:  "ID of a team provider credential to connect through. When omitted, DigitalOcean's shared OAuth application is used.",
			ValidateFunc: validation.NoZeroValues,
		},
		"vpc": {
			Type:        schema.TypeList,
			Optional:    true,
			ForceNew:    true,
			MaxItems:    1,
			Description: "Private network the connection's calls are routed through",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"vpc_uuid": {
						Type:         schema.TypeString,
						Required:     true,
						ForceNew:     true,
						Description:  "UUID of a VPC owned by your team",
						ValidateFunc: validation.IsUUID,
					},
					"region": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Region of the VPC",
					},
					"destination": {
						Type:        schema.TypeList,
						Required:    true,
						ForceNew:    true,
						MinItems:    1,
						Description: "Private hosts the connection's calls may reach through the VPC",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"host": {
									Type:         schema.TypeString,
									Required:     true,
									ForceNew:     true,
									Description:  "Host name or IP address that resolves only to private addresses",
									ValidateFunc: validation.NoZeroValues,
								},
								"port": {
									Type:         schema.TypeInt,
									Required:     true,
									ForceNew:     true,
									Description:  "Destination port",
									ValidateFunc: validation.IsPortNumber,
								},
								"allowed_ip_cidrs": {
									Type:        schema.TypeList,
									Optional:    true,
									ForceNew:    true,
									Description: "CIDR ranges that every address the host resolves to must fall within",
									Elem: &schema.Schema{
										Type:         schema.TypeString,
										ValidateFunc: validation.IsCIDR,
									},
								},
							},
						},
					},
				},
			},
		},
		"tool_provider_display_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Human-readable name of the tool provider",
		},
		"status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Status of the connection: pending, active, or expired",
		},
		"credential_kind": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Kind of credential the connection uses: digitalocean_oauth, private_oauth, or team_api_key",
		},
		"owning_user_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "DigitalOcean user ID of the user who created the connection",
		},
		"authorization": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Authorization the account holder still has to complete. Set only while the connection is pending.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"status": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"connect_url": {
						Type:      schema.TypeString,
						Computed:  true,
						Sensitive: true,
					},
					"verification_code": {
						Type:      schema.TypeString,
						Computed:  true,
						Sensitive: true,
					},
					"expires_at": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"granted_at": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "When the account holder completed authorization",
		},
		"created_at": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"updated_at": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func resourceDigitalOceanActionGatewayConnectionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	req := &godo.ActionGatewayConnectionCreateRequest{
		Provider: d.Get("tool_provider").(string),
		UserID:   d.Get("user_id").(string),
		Scopes:   expandStringSet(d.Get("scopes")),
	}
	if v, ok := d.GetOk("connection_parameters"); ok {
		req.ConnectionParameters = v.(map[string]interface{})
	}
	if v, ok := d.GetOk("credential_id"); ok {
		req.Credential = &godo.ActionGatewayConnectionCredential{
			TeamCredential: &godo.ActionGatewayTeamCredential{CredentialID: v.(string)},
		}
	}
	if v, ok := d.GetOk("vpc"); ok {
		req.Network = expandConnectionNetwork(v.([]interface{}))
	}

	// Instead of creating a connection, the API returns or restarts an
	// existing one for the same tool provider and user ID. Terraform must not
	// take ownership of such a connection: destroying this resource would
	// revoke it while another configuration, or the resource that
	// create_before_destroy is replacing, still uses it.
	existing, err := findConnections(ctx, client, req.Provider, req.UserID)
	if err != nil {
		return diag.Errorf("Error listing Action Gateway connections: %s", err)
	}

	created, _, err := client.ActionGateway.Connections.Create(ctx, req)
	if err != nil {
		return diag.Errorf("Error creating Action Gateway connection: %s", err)
	}
	if created.Connection == nil || created.Connection.ID == "" {
		return diag.Errorf("Error creating Action Gateway connection: response contained no connection")
	}
	for _, conn := range existing {
		if conn.ID == created.Connection.ID {
			return diag.Errorf("Error creating Action Gateway connection: connection %s for tool provider %q and user ID %q already exists (status %s). "+
				"Import it with `terraform import digitalocean_action_gateway_connection.<name> %s`, or revoke it before creating this resource.",
				conn.ID, req.Provider, req.UserID, conn.Status, conn.ID)
		}
	}

	d.SetId(created.Connection.ID)
	log.Printf("[INFO] Action Gateway connection created, ID: %s", d.Id())

	// Omitted scopes default to every scope the provider offers, so record
	// the scopes the API chose.
	var fill []string
	if !isSetInConfig(d, "scopes") {
		fill = append(fill, "scopes")
	}
	return readConnection(ctx, d, meta, fill...)
}

// isSetInConfig reports whether key is present in the configuration. Unlike
// GetOk, it tells an explicitly empty value such as `scopes = []` apart from
// an omitted argument.
func isSetInConfig(d *schema.ResourceData, key string) bool {
	raw := d.GetRawConfig()
	if raw.IsNull() || !raw.IsKnown() {
		_, ok := d.GetOk(key)
		return ok
	}
	return !raw.GetAttr(key).IsNull()
}

func resourceDigitalOceanActionGatewayConnectionImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if diags := readConnection(ctx, d, meta, connectionRequestArguments...); diags.HasError() {
		return nil, fmt.Errorf("%s", diags[0].Summary)
	}
	if d.Id() == "" {
		return nil, fmt.Errorf("Action Gateway connection not found or revoked")
	}
	return []*schema.ResourceData{d}, nil
}

func resourceDigitalOceanActionGatewayConnectionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return readConnection(ctx, d, meta)
}

// readConnection refreshes the connection's state. Of the
// connectionRequestArguments, it sets only those named in fill.
func readConnection(ctx context.Context, d *schema.ResourceData, meta interface{}, fill ...string) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	found, resp, err := client.ActionGateway.Connections.Get(ctx, d.Id())
	if err != nil {
		if isNotFound(resp) {
			log.Printf("[WARN] Action Gateway connection (%s) not found - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error retrieving Action Gateway connection: %s", err)
	}
	if found.Connection == nil || found.Connection.Status == connectionStatusRevoked {
		log.Printf("[WARN] Action Gateway connection (%s) is revoked - removing from state", d.Id())
		d.SetId("")
		return nil
	}

	values := flattenConnection(found.Connection, found.Authorization)
	for _, key := range connectionRequestArguments {
		if !slices.Contains(fill, key) {
			delete(values, key)
		}
	}

	for k, v := range values {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("error setting %s: %s", k, err)
		}
	}
	return nil
}

func resourceDigitalOceanActionGatewayConnectionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	_, resp, err := client.ActionGateway.Connections.Delete(ctx, d.Id())
	if err != nil {
		if isNotFound(resp) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error revoking Action Gateway connection: %s", err)
	}

	log.Printf("[INFO] Action Gateway connection revoked, ID: %s", d.Id())
	d.SetId("")
	return nil
}

func flattenConnection(conn *godo.ActionGatewayConnection, authz *godo.ActionGatewayConnectionAuthorization) map[string]interface{} {
	scopes := conn.Scopes
	grantedAt := conn.GrantedAt
	if conn.OAuth != nil {
		if len(conn.OAuth.Scopes) > 0 {
			scopes = conn.OAuth.Scopes
		}
		if conn.OAuth.GrantedAt != nil {
			grantedAt = conn.OAuth.GrantedAt
		}
	}

	params := make(map[string]interface{}, len(conn.ConnectionParameters))
	for k, v := range conn.ConnectionParameters {
		if s, ok := v.(string); ok {
			params[k] = s
		} else {
			params[k] = fmt.Sprint(v)
		}
	}

	return map[string]interface{}{
		"tool_provider":              conn.Provider,
		"tool_provider_display_name": conn.ProviderDisplayName,
		"user_id":                    conn.UserID,
		"status":                     conn.Status,
		"scopes":                     scopes,
		"connection_parameters":      params,
		"credential_kind":            conn.CredentialKind,
		"credential_id":              conn.CredentialID,
		"owning_user_id":             conn.OwningUserID,
		"vpc":                        flattenConnectionNetwork(conn.Network),
		"authorization":              flattenConnectionAuthorization(authz),
		"granted_at":                 formatTimestamp(grantedAt),
		"created_at":                 formatTimestamp(conn.CreatedAt),
		"updated_at":                 formatTimestamp(conn.UpdatedAt),
	}
}

func expandConnectionNetwork(raw []interface{}) *godo.ActionGatewayConnectionCreateNetwork {
	if len(raw) == 0 || raw[0] == nil {
		return nil
	}
	vpc := raw[0].(map[string]interface{})

	var destinations []godo.ActionGatewayConnectionDestination
	for _, item := range vpc["destination"].([]interface{}) {
		dest := item.(map[string]interface{})
		destinations = append(destinations, godo.ActionGatewayConnectionDestination{
			Host:           dest["host"].(string),
			Port:           uint32(dest["port"].(int)),
			AllowedIPCIDRs: expandStringList(dest["allowed_ip_cidrs"]),
		})
	}

	return &godo.ActionGatewayConnectionCreateNetwork{
		VPC: &godo.ActionGatewayConnectionCreateVPC{
			VPCUUID:      vpc["vpc_uuid"].(string),
			Destinations: destinations,
		},
	}
}

func flattenConnectionNetwork(network *godo.ActionGatewayConnectionNetwork) []interface{} {
	if network == nil || network.VPC == nil {
		return nil
	}

	destinations := make([]interface{}, 0, len(network.VPC.Destinations))
	for _, dest := range network.VPC.Destinations {
		destinations = append(destinations, map[string]interface{}{
			"host":             dest.Host,
			"port":             int(dest.Port),
			"allowed_ip_cidrs": dest.AllowedIPCIDRs,
		})
	}

	return []interface{}{map[string]interface{}{
		"vpc_uuid":    network.VPC.VPCUUID,
		"region":      network.VPC.Region,
		"destination": destinations,
	}}
}

func flattenConnectionAuthorization(authz *godo.ActionGatewayConnectionAuthorization) []interface{} {
	if authz == nil {
		return nil
	}
	return []interface{}{map[string]interface{}{
		"status":            authz.Status,
		"connect_url":       authz.ConnectURL,
		"verification_code": authz.VerificationCode,
		"expires_at":        formatTimestamp(authz.ExpiresAt),
	}}
}
