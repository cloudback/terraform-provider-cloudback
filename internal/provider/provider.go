package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure CloudbackProvider satisfies various provider interfaces.
var _ provider.Provider = &CloudbackProvider{}

// CloudbackProvider defines the provider implementation.
type CloudbackProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// CloudbackProviderModel describes the provider data model.
type CloudbackProviderModel struct {
	ApiKey   types.String `tfsdk:"api_key"`
	Endpoint types.String `tfsdk:"endpoint"`
}

func (p *CloudbackProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "cloudback"
	resp.Version = p.version
}

func (p *CloudbackProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The Cloudback provider allows you to manage backup definitions.",
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				MarkdownDescription: "The API key for authentication. May also be provided via CLOUDBACK_API_KEY environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "The API endpoint URL. May also be provided via CLOUDBACK_ENDPOINT environment variable. Default is https://app.cloudback.it.",
				Required:            false,
				Optional:            true,
			},
		},
	}
}

func (p *CloudbackProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Check environment variables
	apiKey := os.Getenv("CLOUDBACK_API_KEY")
	endpoint := os.Getenv("CLOUDBACK_ENDPOINT")

	var data CloudbackProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Check configuration data, which should take precedence over
	// environment variable data, if found.
	if data.ApiKey.ValueString() != "" {
		apiKey = data.ApiKey.ValueString()
	}

	if data.Endpoint.ValueString() != "" {
		endpoint = data.Endpoint.ValueString()
	}

	if apiKey == "" {
		resp.Diagnostics.AddError(
			"Missing API Key Configuration",
			"While configuring the provider, the API key was not found in "+
				"the CLOUDBACK_API_KEY environment variable or provider "+
				"configuration block api_key attribute.",
		)
		// Not returning early allows the logic to collect all errors.
	}

	if endpoint == "" {
		endpoint = "https://app.cloudback.it"
	}

	// Create data/clients and persist to resp.DataSourceData, resp.ResourceData,
	client := NewCloudbackClient(endpoint, apiKey)
	resp.ResourceData = client
}

func (p *CloudbackProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewBackupDefinitionResource,
	}
}

func (p *CloudbackProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return nil
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &CloudbackProvider{
			version: version,
		}
	}
}
