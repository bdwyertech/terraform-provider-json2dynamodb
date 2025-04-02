package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure JSON2DynamoDBProvider satisfies various provider interfaces.
var _ provider.Provider = &JSON2DynamoDBProvider{}
var _ provider.ProviderWithFunctions = &JSON2DynamoDBProvider{}
var _ provider.ProviderWithEphemeralResources = &JSON2DynamoDBProvider{}

// JSON2DynamoDBProvider defines the provider implementation.
type JSON2DynamoDBProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// JSON2DynamoDBProviderModel describes the provider data model.
type JSON2DynamoDBProviderModel struct {
	// Endpoint types.String `tfsdk:"endpoint"`
}

func (p *JSON2DynamoDBProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "json2dynamodb"
	resp.Version = p.version
}

func (p *JSON2DynamoDBProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			// "endpoint": schema.StringAttribute{
			// 	MarkdownDescription: "Example provider attribute",
			// 	Optional:            true,
			// },
		},
	}
}

func (p *JSON2DynamoDBProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data JSON2DynamoDBProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	// if data.Endpoint.IsNull() { /* ... */ }

	// Example client configuration for data sources and resources
	client := http.DefaultClient
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *JSON2DynamoDBProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		// NewExampleResource,
	}
}

func (p *JSON2DynamoDBProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{
		// NewExampleEphemeralResource,
	}
}

func (p *JSON2DynamoDBProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewJSON2DynamoDBDataSource,
	}
}

func (p *JSON2DynamoDBProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{
		// NewExampleFunction,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &JSON2DynamoDBProvider{
			version: version,
		}
	}
}
