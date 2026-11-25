package digitalocean

import (
	"context"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/database"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
	sdkschema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	_ provider.Provider                       = (*frameworkProvider)(nil)
	_ provider.ProviderWithEphemeralResources = (*frameworkProvider)(nil)
)

// frameworkProvider serves features that the SDKv2 provider cannot, such as
// ephemeral resources. It is muxed with the SDKv2 provider by ProviderServer.
type frameworkProvider struct {
	sdkProvider *sdkschema.Provider
}

// NewFrameworkProvider returns the terraform-plugin-framework provider.
func NewFrameworkProvider(sdkProvider *sdkschema.Provider) provider.Provider {
	return &frameworkProvider{sdkProvider: sdkProvider}
}

func (p *frameworkProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "digitalocean"
}

// Schema must match the SDKv2 provider schema exactly, as required by
// terraform-plugin-mux. Configuration is handled by the SDKv2 provider.
func (p *frameworkProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"token": schema.StringAttribute{
				Optional:    true,
				Description: "The token key for API operations.",
			},
			"api_endpoint": schema.StringAttribute{
				Optional:    true,
				Description: "The URL to use for the DigitalOcean API.",
			},
			"spaces_endpoint": schema.StringAttribute{
				Optional:    true,
				Description: "The URL to use for the DigitalOcean Spaces API.",
			},
			"spaces_access_id": schema.StringAttribute{
				Optional:    true,
				Description: "The access key ID for Spaces API operations.",
			},
			"spaces_secret_key": schema.StringAttribute{
				Optional:    true,
				Description: "The secret access key for Spaces API operations.",
			},
			"requests_per_second": schema.Float64Attribute{
				Optional:    true,
				Description: "The rate of requests per second to limit the HTTP client.",
			},
			"http_retry_max": schema.Int64Attribute{
				Optional:    true,
				Description: "The maximum number of retries on a failed API request.",
			},
			"http_retry_wait_min": schema.Float64Attribute{
				Optional:    true,
				Description: "The minimum wait time (in seconds) between failed API requests.",
			},
			"http_retry_wait_max": schema.Float64Attribute{
				Optional:    true,
				Description: "The maximum wait time (in seconds) between failed API requests.",
			},
		},
	}
}

func (p *frameworkProvider) Configure(_ context.Context, _ provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	if p.sdkProvider == nil {
		return
	}
	resp.EphemeralResourceData = p.sdkProvider.Meta()
}

func (p *frameworkProvider) EphemeralResources(_ context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{
		database.NewDatabaseUserEphemeralResource,
		database.NewDatabaseClusterEphemeralResource,
	}
}

func (p *frameworkProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}

func (p *frameworkProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

// ProviderServer returns a server that muxes the SDKv2 and framework providers.
func ProviderServer(ctx context.Context) (func() tfprotov5.ProviderServer, error) {
	sdkProvider := Provider()

	muxServer, err := tf5muxserver.NewMuxServer(ctx,
		sdkProvider.GRPCProvider,
		providerserver.NewProtocol5(NewFrameworkProvider(sdkProvider)),
	)
	if err != nil {
		return nil, err
	}

	return muxServer.ProviderServer, nil
}
