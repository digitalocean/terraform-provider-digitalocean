package actiongateway

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/digitalocean/godo"
)

// godo has no provider-credentials service yet, so these requests are built
// here on godo's client, which supplies authentication, retries, and error
// decoding. The types follow the public API reference and can be replaced by
// godo's once it has them.

const providerCredentialsPath = "/v2/action-gateway/provider-credentials"

// providerCredential is a stored provider credential. The API returns
// metadata only: never an API key, a client ID, or a client secret.
type providerCredential struct {
	ID             string          `json:"id"`
	Provider       string          `json:"provider"`
	DisplayName    string          `json:"display_name"`
	Kind           string          `json:"kind"`
	BaseURL        string          `json:"base_url"`
	Status         string          `json:"status"`
	Scope          string          `json:"scope"`
	OwningUserID   string          `json:"owning_user_id"`
	ReusableByTeam bool            `json:"reusable_by_team"`
	RedirectURL    string          `json:"redirect_url"`
	CreatedAt      *godo.Timestamp `json:"created_at"`
	UpdatedAt      *godo.Timestamp `json:"updated_at"`
}

type providerCredentialCreateRequest struct {
	Provider             string                 `json:"provider"`
	Kind                 string                 `json:"kind"`
	DisplayName          string                 `json:"display_name"`
	APIKey               string                 `json:"api_key,omitempty"`
	ClientID             string                 `json:"client_id,omitempty"`
	ClientSecret         string                 `json:"client_secret,omitempty"`
	BaseURL              string                 `json:"base_url,omitempty"`
	CredentialParameters map[string]interface{} `json:"credential_parameters,omitempty"`
	ReusableByTeam       bool                   `json:"reusable_by_team,omitempty"`
	RedirectURL          string                 `json:"redirect_url,omitempty"`
}

// providerCredentialUpdateRequest changes only the fields it carries. The API
// rejects a request that carries none.
type providerCredentialUpdateRequest struct {
	DisplayName          string                 `json:"display_name,omitempty"`
	APIKey               string                 `json:"api_key,omitempty"`
	ClientID             string                 `json:"client_id,omitempty"`
	ClientSecret         string                 `json:"client_secret,omitempty"`
	BaseURL              string                 `json:"base_url,omitempty"`
	CredentialParameters map[string]interface{} `json:"credential_parameters,omitempty"`
	ReusableByTeam       bool                   `json:"reusable_by_team,omitempty"`
}

func (r *providerCredentialUpdateRequest) isEmpty() bool {
	return r.DisplayName == "" && r.APIKey == "" && r.ClientID == "" && r.ClientSecret == "" &&
		r.BaseURL == "" && len(r.CredentialParameters) == 0 && !r.ReusableByTeam
}

type providerCredentialResponse struct {
	Credential *providerCredential `json:"credential"`
	// RevokedConnections counts the connections deleted with the credential.
	// It is zero for every operation but a delete.
	RevokedConnections int `json:"revoked_connections"`
}

type providerCredentialsListResponse struct {
	Credentials []providerCredential `json:"credentials"`
}

func providerCredentialRequest[T any](ctx context.Context, client *godo.Client, method, path string, body interface{}) (*T, *godo.Response, error) {
	req, err := client.NewRequest(ctx, method, path, body)
	if err != nil {
		return nil, nil, err
	}
	out := new(T)
	resp, err := client.Do(ctx, req, out)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

func providerCredentialPath(id string) (string, error) {
	if strings.TrimSpace(id) == "" {
		return "", fmt.Errorf("provider credential ID cannot be empty")
	}
	return providerCredentialsPath + "/" + url.PathEscape(id), nil
}

func createProviderCredential(ctx context.Context, client *godo.Client, body *providerCredentialCreateRequest) (*providerCredentialResponse, *godo.Response, error) {
	return providerCredentialRequest[providerCredentialResponse](ctx, client, http.MethodPost, providerCredentialsPath, body)
}

func getProviderCredential(ctx context.Context, client *godo.Client, id string) (*providerCredentialResponse, *godo.Response, error) {
	path, err := providerCredentialPath(id)
	if err != nil {
		return nil, nil, err
	}
	return providerCredentialRequest[providerCredentialResponse](ctx, client, http.MethodGet, path, nil)
}

func updateProviderCredential(ctx context.Context, client *godo.Client, id string, body *providerCredentialUpdateRequest) (*providerCredentialResponse, *godo.Response, error) {
	path, err := providerCredentialPath(id)
	if err != nil {
		return nil, nil, err
	}
	return providerCredentialRequest[providerCredentialResponse](ctx, client, http.MethodPatch, path, body)
}

// deleteProviderCredential deletes the credential and every connection that
// uses it.
func deleteProviderCredential(ctx context.Context, client *godo.Client, id string) (*providerCredentialResponse, *godo.Response, error) {
	path, err := providerCredentialPath(id)
	if err != nil {
		return nil, nil, err
	}
	return providerCredentialRequest[providerCredentialResponse](ctx, client, http.MethodDelete, path, nil)
}

// listProviderCredentials lists the credentials visible to the caller. An
// empty provider lists every provider's. The list is not paginated.
func listProviderCredentials(ctx context.Context, client *godo.Client, provider string) ([]providerCredential, *godo.Response, error) {
	path := providerCredentialsPath
	if provider != "" {
		path += "?" + url.Values{"provider": []string{provider}}.Encode()
	}
	out, resp, err := providerCredentialRequest[providerCredentialsListResponse](ctx, client, http.MethodGet, path, nil)
	if err != nil {
		return nil, resp, err
	}
	return out.Credentials, resp, nil
}
