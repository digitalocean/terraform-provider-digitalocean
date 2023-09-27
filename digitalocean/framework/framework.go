package framework

import (
	"context"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/util"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var (
	_ provider.Provider = &digitaloceanFrameworkProvider{}

	apiTokenEnvVars = []string{
		"DIGITALOCEAN_TOKEN",
		"DIGITALOCEAN_ACCESS_TOKEN",
	}
	apiURLEnvVars      = "DIGITALOCEAN_API_URL"
	defaultAPIEndpoint = "https://api.digitalocean.com"
)

type digitaloceanFrameworkProvider struct {
}

func New() provider.Provider {
	return &digitaloceanFrameworkProvider{}
}

func (p *digitaloceanFrameworkProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "digitalocean"
}

func (p *digitaloceanFrameworkProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"token": schema.StringAttribute{
				Optional:    true,
				Description: "The token key for API operations.",
			},
			"api_endpoint": schema.StringAttribute{
				Optional: true,
				// DefaultFunc doesn't exist in the framework
				// DefaultFunc: schema.EnvDefaultFunc("DIGITALOCEAN_API_URL", "https://api.digitalocean.com"),
				Description: "The URL to use for the DigitalOcean API.",
			},
			"spaces_endpoint": schema.StringAttribute{
				Optional: true,
				// DefaultFunc: schema.EnvDefaultFunc("SPACES_ENDPOINT_URL", "https://{{.Region}}.digitaloceanspaces.com"),
				Description: "The URL to use for the DigitalOcean Spaces API.",
			},
			"spaces_access_id": schema.StringAttribute{
				Optional: true,
				// DefaultFunc: schema.EnvDefaultFunc("SPACES_ACCESS_KEY_ID", nil),
				Description: "The access key ID for Spaces API operations.",
			},
			"spaces_secret_key": schema.StringAttribute{
				Optional: true,
				// DefaultFunc: schema.EnvDefaultFunc("SPACES_SECRET_ACCESS_KEY", nil),
				Description: "The secret access key for Spaces API operations.",
			},
			"requests_per_second": schema.NumberAttribute{
				Optional: true,
				// DefaultFunc: schema.EnvDefaultFunc("DIGITALOCEAN_REQUESTS_PER_SECOND", 0.0),
				Description: "The rate of requests per second to limit the HTTP client.",
			},
			"http_retry_max": schema.NumberAttribute{
				Optional: true,
				// DefaultFunc: schema.EnvDefaultFunc("DIGITALOCEAN_HTTP_RETRY_MAX", 4),
				Description: "The maximum number of retries on a failed API request.",
			},
			"http_retry_wait_min": schema.NumberAttribute{
				Optional: true,
				// DefaultFunc: schema.EnvDefaultFunc("DIGITALOCEAN_HTTP_RETRY_WAIT_MIN", 1.0),
				Description: "The minimum wait time (in seconds) between failed API requests.",
			},
			"http_retry_wait_max": schema.NumberAttribute{
				Optional: true,
				// DefaultFunc: schema.EnvDefaultFunc("DIGITALOCEAN_HTTP_RETRY_WAIT_MAX", 30.0),
				Description: "The maximum wait time (in seconds) between failed API requests.",
			},
		},
	}
}

func (p *digitaloceanFrameworkProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config config.Config
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags := p.ConfigureConfigDefaults(ctx, &config)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
	}

	// resp.Diagnostics.Append(p.ConfigureCallbackFunc(p, &req, &config)...)
	// if resp.Diagnostics.HasError() {
	// 	return
	// }

	client, err := config.Client("dev")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error configuring provider",
			err.Error(),
		)
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *digitaloceanFrameworkProvider) ConfigureConfigDefaults(ctx context.Context, config *config.Config) diag.Diagnostics {
	var diags diag.Diagnostics

	if config.Token == nil {
		token, err := util.GetMultiEnvVar(apiTokenEnvVars...)
		if err == nil {
			config.Token = godo.PtrTo(token)
		}
	}

	if config.APIEndpoint == nil {
		url, err := util.GetMultiEnvVar(apiURLEnvVars)
		if err == nil {
			config.APIEndpoint = godo.PtrTo(url)
		}

		if config.APIEndpoint == nil {
			config.APIEndpoint = godo.PtrTo(defaultAPIEndpoint)
		}
	}

	// TODO: Load all other defaults and/or env vars

	return diags
}

func (p *digitaloceanFrameworkProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
	// return []func() resource.Resource{
	// 	func() resource.Resource {
	// 		return resourceExample{}
	// 	},
	// }
}

func (p *digitaloceanFrameworkProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
	// return []func() datasource.DataSource{
	// 	func() datasource.DataSource {
	// 		return dataSourceExample{}
	// 	},
	// }
}
