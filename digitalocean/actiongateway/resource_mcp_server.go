package actiongateway

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	mcpCredentialSourceNone       = "none"
	mcpCredentialSourceSecret     = "secret"
	mcpCredentialSourceConnection = "connection"
)

var mcpOAuthArguments = []string{
	"oauth_client_id",
	"oauth_client_secret",
	"oauth_authorize_url",
	"oauth_token_url",
	"oauth_scopes",
}

// ResourceDigitalOceanActionGatewayMCPServer returns the
// digitalocean_action_gateway_mcp_server resource. Only the description and
// the enabled tools can change in place; the endpoint and credentials cannot
// be changed without registering the server again.
func ResourceDigitalOceanActionGatewayMCPServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanActionGatewayMCPServerCreate,
		ReadContext:   resourceDigitalOceanActionGatewayMCPServerRead,
		UpdateContext: resourceDigitalOceanActionGatewayMCPServerUpdate,
		DeleteContext: resourceDigitalOceanActionGatewayMCPServerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
		},
		CustomizeDiff: customdiff.All(
			validateMCPServerCredentials,
			customdiff.ComputedIf("discovered_tools", enabledToolsChanged),
			customdiff.ComputedIf("tool_count", enabledToolsChanged),
			customdiff.ComputedIf("sync_status", enabledToolsChanged),
			customdiff.ComputedIf("sync_error", enabledToolsChanged),
			customdiff.ComputedIf("last_synced_at", enabledToolsChanged),
			customdiff.ComputedIf("protocol_version", enabledToolsChanged),
			customdiff.ComputedIf("updated_at", enabledToolsChanged),
		),
		Schema: map[string]*schema.Schema{
			"server_ref": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Identifier for the server, unique within your team. Prefixes the slugs of the server's tools.",
				ValidateFunc: validation.StringMatch(serverRefRegexp, "must start with a lowercase letter and contain at most 64 characters from [a-z0-9-]"),
			},
			"endpoint": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "HTTPS URL of the server's MCP endpoint",
				ValidateFunc: validation.IsURLWithHTTPS,
			},
			"transport": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  "MCP transport. Only streamable_http is supported.",
				ValidateFunc: validation.StringInSlice([]string{"streamable_http"}, false),
			},
			"credential_ref_source": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  "How requests to the server authenticate: none (the default), secret, or connection",
				ValidateFunc: validation.StringInSlice([]string{mcpCredentialSourceNone, mcpCredentialSourceSecret, mcpCredentialSourceConnection}, false),
			},
			"credential_ref": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				Description:   "For credential_ref_source secret: reference to a secret your team stores with DigitalOcean",
				ConflictsWith: []string{"api_key"},
				ValidateFunc:  validation.NoZeroValues,
			},
			"api_key": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				ForceNew:      true,
				Description:   "For credential_ref_source secret: the key or token sent to the server. It is write-only.",
				ConflictsWith: []string{"credential_ref"},
				ValidateFunc:  validation.NoZeroValues,
			},
			"oauth_client_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Description:  "For credential_ref_source connection: client ID of your team's OAuth client for the server",
				ValidateFunc: validation.NoZeroValues,
			},
			"oauth_client_secret": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				ForceNew:     true,
				Description:  "For credential_ref_source connection: client secret of your team's OAuth client. It is write-only.",
				ValidateFunc: validation.NoZeroValues,
			},
			"oauth_authorize_url": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Description:  "For credential_ref_source connection: the OAuth authorization endpoint",
				ValidateFunc: validation.IsURLWithHTTPS,
			},
			"oauth_token_url": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Description:  "For credential_ref_source connection: the OAuth token endpoint",
				ValidateFunc: validation.IsURLWithHTTPS,
			},
			"oauth_scopes": {
				Type:        schema.TypeSet,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    64,
				Description: "For credential_ref_source connection: the scopes to request from each user",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"oauth_authorization_ttl_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  "How long, in seconds, a user's authorization may be reused before they must consent again. 0 means 30 days.",
				ValidateFunc: validation.IntAtLeast(0),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Description shown on the server's catalog card",
				ValidateFunc: validation.StringLenBetween(0, 1024),
			},
			"enabled_tools": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Description: "Slugs of the discovered tools to enable, such as <server_ref>_<tool>. Every other tool of the server is disabled.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"sync_user_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "For credential_ref_source connection: the user whose authorization is used to discover and verify the server's tools",
				ValidateFunc: validateUserID(),
			},
			"sync_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the most recent tool discovery",
			},
			"sync_error": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Error from the most recent tool discovery, if any",
			},
			"last_synced_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "When the server's tools were last discovered",
			},
			"protocol_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "MCP protocol version the server reported",
			},
			"tool_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of tools discovered on the server",
			},
			"discovered_tools": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Tools discovered on the server",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tool_slug": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"quarantined": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"quarantine_reason": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
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

func enabledToolsChanged(_ context.Context, d *schema.ResourceDiff, _ interface{}) bool {
	return d.Id() != "" && d.HasChange("enabled_tools")
}

func validateMCPServerCredentials(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	return validateMCPServerCredentialConfig(d.GetRawConfig())
}

// validateMCPServerCredentialConfig checks that the arguments a credential
// source needs are configured. Values unknown until apply are not checked.
func validateMCPServerCredentialConfig(raw cty.Value) error {
	if raw.IsNull() || !raw.IsKnown() {
		return nil
	}

	source := raw.GetAttr("credential_ref_source")
	if !source.IsWhollyKnown() {
		return nil
	}
	sourceValue := mcpCredentialSourceNone
	if !source.IsNull() {
		sourceValue = source.AsString()
	}

	switch sourceValue {
	case mcpCredentialSourceSecret:
		ref, apiKey := raw.GetAttr("credential_ref"), raw.GetAttr("api_key")
		if ref.IsWhollyKnown() && apiKey.IsWhollyKnown() && !isConfigured(ref) && !isConfigured(apiKey) {
			return fmt.Errorf("credential_ref_source %q requires one of credential_ref or api_key", sourceValue)
		}
	case mcpCredentialSourceConnection:
		for _, key := range mcpOAuthArguments {
			if v := raw.GetAttr(key); v.IsWhollyKnown() && !isConfigured(v) {
				return fmt.Errorf("credential_ref_source %q requires %s", sourceValue, key)
			}
		}
	}
	return nil
}

func isConfigured(v cty.Value) bool {
	if v.IsNull() {
		return false
	}
	if v.Type() == cty.String {
		return v.AsString() != ""
	}
	if v.Type().IsSetType() || v.Type().IsListType() {
		return v.LengthInt() > 0
	}
	return true
}

func resourceDigitalOceanActionGatewayMCPServerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	req := &godo.ActionGatewayMCPServerCreateRequest{
		ServerRef:           d.Get("server_ref").(string),
		Endpoint:            d.Get("endpoint").(string),
		Transport:           d.Get("transport").(string),
		CredentialRefSource: d.Get("credential_ref_source").(string),
		CredentialRef:       d.Get("credential_ref").(string),
		APIKey:              d.Get("api_key").(string),
		Description:         d.Get("description").(string),
		OAuthClientID:       d.Get("oauth_client_id").(string),
		OAuthClientSecret:   d.Get("oauth_client_secret").(string),
		OAuthAuthorizeURL:   d.Get("oauth_authorize_url").(string),
		OAuthTokenURL:       d.Get("oauth_token_url").(string),
		OAuthScopes:         expandStringSet(d.Get("oauth_scopes")),
	}
	if v, ok := d.GetOk("oauth_authorization_ttl_seconds"); ok {
		req.OAuthAuthorizationTTLSeconds = strconv.Itoa(v.(int))
	}

	created, _, err := client.ActionGateway.MCPServers.Create(ctx, req)
	if err != nil {
		return diag.Errorf("Error registering Action Gateway MCP server: %s", err)
	}
	if created.MCPServer == nil || created.MCPServer.ServerRef == "" {
		return diag.Errorf("Error registering Action Gateway MCP server: response contained no server")
	}

	d.SetId(created.MCPServer.ServerRef)
	log.Printf("[INFO] Action Gateway MCP server registered: %s", d.Id())

	// A discovery that cannot complete yet, such as one waiting on a user's
	// authorization, is reported as a warning rather than an error so the
	// server is not tainted. Read leaves enabled_tools empty, and the next
	// apply discovers and enables the tools again.
	var diags diag.Diagnostics
	source := d.Get("credential_ref_source").(string)
	syncUserID := d.Get("sync_user_id").(string)
	if source == mcpCredentialSourceConnection && syncUserID == "" {
		if _, ok := d.GetOk("enabled_tools"); ok {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Action Gateway MCP server tools were not enabled",
				Detail:   "Discovering the tools of a server registered with credential_ref_source \"connection\" requires sync_user_id.",
			})
		}
		return append(diags, resourceDigitalOceanActionGatewayMCPServerRead(ctx, d, meta)...)
	}

	if err := discoverMCPServerTools(ctx, client, d.Id(), syncUserID, d.Timeout(schema.TimeoutCreate)); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Action Gateway MCP server tools were not discovered",
			Detail:   err.Error(),
		})
	} else if v, ok := d.GetOk("enabled_tools"); ok {
		if err := setMCPServerEnabledTools(ctx, client, d.Id(), expandStringSet(v), syncUserID); err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	}

	return append(diags, resourceDigitalOceanActionGatewayMCPServerRead(ctx, d, meta)...)
}

func resourceDigitalOceanActionGatewayMCPServerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	found, resp, err := client.ActionGateway.MCPServers.Get(ctx, d.Id())
	if err != nil {
		if isNotFound(resp) {
			log.Printf("[WARN] Action Gateway MCP server (%s) not found - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error retrieving Action Gateway MCP server: %s", err)
	}
	if found.MCPServer == nil {
		d.SetId("")
		return nil
	}

	tools, _, err := client.ActionGateway.MCPServers.ListTools(ctx, d.Id())
	if err != nil {
		return diag.Errorf("Error listing tools of Action Gateway MCP server (%s): %s", d.Id(), err)
	}

	values, err := flattenMCPServer(found.MCPServer, tools.Tools)
	if err != nil {
		return diag.FromErr(err)
	}
	for k, v := range values {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("error setting %s: %s", k, err)
		}
	}
	return nil
}

func resourceDigitalOceanActionGatewayMCPServerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	if d.HasChange("description") {
		req := &godo.ActionGatewayMCPServerUpdateRequest{Description: d.Get("description").(string)}
		if _, _, err := client.ActionGateway.MCPServers.Update(ctx, d.Id(), req); err != nil {
			return diag.Errorf("Error updating Action Gateway MCP server (%s): %s", d.Id(), err)
		}
	}

	if d.HasChange("enabled_tools") {
		syncUserID := d.Get("sync_user_id").(string)
		if err := discoverMCPServerTools(ctx, client, d.Id(), syncUserID, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return diag.FromErr(err)
		}
		if err := setMCPServerEnabledTools(ctx, client, d.Id(), expandStringSet(d.Get("enabled_tools")), syncUserID); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceDigitalOceanActionGatewayMCPServerRead(ctx, d, meta)
}

func resourceDigitalOceanActionGatewayMCPServerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	_, resp, err := client.ActionGateway.MCPServers.Delete(ctx, d.Id())
	if err != nil {
		if isNotFound(resp) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting Action Gateway MCP server: %s", err)
	}

	log.Printf("[INFO] Action Gateway MCP server deleted: %s", d.Id())
	d.SetId("")
	return nil
}

// discoverMCPServerTools refreshes the server's discovered tools. Discovery
// that outlasts the request continues in the background; the API documents
// polling the server until last_synced_at or sync_error changes.
func discoverMCPServerTools(ctx context.Context, client *godo.Client, ref, userID string, timeout time.Duration) error {
	before, _, err := client.ActionGateway.MCPServers.Get(ctx, ref)
	if err != nil {
		return fmt.Errorf("error retrieving Action Gateway MCP server (%s): %s", ref, err)
	}
	if before.MCPServer == nil {
		return fmt.Errorf("error retrieving Action Gateway MCP server (%s): response contained no server", ref)
	}
	previous := *before.MCPServer

	result, _, err := client.ActionGateway.MCPServers.Resync(ctx, ref, &godo.ActionGatewayMCPServerResyncRequest{UserID: userID})
	if err != nil {
		return fmt.Errorf("error discovering tools of Action Gateway MCP server (%s): %s", ref, err)
	}
	if result.Authorization != nil {
		return fmt.Errorf("user %q has not authorized Action Gateway MCP server %q yet; complete that user's authorization and apply again", userID, ref)
	}

	synced := result.MCPServer
	if result.Pending || synced == nil {
		stateConf := &retry.StateChangeConf{
			Pending: []string{"discovering"},
			Target:  []string{"done"},
			Refresh: func() (interface{}, string, error) {
				current, _, err := client.ActionGateway.MCPServers.Get(ctx, ref)
				if err != nil {
					return nil, "", err
				}
				s := current.MCPServer
				if s == nil || (s.LastSyncedAt == previous.LastSyncedAt && s.SyncError == previous.SyncError) {
					return current, "discovering", nil
				}
				return s, "done", nil
			},
			Timeout:    timeout,
			Delay:      2 * time.Second,
			MinTimeout: 2 * time.Second,
		}
		raw, err := stateConf.WaitForStateContext(ctx)
		if err != nil {
			return fmt.Errorf("error waiting for tool discovery of Action Gateway MCP server (%s): %s", ref, err)
		}
		synced = raw.(*godo.ActionGatewayMCPServer)
	}

	if synced.SyncError != "" {
		return fmt.Errorf("tool discovery of Action Gateway MCP server (%s) failed: %s", ref, synced.SyncError)
	}
	return nil
}

func setMCPServerEnabledTools(ctx context.Context, client *godo.Client, ref string, slugs []string, userID string) error {
	req := &godo.ActionGatewayMCPServerToolsUpdateRequest{EnabledToolSlugs: slugs, UserID: userID}
	if req.EnabledToolSlugs == nil {
		req.EnabledToolSlugs = []string{}
	}
	if _, _, err := client.ActionGateway.MCPServers.UpdateTools(ctx, ref, req); err != nil {
		return fmt.Errorf("error enabling tools of Action Gateway MCP server (%s): %s", ref, err)
	}
	return nil
}

func flattenMCPServer(s *godo.ActionGatewayMCPServer, tools []godo.ActionGatewayMCPServerTool) (map[string]interface{}, error) {
	ttl := 0
	if s.OAuthAuthorizationTTLSeconds != "" {
		parsed, err := strconv.Atoi(s.OAuthAuthorizationTTLSeconds)
		if err != nil {
			return nil, fmt.Errorf("error parsing oauth_authorization_ttl_seconds %q: %s", s.OAuthAuthorizationTTLSeconds, err)
		}
		ttl = parsed
	}

	enabled := make([]string, 0, len(tools))
	discovered := make([]interface{}, 0, len(tools))
	for _, tool := range tools {
		if tool.Enabled {
			enabled = append(enabled, tool.ToolSlug)
		}
		discovered = append(discovered, map[string]interface{}{
			"name":              tool.Name,
			"tool_slug":         tool.ToolSlug,
			"description":       tool.Description,
			"enabled":           tool.Enabled,
			"quarantined":       tool.Quarantined,
			"quarantine_reason": tool.QuarantineReason,
		})
	}

	return map[string]interface{}{
		"server_ref":                      s.ServerRef,
		"endpoint":                        s.Endpoint,
		"transport":                       s.Transport,
		"credential_ref_source":           s.CredentialRefSource,
		"credential_ref":                  s.CredentialRef,
		"oauth_authorize_url":             s.OAuthAuthorizeURL,
		"oauth_scopes":                    s.OAuthScopes,
		"oauth_authorization_ttl_seconds": ttl,
		"description":                     s.Description,
		"enabled_tools":                   enabled,
		"sync_status":                     s.SyncStatus,
		"sync_error":                      s.SyncError,
		"last_synced_at":                  s.LastSyncedAt,
		"protocol_version":                s.ProtocolVersion,
		"tool_count":                      s.ToolCount,
		"discovered_tools":                discovered,
		"created_at":                      s.CreatedAt,
		"updated_at":                      s.UpdatedAt,
	}, nil
}
