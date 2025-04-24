package provider

import (
	"context"
	"fmt"
	// "strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource		   = &engineerResource{}
	_ resource.ResourceWithConfigure   = &engineerResource{}
	_ resource.ResourceWithImportState = &engineerResource{}
)

// NewengineerResource is a helper function to simplify the provider implementation.
func NewEngineerResource() resource.Resource {
	return &engineerResource{}
}

// engineerResource is the resource implementation.
type engineerResource struct{
	client *Client
}

// engineerResourceModel maps the resource schema data.
type engineerResourceModel struct{
	Id types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Email types.String `tfsdk:"email"`
	LastUpdated types.String `tfsdk:"last_updated"`
}

// Metadata returns the resource type name.
func (r *engineerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_engineer"
}

// Schema defines the schema for the resource.
func (r *engineerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
			},
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"email": schema.StringAttribute{
				Required: true,
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *engineerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from new_engineer
	var plan engineerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}



	// // Generate API request body from plan
	// var items engineerResourceModel
	// for _, item := range engineerResourceModel{
	// 	items = append(items, hashicups.OrderItem{
	// 		Coffee: hashicups.Coffee{
	// 			Id: int(item.Coffee.Id.ValueInt64()),
	// 		},
	// 		Quantity: int(item.Quantity.ValueInt64()),
	// 	})
	// }

	// Create new engineer
	engineer, err := r.client.CreateEngineer(plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating engineer",
			"Could not create engineer, unexpected error: "+err.Error(),
			)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.Name = types.StringValue(engineer.Name) //engineer.Name.
	plan.Id = types.StringValue(engineer.Id) //engineer.Id
	plan.Email = types.StringValue(engineer.Email) //engineer.Email
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *engineerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state engineerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Fetch engineer by Id
	engineer, err := r.client.GetEngineer(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading engineer",
			"Could not read engineer, unexpected error: "+err.Error(),
			)
		return
	}

	// overwrite state with data from API
	state.Name = types.StringValue(engineer.Name) //engineer.Name
	state.Id = types.StringValue(engineer.Id) //engineer.Id
	state.Email = types.StringValue(engineer.Email) //engineer.Email

	// set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *engineerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan engineerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// // Generate API request body from plan
	// var hashicupsItems []hashicups.OrderItem
	// for _, item := range plan.Items {
	//     hashicupsItems = append(hashicupsItems, hashicups.OrderItem{
	//         Coffee: hashicups.Coffee{
	//             Id: int(item.Coffee.Id.ValueInt64()),
	//         },
	//         Quantity: int(item.Quantity.ValueInt64()),
	//     })
	// }


	// Update existing order
	_, err := r.client.UpdateEngineer(plan.Id.ValueString(), plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Engineer",
			"Could not update engineer, unexpected error: "+err.Error(),
			)
		return
	}

	tflog.Info(ctx, "Engineer updated", map[string]interface{}{
                "Id": plan.Id.ValueString(),
        })

	// Fetch updated items from Engineer
	updatedEngineer, err := r.client.GetEngineer(plan.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Engineer",
			"Could not read Engineer Id "+plan.Name.ValueString()+": "+err.Error(),
			)
		return
	}

	// Update resource state with updated items and timestamp
	plan.Id = types.StringValue(updatedEngineer.Id)
	plan.Name = types.StringValue(updatedEngineer.Name)
	plan.Email = types.StringValue(updatedEngineer.Email)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *engineerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
        // Retrieve values from state
        var state engineerResourceModel
        diags := req.State.Get(ctx, &state)
        resp.Diagnostics.Append(diags...)
        if resp.Diagnostics.HasError() {
                return
        }

        // Delete existing order
        err := r.client.DeleteEngineer(state.Id.ValueString())
        if err != nil {
                resp.Diagnostics.AddError(
                        "Error deleting engineer",
                        "Could not delete engineer, unexpected error: "+err.Error(),
                        )
                return
        }
}

func (r *engineerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nill check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
			)

		return
	}

	r.client = client
}


func (r *engineerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    // Retrieve import ID and save to id attribute
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
