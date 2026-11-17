package actiongateway

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"slices"
	"strings"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	providerCredentialKindOAuthClient = "oauth_client"
	providerCredentialKindAPIKey      = "api_key"

	providerCredentialKindPrefix   = "PROVIDER_CREDENTIAL_KIND_"
	providerCredentialStatusPrefix = "PROVIDER_CREDENTIAL_STATUS_"

	providerCredentialStatusDeleted = "deleted"
)

// providerCredentialRegistrationArguments cannot change after registration, so
// Read keeps them as configured and only an import fills them from the API.
// Replacing a credential deletes every connection that uses it, and a value
// the API reported in another spelling must never plan that.
var providerCredentialRegistrationArguments = []string{"tool_provider", "kind", "redirect_url"}

// providerCredentialUpdateArguments are the arguments an update sends.
var providerCredentialUpdateArguments = []string{
	"display_name",
	"api_key",
	"client_id",
	"client_secret",
	"base_url",
	"credential_parameters",
	"reusable_by_team",
}

// ResourceDigitalOceanActionGatewayProviderCredential returns the
// digitalocean_action_gateway_provider_credential resource: an API key or an
// OAuth client your team registers for a tool provider. Secrets rotate in
// place, so a credential keeps its ID and its connections.
func ResourceDigitalOceanActionGatewayProviderCredential() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanActionGatewayProviderCredentialCreate,
		ReadContext:   resourceDigitalOceanActionGatewayProviderCredentialRead,
		UpdateContext: resourceDigitalOceanActionGatewayProviderCredentialUpdate,
		DeleteContext: resourceDigitalOceanActionGatewayProviderCredentialDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDigitalOceanActionGatewayProviderCredentialImport,
		},
		CustomizeDiff: customdiff.All(
			validateProviderCredentialDiff,
			// The API does not return the client ID, so an imported credential
			// has none in state. Setting it then is an update, which the API
			// accepts only for the registered ID. Changing a known ID means
			// another OAuth client, which needs a new credential.
			customdiff.ForceNewIfChange("client_id", func(_ context.Context, old, _, _ interface{}) bool {
				return old.(string) != ""
			}),
			// Sharing with the team cannot be undone in place.
			customdiff.ForceNewIfChange("reusable_by_team", func(_ context.Context, old, new, _ interface{}) bool {
				return old.(bool) && !new.(bool)
			}),
			customdiff.ComputedIf("base_url", derivesProviderCredentialBaseURL),
			customdiff.ComputedIf("scope", sharesProviderCredential),
			customdiff.ComputedIf("owning_user_id", sharesProviderCredential),
			customdiff.ComputedIf("updated_at", updatesProviderCredential),
		),
		Schema: providerCredentialResourceSchema(),
	}
}

func providerCredentialResourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"tool_provider": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			Description:  "Slug of the tool provider the credential is for, such as github or jira",
			ValidateFunc: validation.NoZeroValues,
		},
		"kind": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			Description:  "Kind of credential: oauth_client or api_key",
			ValidateFunc: validation.StringInSlice([]string{providerCredentialKindOAuthClient, providerCredentialKindAPIKey}, false),
		},
		"display_name": {
			Type:         schema.TypeString,
			Required:     true,
			Description:  "Label, unique among your team's credentials for the tool provider, ignoring case",
			ValidateFunc: validation.StringLenBetween(1, 255),
		},
		"api_key": {
			Type:         schema.TypeString,
			Optional:     true,
			Sensitive:    true,
			Description:  "For kind api_key: the key or token sent to the tool provider. It is write-only, and changing it rotates the key in place.",
			ValidateFunc: validation.StringLenBetween(1, 4096),
		},
		"client_id": {
			Type:         schema.TypeString,
			Optional:     true,
			Description:  "For kind oauth_client: client ID of your OAuth client at the tool provider. Changing it replaces the credential.",
			ValidateFunc: validation.StringLenBetween(1, 255),
		},
		"client_secret": {
			Type:         schema.TypeString,
			Optional:     true,
			Sensitive:    true,
			Description:  "For kind oauth_client: client secret of your OAuth client. It is write-only, and changing it rotates the secret in place.",
			ValidateFunc: validation.StringLenBetween(1, 1024),
		},
		"base_url": {
			Type:             schema.TypeString,
			Optional:         true,
			Computed:         true,
			Description:      "HTTPS origin of the tool provider's API, for a provider that accepts one. Derived from credential_parameters when those are set.",
			ConflictsWith:    []string{"credential_parameters"},
			ValidateFunc:     validation.All(validation.IsURLWithHTTPS, validation.StringLenBetween(1, 255)),
			DiffSuppressFunc: suppressSameOrigin,
		},
		"credential_parameters": {
			Type:          schema.TypeMap,
			Optional:      true,
			Description:   "For kind api_key: values for the credential parameters the tool provider declares, such as subdomain",
			ConflictsWith: []string{"base_url"},
			Elem:          &schema.Schema{Type: schema.TypeString},
		},
		"reusable_by_team": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Whether everyone on your team can see and use the credential instead of only its creator. Sharing cannot be undone, so turning it off replaces the credential.",
		},
		"redirect_url": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			Description:  "For kind oauth_client: HTTPS redirect URI registered with your OAuth client. Omit it to use DigitalOcean's callback.",
			ValidateFunc: validation.All(validation.IsURLWithHTTPS, validation.StringLenBetween(1, 1024)),
		},
		"status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Status of the credential: active or disabled",
		},
		"scope": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Who can use the credential: user for only its creator, or team",
		},
		"owning_user_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "DigitalOcean user ID of the creator. Empty for a credential shared with the team.",
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

// suppressSameOrigin hides a difference between two spellings of one HTTPS
// origin. The API stores the origin normalized.
func suppressSameOrigin(_, old, new string, _ *schema.ResourceData) bool {
	return old != "" && new != "" && normalizeOrigin(old) == normalizeOrigin(new)
}

func normalizeOrigin(raw string) string {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || u.Host == "" {
		return raw
	}
	if strings.Trim(u.Path, "/") != "" || u.RawQuery != "" || u.Fragment != "" || u.User != nil {
		return raw
	}
	scheme := strings.ToLower(u.Scheme)
	host := strings.ToLower(u.Hostname())
	if strings.Contains(host, ":") {
		host = "[" + host + "]"
	}
	if port := u.Port(); port != "" && !(scheme == "https" && port == "443") {
		host += ":" + port
	}
	return scheme + "://" + host
}

func validateProviderCredentialDiff(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	return validateProviderCredentialConfig(d.GetRawConfig())
}

// validateProviderCredentialConfig checks that the arguments a kind needs are
// configured and that those of the other kind are not, so a mismatch fails at
// plan time instead of as a 400 on apply. Values unknown until apply are not
// checked.
func validateProviderCredentialConfig(raw cty.Value) error {
	kindValue, ok := rawAttr(raw, "kind")
	if !ok || kindValue.IsNull() || !kindValue.IsKnown() || kindValue.Type() != cty.String {
		return nil
	}
	kind := kindValue.AsString()

	var required, rejected []string
	switch kind {
	case providerCredentialKindAPIKey:
		required = []string{"api_key"}
		rejected = []string{"client_id", "client_secret", "redirect_url"}
	case providerCredentialKindOAuthClient:
		required = []string{"client_id", "client_secret"}
		rejected = []string{"api_key", "credential_parameters"}
	default:
		return nil
	}

	for _, key := range required {
		if v, ok := rawAttr(raw, key); ok && v.IsWhollyKnown() && !isConfigured(v) {
			return fmt.Errorf("kind %q requires %s", kind, key)
		}
	}
	for _, key := range rejected {
		if v, ok := rawAttr(raw, key); ok && v.IsWhollyKnown() && isRawAttrSet(v) {
			return fmt.Errorf("%s cannot be set for kind %q", key, kind)
		}
	}
	return nil
}

// isRawAttrSet reports whether a configuration value holds anything. Unlike
// isConfigured, it counts the entries of a map.
func isRawAttrSet(v cty.Value) bool {
	if v.IsNull() {
		return false
	}
	if v.Type().IsMapType() || v.Type().IsObjectType() {
		return v.LengthInt() > 0
	}
	return isConfigured(v)
}

// plansChange reports whether the plan changes any of keys. Unlike HasChange,
// it leaves out a difference a DiffSuppressFunc hid, such as another spelling
// of the stored base URL.
func plansChange(d *schema.ResourceDiff, keys ...string) bool {
	for _, key := range keys {
		if len(d.GetChangedKeysPrefix(key)) > 0 {
			return true
		}
	}
	return false
}

// derivesProviderCredentialBaseURL reports whether an update changes the
// parameters the API derives the base URL from.
func derivesProviderCredentialBaseURL(_ context.Context, d *schema.ResourceDiff, _ interface{}) bool {
	if d.Id() == "" || !plansChange(d, "credential_parameters") {
		return false
	}
	// A configured base URL conflicts with credential_parameters, so only an
	// omitted one can be derived.
	baseURL, ok := rawAttr(d.GetRawConfig(), "base_url")
	return ok && baseURL.IsNull()
}

func sharesProviderCredential(_ context.Context, d *schema.ResourceDiff, _ interface{}) bool {
	return d.Id() != "" && plansChange(d, "reusable_by_team")
}

func updatesProviderCredential(_ context.Context, d *schema.ResourceDiff, _ interface{}) bool {
	return d.Id() != "" && plansChange(d, providerCredentialUpdateArguments...)
}

func resourceDigitalOceanActionGatewayProviderCredentialCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	req := &godo.ActionGatewayProviderCredentialCreateRequest{
		Provider:       d.Get("tool_provider").(string),
		Kind:           expandProviderCredentialKind(d.Get("kind").(string)),
		DisplayName:    d.Get("display_name").(string),
		APIKey:         d.Get("api_key").(string),
		ClientID:       d.Get("client_id").(string),
		ClientSecret:   d.Get("client_secret").(string),
		ReusableByTeam: d.Get("reusable_by_team").(bool),
		RedirectURL:    d.Get("redirect_url").(string),
	}
	// The API rejects a base URL sent alongside the parameters it is derived
	// from.
	if params := expandCredentialParameters(d.Get("credential_parameters")); len(params) > 0 {
		req.CredentialParameters = params
	} else {
		req.BaseURL = d.Get("base_url").(string)
	}

	created, _, err := client.ActionGateway.ProviderCredentials.Create(ctx, req)
	if err != nil {
		return diag.Errorf("Error creating Action Gateway provider credential: %s", err)
	}
	if created.Credential == nil || created.Credential.ID == "" {
		return diag.Errorf("Error creating Action Gateway provider credential: response contained no credential")
	}

	d.SetId(created.Credential.ID)
	log.Printf("[INFO] Action Gateway provider credential created, ID: %s", d.Id())

	return readProviderCredential(ctx, d, meta)
}

func resourceDigitalOceanActionGatewayProviderCredentialImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if diags := readProviderCredential(ctx, d, meta, providerCredentialRegistrationArguments...); diags.HasError() {
		return nil, fmt.Errorf("%s", diags[0].Summary)
	}
	if d.Id() == "" {
		return nil, fmt.Errorf("Action Gateway provider credential not found")
	}
	return []*schema.ResourceData{d}, nil
}

func resourceDigitalOceanActionGatewayProviderCredentialRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return readProviderCredential(ctx, d, meta)
}

// readProviderCredential refreshes the credential's state. Of the
// providerCredentialRegistrationArguments, it sets only those named in fill.
func readProviderCredential(ctx context.Context, d *schema.ResourceData, meta interface{}, fill ...string) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	found, resp, err := client.ActionGateway.ProviderCredentials.Get(ctx, d.Id())
	if err != nil {
		if isNotFound(resp) {
			log.Printf("[WARN] Action Gateway provider credential (%s) not found - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error retrieving Action Gateway provider credential: %s", err)
	}
	// Dropping the credential from state plans its creation again, so an
	// answer without one is an error rather than a missing credential.
	if found.Credential == nil {
		return diag.Errorf("Error retrieving Action Gateway provider credential (%s): response contained no credential", d.Id())
	}
	if normalizeProviderCredentialStatus(found.Credential.Status) == providerCredentialStatusDeleted {
		log.Printf("[WARN] Action Gateway provider credential (%s) is deleted - removing from state", d.Id())
		d.SetId("")
		return nil
	}

	values := flattenProviderCredential(found.Credential)
	for _, key := range providerCredentialRegistrationArguments {
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

func resourceDigitalOceanActionGatewayProviderCredentialUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	req := expandProviderCredentialUpdate(d)
	if !isEmptyProviderCredentialUpdate(req) {
		if _, _, err := client.ActionGateway.ProviderCredentials.Update(ctx, d.Id(), req); err != nil {
			// Without this the SDK records the configuration as applied, and
			// a secret that was never rotated would stop showing as a change.
			d.Partial(true)
			return diag.Errorf("Error updating Action Gateway provider credential (%s): %s", d.Id(), err)
		}
		log.Printf("[INFO] Action Gateway provider credential updated, ID: %s", d.Id())
	}

	return readProviderCredential(ctx, d, meta)
}

// expandProviderCredentialUpdate builds a request carrying only what changed.
func expandProviderCredentialUpdate(d *schema.ResourceData) *godo.ActionGatewayProviderCredentialUpdateRequest {
	req := &godo.ActionGatewayProviderCredentialUpdateRequest{}
	oauthClient := d.Get("kind").(string) == providerCredentialKindOAuthClient

	if d.HasChange("display_name") {
		req.DisplayName = d.Get("display_name").(string)
	}
	if d.HasChange("api_key") {
		req.APIKey = d.Get("api_key").(string)
	}
	// The API rotates an OAuth client's secret only alongside its registered
	// client ID, and a client ID on its own changes nothing. Both are sent
	// whenever either differs, which is the case after an import.
	if oauthClient && d.HasChanges("client_id", "client_secret") {
		req.ClientID = d.Get("client_id").(string)
		req.ClientSecret = d.Get("client_secret").(string)
	}
	if params := expandCredentialParameters(d.Get("credential_parameters")); d.HasChange("credential_parameters") && len(params) > 0 {
		req.CredentialParameters = params
	} else if baseURL := d.Get("base_url").(string); d.HasChange("base_url") && baseURL != "" {
		req.BaseURL = godo.PtrTo(baseURL)
		if oauthClient {
			req.ClientID = d.Get("client_id").(string)
		}
	}
	// Only sharing is sent: the API rejects false, and turning sharing off
	// replaces the credential instead.
	if d.HasChange("reusable_by_team") && d.Get("reusable_by_team").(bool) {
		req.ReusableByTeam = true
	}
	return req
}

// isEmptyProviderCredentialUpdate reports whether the request changes
// nothing. The API rejects such a request.
func isEmptyProviderCredentialUpdate(req *godo.ActionGatewayProviderCredentialUpdateRequest) bool {
	return req.DisplayName == "" && req.APIKey == "" && req.ClientID == "" && req.ClientSecret == "" &&
		req.BaseURL == nil && len(req.CredentialParameters) == 0 && !req.ReusableByTeam
}

func resourceDigitalOceanActionGatewayProviderCredentialDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	deleted, resp, err := client.ActionGateway.ProviderCredentials.Delete(ctx, d.Id())
	if err != nil {
		if isNotFound(resp) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting Action Gateway provider credential: %s", err)
	}

	id := d.Id()
	log.Printf("[INFO] Action Gateway provider credential deleted, ID: %s", id)
	d.SetId("")

	if deleted.RevokedConnections > 0 {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "Action Gateway connections were deleted with the provider credential",
			Detail:   fmt.Sprintf("Deleting provider credential %s also deleted the %d connection(s) that used it.", id, deleted.RevokedConnections),
		}}
	}
	return nil
}

func flattenProviderCredential(cred *godo.ActionGatewayProviderCredential) map[string]interface{} {
	return map[string]interface{}{
		"tool_provider":    cred.Provider,
		"kind":             normalizeProviderCredentialKind(cred.Kind),
		"display_name":     cred.DisplayName,
		"base_url":         cred.BaseURL,
		"reusable_by_team": cred.ReusableByTeam,
		"redirect_url":     cred.RedirectURL,
		"status":           normalizeProviderCredentialStatus(cred.Status),
		"scope":            cred.Scope,
		"owning_user_id":   cred.OwningUserID,
		"created_at":       formatTimestamp(cred.CreatedAt),
		"updated_at":       formatTimestamp(cred.UpdatedAt),
	}
}

// expandProviderCredentialKind returns the API's spelling of a kind, such as
// PROVIDER_CREDENTIAL_KIND_API_KEY for api_key.
func expandProviderCredentialKind(kind string) string {
	return providerCredentialKindPrefix + strings.ToUpper(kind)
}

func normalizeProviderCredentialKind(kind string) string {
	return strings.ToLower(strings.TrimPrefix(kind, providerCredentialKindPrefix))
}

func normalizeProviderCredentialStatus(status string) string {
	return strings.ToLower(strings.TrimPrefix(status, providerCredentialStatusPrefix))
}

func expandCredentialParameters(v interface{}) map[string]interface{} {
	raw, _ := v.(map[string]interface{})
	if len(raw) == 0 {
		return nil
	}
	out := make(map[string]interface{}, len(raw))
	for k, item := range raw {
		out[k] = item
	}
	return out
}
