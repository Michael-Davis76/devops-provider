package provider

import (
    "context"
    "fmt"
    "os"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/path"
    "github.com/hashicorp/terraform-plugin-framework/provider"
    "github.com/hashicorp/terraform-plugin-framework/provider/schema"
    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
    _ provider.Provider = &devopsProvider{}
)

// New is a helper function hashicupsProvidermplify provider server and testing implementation.
func New(version string) func() provider.Provider {
    return func() provider.Provider {
        return &devopsProvider{
            version: version,
        }
    }
}

// hashicupsProvider is the provider implementation.
type devopsProvider struct {
    // version is set to the provider version on release, "dev" when the
    // provider is built and ran locally, and "test" when running acceptance
    // testing.
    version string
}

// devopsProviderModel is the provider data model.
type devopsProviderModel struct {
    Host types.String `tfsdk:"host"`
}

// Metadata returns the provider type name.
func (p *devopsProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
    resp.TypeName = "devops"
    resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *devopsProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
    resp.Schema = schema.Schema{
        Attributes: map[string]schema.Attribute{
            "host": schema.StringAttribute{
                Optional: true,
            },
        },
    }
}

// Configure prepares a devops API client for data sources and resources.
func (p *devopsProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
    // log info for debugging
    tflog.Info(ctx, "Configuring DevOps client")

    // Retrieve provider data from configuration
    var config devopsProviderModel
    diags := req.Config.Get(ctx, &config)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    if config.Host.IsUnknown() {
        resp.Diagnostics.AddAttributeError(
            path.Root("host"),
            "Unknown DevOps API Host",
            "The provider cannot create the HashiCups API client as there is an unknown configuration value for the DevOps API host. Fix it",
        )
    }

    if resp.Diagnostics.HasError() {
        return
    }


    // Default values to environment varaibles,
    // but override with Terraform configuration value if set

    host:= os.Getenv("DEVOPS_HOST")

    if !config.Host.IsNull() {
        host = config.Host.ValueString()
    }

    if host == "" {
        resp.Diagnostics.AddAttributeError(
            path.Root("host"),
            "Missing DevOpsAPI Host",
            "The provider cannot create the DevOps API client as there is a missing or empty value for the DevOps API host. "+
                "Set the host value in the configuration or use the DEVOPS_HOST environment variable. "+
                "If either is already set, ensure the value is not empty.",
        )
    }

    if resp.Diagnostics.HasError() {
        return
    }

    ctx = tflog.SetField(ctx, "devops_host", config.Host.ValueString())

    tflog.Debug(ctx, "Creating DevOps Client")

    // Create a new DevOps client using the configuration values
    client, err := NewClient(&host)
    if err != nil {
        resp.Diagnostics.AddError(
            "Unable to Create DevOps API Client",
            fmt.Sprintf("Unable to create the DevOps API client, got error: %s", err),
        )
        return
    }

    // Make the DevOps client available during DataSource and Resource
    // type Configure methods.

    resp.DataSourceData = client
    resp.ResourceData = client

    tflog.Info(ctx, "Configured DevOps Client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *devopsProvider) DataSources(_ context.Context) []func() datasource.DataSource {
    return []func() datasource.DataSource {
        NewEngineerDataSource,
    }
}

// Resources defines the resources implemented in the provider.
func (p *devopsProvider) Resources(_ context.Context) []func() resource.Resource {
    return []func() resource.Resource {
        NewEngineerResource,
    }
}
