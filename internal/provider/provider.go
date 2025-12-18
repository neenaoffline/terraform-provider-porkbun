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

// Ensure PorkbunProvider satisfies various provider interfaces.
var _ provider.Provider = &PorkbunProvider{}

// PorkbunProvider defines the provider implementation.
type PorkbunProvider struct {
	version string
}

// PorkbunProviderModel describes the provider data model.
type PorkbunProviderModel struct {
	APIKey       types.String `tfsdk:"api_key"`
	SecretAPIKey types.String `tfsdk:"secret_api_key"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &PorkbunProvider{
			version: version,
		}
	}
}

func (p *PorkbunProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "porkbun"
	resp.Version = p.version
}

func (p *PorkbunProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with Porkbun DNS API.",
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				Description: "Porkbun API key. Can also be set via the PORKBUN_API_KEY environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
			"secret_api_key": schema.StringAttribute{
				Description: "Porkbun Secret API key. Can also be set via the PORKBUN_SECRET_API_KEY environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

func (p *PorkbunProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config PorkbunProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get API key from config or environment
	apiKey := os.Getenv("PORKBUN_API_KEY")
	if !config.APIKey.IsNull() {
		apiKey = config.APIKey.ValueString()
	}

	// Get Secret API key from config or environment
	secretAPIKey := os.Getenv("PORKBUN_SECRET_API_KEY")
	if !config.SecretAPIKey.IsNull() {
		secretAPIKey = config.SecretAPIKey.ValueString()
	}

	if apiKey == "" {
		resp.Diagnostics.AddError(
			"Missing API Key",
			"The provider cannot create the Porkbun API client as there is a missing or empty value for the Porkbun API key. "+
				"Set the api_key value in the configuration or use the PORKBUN_API_KEY environment variable.",
		)
	}

	if secretAPIKey == "" {
		resp.Diagnostics.AddError(
			"Missing Secret API Key",
			"The provider cannot create the Porkbun API client as there is a missing or empty value for the Porkbun Secret API key. "+
				"Set the secret_api_key value in the configuration or use the PORKBUN_SECRET_API_KEY environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create a new Porkbun client
	client := NewClient(apiKey, secretAPIKey)

	// Test the connection
	if err := client.Ping(); err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Porkbun API Client",
			"An unexpected error occurred when creating the Porkbun API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Porkbun Client Error: "+err.Error(),
		)
		return
	}

	// Make the client available during DataSource and Resource type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *PorkbunProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewDNSRecordResource,
	}
}

func (p *PorkbunProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewDNSRecordDataSource,
	}
}
