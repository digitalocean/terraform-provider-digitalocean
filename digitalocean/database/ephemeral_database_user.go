package database

import (
	"context"
	"fmt"
	"net/http"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ ephemeral.EphemeralResourceWithConfigure = (*databaseUserEphemeralResource)(nil)

// NewDatabaseUserEphemeralResource returns an ephemeral resource that reads
// database user credentials without persisting them to Terraform state.
func NewDatabaseUserEphemeralResource() ephemeral.EphemeralResource {
	return &databaseUserEphemeralResource{}
}

type databaseUserEphemeralResource struct {
	client *godo.Client
}

type databaseUserEphemeralModel struct {
	ClusterID       types.String `tfsdk:"cluster_id"`
	Name            types.String `tfsdk:"name"`
	Role            types.String `tfsdk:"role"`
	Password        types.String `tfsdk:"password"`
	AccessCert      types.String `tfsdk:"access_cert"`
	AccessKey       types.String `tfsdk:"access_key"`
	MySQLAuthPlugin types.String `tfsdk:"mysql_auth_plugin"`
}

func (r *databaseUserEphemeralResource) Metadata(_ context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database_user"
}

func (r *databaseUserEphemeralResource) Schema(_ context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Reads a DigitalOcean database user's credentials without storing them in Terraform state. Use the digitalocean_database_user data source for non-sensitive metadata such as ACL settings.",
		Attributes: map[string]schema.Attribute{
			"cluster_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the database cluster.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the database user.",
			},
			"role": schema.StringAttribute{
				Computed:    true,
				Description: "The role of the database user, either primary or normal.",
			},
			"password": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "The password of the database user. This is not set for MongoDB users.",
			},
			"access_cert": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "Access certificate for TLS client authentication. Only set for Kafka users.",
			},
			"access_key": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "Access key for TLS client authentication. Only set for Kafka users.",
			},
			"mysql_auth_plugin": schema.StringAttribute{
				Computed:    true,
				Description: "The authentication method of the MySQL user.",
			},
		},
	}
}

func (r *databaseUserEphemeralResource) Configure(_ context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
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

func (r *databaseUserEphemeralResource) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Provider not configured", "The DigitalOcean provider client is nil.")
		return
	}

	var data databaseUserEphemeralModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterID := data.ClusterID.ValueString()
	name := data.Name.ValueString()

	user, apiResp, err := r.client.Databases.GetUser(ctx, clusterID, name)
	if err != nil {
		if apiResp != nil && apiResp.StatusCode == http.StatusNotFound {
			resp.Diagnostics.AddError("Database user not found", fmt.Sprintf("Database user %q was not found in cluster %s.", name, clusterID))
			return
		}
		resp.Diagnostics.AddError("Error retrieving database user", err.Error())
		return
	}
	if user == nil {
		resp.Diagnostics.AddError("Error retrieving database user", fmt.Sprintf("The API returned no database user for %q in cluster %s.", name, clusterID))
		return
	}

	populateDatabaseUserEphemeralModel(&data, user)

	if user.Password == "" && user.AccessCert == "" && user.AccessKey == "" {
		resp.Diagnostics.AddWarning(
			"No credentials returned",
			fmt.Sprintf("The API returned no password, access certificate, or access key for database user %q. MongoDB user passwords are only available when the user is created.", name),
		)
	}

	resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
}

func populateDatabaseUserEphemeralModel(data *databaseUserEphemeralModel, user *godo.DatabaseUser) {
	data.Role = types.StringValue(user.Role)
	data.Password = stringValueOrNull(user.Password)
	data.AccessCert = stringValueOrNull(user.AccessCert)
	data.AccessKey = stringValueOrNull(user.AccessKey)

	if user.MySQLSettings != nil && user.MySQLSettings.AuthPlugin != "" {
		data.MySQLAuthPlugin = types.StringValue(user.MySQLSettings.AuthPlugin)
	} else {
		data.MySQLAuthPlugin = types.StringNull()
	}
}

func stringValueOrNull(value string) types.String {
	if value == "" {
		return types.StringNull()
	}
	return types.StringValue(value)
}
