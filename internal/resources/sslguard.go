package resources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"gitlab.wpj.cz/terraform-provider-nethost/internal/client"
)

func NewSSLGuardResource() resource.Resource {
	return &SSLGuardResource{}
}

type SSLGuardResource struct {
	client *client.Client
}

type SSLGuardResourceModel struct {
	ID     types.String `tfsdk:"id"`
	Domain types.String `tfsdk:"domain"`
	Port   types.Int64  `tfsdk:"port"`
}

func (r *SSLGuardResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sslguard"
}

func (r *SSLGuardResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"domain": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"port": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(443),
				Validators: []validator.Int64{
					int64validator.Between(1, 65535),
				},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *SSLGuardResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SSLGuardResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SSLGuardResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.CreateSSLGuard(ctx, plan.Domain.ValueString(), plan.Port.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Create SSLGuard failed", err.Error())
		return
	}

	item, err := r.client.FindSSLGuardByDomain(ctx, plan.Domain.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Read SSLGuard after create failed", err.Error())
		return
	}

	if item == nil || item.ID == "" {
		resp.Diagnostics.AddError(
			"SSLGuard ID not found",
			"SSLGuard was created, but API /search did not return matching domain.",
		)
		return
	}

	plan.ID = types.StringValue(item.ID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SSLGuardResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SSLGuardResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := r.client.FindSSLGuardByDomain(ctx, state.Domain.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Read SSLGuard failed", err.Error())
		return
	}

	if item == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.ID = types.StringValue(item.ID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *SSLGuardResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update SSLGuard is not supported",
		"Changing domain or port requires replacing the resource.",
	)
}

func (r *SSLGuardResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state SSLGuardResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.ID.IsNull() || state.ID.ValueString() == "" {
		resp.Diagnostics.AddError("Missing SSLGuard ID", "Cannot delete SSLGuard without ID in Terraform state.")
		return
	}

	err := r.client.DeleteSSLGuard(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Delete SSLGuard failed", err.Error())
		return
	}
}
