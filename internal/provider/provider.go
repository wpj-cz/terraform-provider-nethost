package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"gitlab.wpj.cz/terraform-provider-nethost/internal/client"
	"gitlab.wpj.cz/terraform-provider-nethost/internal/resources"
)

func NewProvider() provider.Provider {
	return &NethostProvider{}
}

type NethostProvider struct{}

type NethostProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	APIKey   types.String `tfsdk:"api_key"`
}

func (p *NethostProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "nethost"
}

func (p *NethostProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				Required:    true,
				Description: "The URL endpoint for the Nethost API.",
			},
			"api_key": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "The API key for Nethost authentication.",
			},
		},
	}
}

func (p *NethostProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config NethostProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// 1. Endpoint validation
	var endpoint string
	if config.Endpoint.IsNull() || config.Endpoint.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Missing Nethost Endpoint",
			"The 'endpoint' attribute is required in the provider configuration block.",
		)
		return
	}

	endpoint = config.Endpoint.ValueString()

	// 2. API key validation
	apiKey := os.Getenv("NETHOST_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("API_KEY")
	}
	if !config.APIKey.IsNull() && !config.APIKey.IsUnknown() {
		apiKey = config.APIKey.ValueString()
	}

	if apiKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Missing Nethost API key",
			"Set api_key in provider config or use NETHOST_API_KEY/API_KEY environment variable.",
		)
		return
	}

	// 3. Initializace client
	c := client.NewClient(endpoint, apiKey)

	resp.DataSourceData = c
	resp.ResourceData = c
}

func (p *NethostProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		resources.NewSSLGuardResource,
		resources.NewSpongefingerResource,
	}
}

func (p *NethostProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return nil
}
