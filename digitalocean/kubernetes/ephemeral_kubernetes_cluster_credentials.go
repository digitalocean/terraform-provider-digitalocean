package kubernetes

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ ephemeral.EphemeralResourceWithConfigure = (*clusterCredentialsEphemeralResource)(nil)

// NewClusterCredentialsEphemeralResource returns an ephemeral resource that
// fetches Kubernetes cluster credentials without persisting them to state.
func NewClusterCredentialsEphemeralResource() ephemeral.EphemeralResource {
	return &clusterCredentialsEphemeralResource{}
}

type clusterCredentialsEphemeralResource struct {
	client *godo.Client
}

type clusterCredentialsModel struct {
	ClusterID            types.String `tfsdk:"cluster_id"`
	ExpirySeconds        types.Int64  `tfsdk:"expiry_seconds"`
	Host                 types.String `tfsdk:"host"`
	Token                types.String `tfsdk:"token"`
	ClusterCACertificate types.String `tfsdk:"cluster_ca_certificate"`
	ClientCertificate    types.String `tfsdk:"client_certificate"`
	ClientKey            types.String `tfsdk:"client_key"`
	RawConfig            types.String `tfsdk:"raw_config"`
	ExpiresAt            types.String `tfsdk:"expires_at"`
}

func (r *clusterCredentialsEphemeralResource) Metadata(_ context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kubernetes_cluster_credentials"
}

func (r *clusterCredentialsEphemeralResource) Schema(_ context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches credentials for a DigitalOcean Kubernetes cluster without storing them in Terraform state.",
		Attributes: map[string]schema.Attribute{
			"cluster_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the Kubernetes cluster.",
			},
			"expiry_seconds": schema.Int64Attribute{
				Optional:    true,
				Description: "The duration in seconds that the returned credentials will be valid. Defaults to the API's default expiry.",
			},
			"host": schema.StringAttribute{
				Computed:    true,
				Description: "The URL of the API server on the Kubernetes master node.",
			},
			"token": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "The DigitalOcean API access token used by clients to access the cluster.",
			},
			"cluster_ca_certificate": schema.StringAttribute{
				Computed:    true,
				Description: "The base64 encoded public certificate for the cluster's certificate authority.",
			},
			"client_certificate": schema.StringAttribute{
				Computed:    true,
				Description: "The PEM-encoded public certificate used by clients to access the cluster. Only available if token authentication is not supported on your cluster.",
			},
			"client_key": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "The PEM-encoded private key used by clients to access the cluster. Only available if token authentication is not supported on your cluster.",
			},
			"raw_config": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "The full contents of the Kubernetes cluster's kubeconfig file.",
			},
			"expires_at": schema.StringAttribute{
				Computed:    true,
				Description: "The date and time when the credentials will expire and need to be regenerated.",
			},
		},
	}
}

func (r *clusterCredentialsEphemeralResource) Configure(_ context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	meta, ok := req.ProviderData.(*config.CombinedConfig)
	if !ok {
		resp.Diagnostics.AddError("Unexpected provider data", fmt.Sprintf("Expected *config.CombinedConfig, got %T.", req.ProviderData))
		return
	}

	r.client = meta.GodoClient()
}

func (r *clusterCredentialsEphemeralResource) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var data clusterCredentialsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterID := data.ClusterID.ValueString()

	cluster, _, err := r.client.Kubernetes.Get(ctx, clusterID)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving Kubernetes cluster", err.Error())
		return
	}
	if cluster.SSO != nil && cluster.SSO.Required {
		resp.Diagnostics.AddError("Credentials unavailable", fmt.Sprintf("Kubernetes cluster %s requires SSO, so the API does not issue token credentials.", clusterID))
		return
	}

	credsReq := &godo.KubernetesClusterCredentialsGetRequest{}
	if !data.ExpirySeconds.IsNull() {
		expirySeconds := int(data.ExpirySeconds.ValueInt64())
		credsReq.ExpirySeconds = &expirySeconds
	}

	creds, _, err := r.client.Kubernetes.GetCredentials(ctx, clusterID, credsReq)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving Kubernetes cluster credentials", err.Error())
		return
	}

	rawConfig, err := RenderKubeconfig(cluster.Name, cluster.RegionSlug, creds)
	if err != nil {
		resp.Diagnostics.AddError("Error rendering kubeconfig", err.Error())
		return
	}

	data.Host = types.StringValue(creds.Server)
	data.Token = types.StringValue(creds.Token)
	data.ClusterCACertificate = types.StringValue(base64.StdEncoding.EncodeToString(creds.CertificateAuthorityData))
	data.ClientCertificate = types.StringValue(string(creds.ClientCertificateData))
	data.ClientKey = types.StringValue(string(creds.ClientKeyData))
	data.RawConfig = types.StringValue(string(rawConfig))
	data.ExpiresAt = types.StringValue(creds.ExpiresAt.Format(time.RFC3339))

	resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
}
