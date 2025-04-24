package provider

import (
        "context"
        "fmt"

        "github.com/hashicorp/terraform-plugin-framework/datasource"
        "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
        "github.com/hashicorp/terraform-plugin-framework/path"
        // "github.com/hashicorp/terraform-plugin-framework/diag"
)

// Ensure the implementation satisfies the expected interfaces.
var (
        _ datasource.DataSource = &engineerDataSource{}
        _ datasource.DataSourceWithConfigure = &engineerDataSource{}
)

// helper function to simplify the provider implementation
func NewEngineerDataSource() datasource.DataSource {
        return &engineerDataSource{}
}

// engineerDataSource defines the data source implementation.
type engineerDataSource struct{
        client *Client
}

// engineerDataSourceModel defines the data model for the data source.
type engineerDataSourceModel struct {
        //Engineer engineerModel `tfsdk:"engineer"`
        Name    string `tfsdk:"name"`
        Id      string `tfsdk:"id"`
        Email   string `tfsdk:"email"`
}

// engineerModel maps engineer schema data
// type engineerModel struct {
//         Name    string `tfsdk:"name"`
//         Id      string `tfsdk:"id"`
//         Email   string `tfsdk:"email"`
// }

// Metadata returns the data source type name.
func (d *engineerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
        resp.TypeName = req.ProviderTypeName + "_engineer"
}


// Configure adds the provider configured client to the data source.
func (d *engineerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
        // Add a nil check when handling ProviderData because Terraform
        // sets that data after it calls the ConfigureProvider RPC.
        if req.ProviderData == nil {
                return
        }

        client, ok := req.ProviderData.(*Client)
        if !ok {
                resp.Diagnostics.AddError(
                        "Unexpected Data Source Configure Type",
                        fmt.Sprintf("Expected Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
                        )

                return
        }

        d.client = client
}


// Schema defines the schema for the data source.
func (d *engineerDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
        resp.Schema = schema.Schema{
                Attributes: map[string]schema.Attribute{
                        "name": schema.StringAttribute{
                                Required: true,    
                        },
                        "id": schema.StringAttribute{
                                Computed: true,
                        },
                        "email": schema.StringAttribute{
                                Required: true,
                        },
                },
        }
}

// Read refreshes the Terraform state with the latest data.
func (d *engineerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
        var state engineerDataSourceModel

        // Get the "id" attribute from the configuration
        var id string
        diags := req.Config.GetAttribute(ctx, path.Root("id"), &id)
        if diags.HasError() {
                resp.Diagnostics.AddError(
                        "Unable to Read DevOps engineer",
                        "The 'id' attribute is missing or invalid",
                        )
                return
        }

        // // Convert attribute to a string
        // id := idAttr.(string)

        engineer, err := d.client.GetEngineer(id)
        if err != nil {
                resp.Diagnostics.AddError(
                        "Unable to Read DevOps engineer",
                        err.Error(),
                        )
                return
        }

        // Map the engineer data to the state model
        state.Name =    engineer.Name
        state.Id =      engineer.Id
        state.Email =   engineer.Email

        // Set state
        diags = resp.State.Set(ctx, &state)
        resp.Diagnostics.Append(diags...)
        if resp.Diagnostics.HasError() {
                return
        }
}
