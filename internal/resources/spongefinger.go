package resources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"gitlab.wpj.cz/terraform-provider-nethost/internal/client"
	"gitlab.wpj.cz/terraform-provider-nethost/internal/tfconvert"
)

func NewSpongefingerResource() resource.Resource {
	return &SpongefingerResource{}
}

type SpongefingerResource struct {
	client *client.Client
}

type SpongefingerResourceModel struct {
	ID         types.String   `tfsdk:"id"`
	Name       types.String   `tfsdk:"name"`
	VariantID  types.Int64    `tfsdk:"variant_id"`
	Domain     types.String   `tfsdk:"domain"`
	Subdomains []types.String `tfsdk:"subdomains"`
	Backends   []types.String `tfsdk:"backends"`
}

func (r *SpongefingerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_spongefinger"
}

func (r *SpongefingerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"variant_id": schema.Int64Attribute{
				Required:    true,
				Description: "Variants (plans)",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"domain": schema.StringAttribute{
				Required:    true,
				Description: "Protected domain. Enter the main (root) domain you want to protect.",
			},
			"subdomains": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Protected subdomains",
			},
			"backends": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
				Description: "Target servers. The origin server address in `host:port` format.",
			},
		},
	}
}

func (r *SpongefingerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected provider data",
			fmt.Sprintf("Expected *Client, got: %T", req.ProviderData),
		)
		return
	}

	r.client = c
}

func (r *SpongefingerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SpongefingerResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.CreateSpongefinger(
		ctx,
		plan.Name.ValueString(),
		plan.VariantID.ValueInt64(),
		plan.Domain.ValueString(),
		tfconvert.TypesToStrings(plan.Subdomains),
		tfconvert.TypesToStrings(plan.Backends),
	)
	if err != nil {
		resp.Diagnostics.AddError("Create Spongefinger failed", err.Error())
		return
	}

	item, err := r.client.FindSpongefingerByDomain(ctx, plan.Domain.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Read Spongefinger after create failed", err.Error())
		return
	}

	if item == nil || item.ID == "" {
		resp.Diagnostics.AddError(
			"Spongefinger ID not found",
			"Spongefinger was created, but API /search did not return matching domain.",
		)
		return
	}

	detail, err := r.client.FindSpongefingerByID(ctx, item.ID)
	if err != nil {
		resp.Diagnostics.AddError("Read Spongefinger after create failed", err.Error())
		return
	}

	if detail == nil {
		resp.Diagnostics.AddError(
			"Spongefinger detail not found",
			"Spongefinger was created, but API /detail did not return matching ID.",
		)
		return
	}

	plan.ID = types.StringValue(detail.ID)
	plan.Name = types.StringValue(detail.Name)
	plan.Domain = types.StringValue(detail.Domain)
	plan.VariantID = types.Int64Value(detail.VariantID)
	plan.Subdomains = tfconvert.StringsToTypes(detail.Subdomains)
	plan.Backends = tfconvert.StringsToTypes(detail.Backends)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SpongefingerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SpongefingerResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	if id == "" {
		item, err := r.client.FindSpongefingerByDomain(ctx, state.Domain.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Read Spongefinger failed", err.Error())
			return
		}
		if item == nil {
			resp.State.RemoveResource(ctx)
			return
		}
		id = item.ID
	}

	item, err := r.client.FindSpongefingerByID(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError("Read Spongefinger failed", err.Error())
		return
	}

	if item == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.ID = types.StringValue(item.ID)
	state.Name = types.StringValue(item.Name)
	state.Domain = types.StringValue(item.Domain)
	state.VariantID = types.Int64Value(item.VariantID)
	state.Subdomains = tfconvert.StringsToTypes(item.Subdomains)
	state.Backends = tfconvert.StringsToTypes(item.Backends)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *SpongefingerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan SpongefingerResourceModel
	var state SpongefingerResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !state.VariantID.IsNull() && !plan.VariantID.IsNull() && state.VariantID.ValueInt64() != plan.VariantID.ValueInt64() {
		resp.Diagnostics.AddError(
			"Update Spongefinger requires replacement",
			"Changing variant_id is not supported by the update endpoint.",
		)
		return
	}

	id := state.ID.ValueString()
	if id == "" {
		resp.Diagnostics.AddError("Missing Spongefinger ID", "Cannot update Spongefinger without ID in Terraform state.")
		return
	}

	err := r.client.UpdateSpongefinger(
		ctx,
		id,
		plan.Name.ValueString(),
		plan.Domain.ValueString(),
		tfconvert.TypesToStrings(plan.Subdomains),
		tfconvert.TypesToStrings(plan.Backends),
	)
	if err != nil {
		resp.Diagnostics.AddError("Update Spongefinger failed", err.Error())
		return
	}

	detail, err := r.client.FindSpongefingerByID(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError("Read Spongefinger after update failed", err.Error())
		return
	}

	if detail == nil {
		resp.Diagnostics.AddError(
			"Spongefinger detail not found",
			"Spongefinger was updated, but API /detail did not return matching ID.",
		)
		return
	}

	plan.ID = types.StringValue(detail.ID)
	plan.Name = types.StringValue(detail.Name)
	plan.Domain = types.StringValue(detail.Domain)
	plan.VariantID = types.Int64Value(detail.VariantID)
	plan.Subdomains = tfconvert.StringsToTypes(detail.Subdomains)
	plan.Backends = tfconvert.StringsToTypes(detail.Backends)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SpongefingerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state SpongefingerResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.ID.IsNull() || state.ID.ValueString() == "" {
		resp.Diagnostics.AddError("Missing Spongefinger ID", "Cannot delete Spongefinger without ID in Terraform state.")
		return
	}

	err := r.client.DeleteSpongefinger(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Delete Spongefinger failed", err.Error())
		return
	}
}
