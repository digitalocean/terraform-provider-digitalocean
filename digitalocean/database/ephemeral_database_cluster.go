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

var _ ephemeral.EphemeralResourceWithConfigure = (*databaseClusterEphemeralResource)(nil)

// NewDatabaseClusterEphemeralResource returns an ephemeral resource that reads
// database cluster connection credentials without persisting them to state.
func NewDatabaseClusterEphemeralResource() ephemeral.EphemeralResource {
	return &databaseClusterEphemeralResource{}
}

type databaseClusterEphemeralResource struct {
	client *godo.Client
}

type databaseClusterEphemeralModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Engine      types.String `tfsdk:"engine"`
	Version     types.String `tfsdk:"version"`
	Region      types.String `tfsdk:"region"`
	Host        types.String `tfsdk:"host"`
	PrivateHost types.String `tfsdk:"private_host"`
	UIHost      types.String `tfsdk:"ui_host"`
	Port        types.Int64  `tfsdk:"port"`
	UIPort      types.Int64  `tfsdk:"ui_port"`
	Database    types.String `tfsdk:"database"`
	UIDatabase  types.String `tfsdk:"ui_database"`
	User        types.String `tfsdk:"user"`
	UIUser      types.String `tfsdk:"ui_user"`
	Password    types.String `tfsdk:"password"`
	UIPassword  types.String `tfsdk:"ui_password"`
	URI         types.String `tfsdk:"uri"`
	UIURI       types.String `tfsdk:"ui_uri"`
	PrivateURI  types.String `tfsdk:"private_uri"`
}

func (r *databaseClusterEphemeralResource) Metadata(_ context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database_cluster"
}

func (r *databaseClusterEphemeralResource) Schema(_ context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Reads a DigitalOcean database cluster's connection credentials without storing them in Terraform state. Use the digitalocean_database_cluster data source for non-sensitive metadata.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the database cluster.",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "The name of the database cluster.",
			},
			"engine": schema.StringAttribute{
				Computed:    true,
				Description: "The database engine slug.",
			},
			"version": schema.StringAttribute{
				Computed:    true,
				Description: "The database engine version.",
			},
			"region": schema.StringAttribute{
				Computed:    true,
				Description: "The region slug where the database cluster is running.",
			},
			"host": schema.StringAttribute{
				Computed:    true,
				Description: "The public hostname for the database cluster.",
			},
			"private_host": schema.StringAttribute{
				Computed:    true,
				Description: "The private hostname for the database cluster.",
			},
			"ui_host": schema.StringAttribute{
				Computed:    true,
				Description: "The hostname for the database cluster's web UI.",
			},
			"port": schema.Int64Attribute{
				Computed:    true,
				Description: "The public port for the database cluster.",
			},
			"ui_port": schema.Int64Attribute{
				Computed:    true,
				Description: "The port for the database cluster's web UI.",
			},
			"database": schema.StringAttribute{
				Computed:    true,
				Description: "The name of the default database.",
			},
			"ui_database": schema.StringAttribute{
				Computed:    true,
				Description: "The name of the database used by the web UI.",
			},
			"user": schema.StringAttribute{
				Computed:    true,
				Description: "The default database user.",
			},
			"ui_user": schema.StringAttribute{
				Computed:    true,
				Description: "The database user used by the web UI.",
			},
			"password": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "The password for the default database user.",
			},
			"ui_password": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "The password for the web UI database user.",
			},
			"uri": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "The public connection URI for the database cluster.",
			},
			"ui_uri": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "The connection URI for the database cluster's web UI.",
			},
			"private_uri": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "The private connection URI for the database cluster.",
			},
		},
	}
}

func (r *databaseClusterEphemeralResource) Configure(_ context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
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

func (r *databaseClusterEphemeralResource) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Provider not configured", "The DigitalOcean provider client is nil.")
		return
	}

	var data databaseClusterEphemeralModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterID := data.ID.ValueString()

	database, apiResp, err := r.client.Databases.Get(ctx, clusterID)
	if err != nil {
		if apiResp != nil && apiResp.StatusCode == http.StatusNotFound {
			resp.Diagnostics.AddError("Database cluster not found", fmt.Sprintf("Database cluster %s was not found.", clusterID))
			return
		}
		resp.Diagnostics.AddError("Error retrieving database cluster", err.Error())
		return
	}
	if database == nil {
		resp.Diagnostics.AddError("Error retrieving database cluster", fmt.Sprintf("The API returned no database cluster for %s.", clusterID))
		return
	}

	if err := populateDatabaseClusterEphemeralModel(&data, database); err != nil {
		resp.Diagnostics.AddError("Error building database cluster credentials", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
}

func populateDatabaseClusterEphemeralModel(data *databaseClusterEphemeralModel, database *godo.Database) error {
	data.Name = stringValueOrNull(database.Name)
	data.Engine = stringValueOrNull(database.EngineSlug)
	data.Version = stringValueOrNull(database.VersionSlug)
	data.Region = stringValueOrNull(database.RegionSlug)

	password := ""
	if database.Connection != nil {
		password = database.Connection.Password
		data.Host = stringValueOrNull(database.Connection.Host)
		data.Port = types.Int64Value(int64(database.Connection.Port))
		data.Database = stringValueOrNull(database.Connection.Database)
		data.User = stringValueOrNull(database.Connection.User)
		data.Password = stringValueOrNull(database.Connection.Password)

		uri, err := buildDBConnectionURIWithPassword(database.Connection, database.EngineSlug, database.Connection.Password)
		if err != nil {
			return fmt.Errorf("building public connection URI: %w", err)
		}
		data.URI = stringValueOrNull(uri)
	} else {
		data.Host = types.StringNull()
		data.Port = types.Int64Null()
		data.Database = types.StringNull()
		data.User = types.StringNull()
		data.Password = types.StringNull()
		data.URI = types.StringNull()
	}

	if database.PrivateConnection != nil {
		data.PrivateHost = stringValueOrNull(database.PrivateConnection.Host)

		privateURI, err := buildDBConnectionURIWithPassword(database.PrivateConnection, database.EngineSlug, password)
		if err != nil {
			return fmt.Errorf("building private connection URI: %w", err)
		}
		data.PrivateURI = stringValueOrNull(privateURI)
	} else {
		data.PrivateHost = types.StringNull()
		data.PrivateURI = types.StringNull()
	}

	if database.UIConnection != nil {
		data.UIHost = stringValueOrNull(database.UIConnection.Host)
		data.UIPort = types.Int64Value(int64(database.UIConnection.Port))
		data.UIDatabase = stringValueOrNull(database.UIConnection.Database)
		data.UIUser = stringValueOrNull(database.UIConnection.User)
		data.UIPassword = stringValueOrNull(database.UIConnection.Password)
		data.UIURI = stringValueOrNull(database.UIConnection.URI)
	} else {
		data.UIHost = types.StringNull()
		data.UIPort = types.Int64Null()
		data.UIDatabase = types.StringNull()
		data.UIUser = types.StringNull()
		data.UIPassword = types.StringNull()
		data.UIURI = types.StringNull()
	}

	return nil
}
