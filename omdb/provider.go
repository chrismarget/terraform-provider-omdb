package omdb

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Provider fulfils the provider.Provider interface
type Provider struct {
	Version string
	Commit  string
}

// dataSourceProviderData gets instantiated in the provider's Configure()
// method and is made available to methods on datasource.DataSource
type dataSourceProviderData struct {
	apiKey     string
	configured bool
}

// resourceProviderData gets instantiated in the provider's Configure()
// method and is made available to methods on resource.Resource
type resourceProviderData struct {
	// no resources provisioned in this provider
}

func (p *Provider) Metadata(_ context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "omdb"
	if p.Version != "" {
		resp.Version = "v" + p.Version
	} else {
		resp.Version = p.Commit
	}
}

// GetSchema returns provider schema
func (p *Provider) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Top level provider markdown description.",
		Attributes: map[string]tfsdk.Attribute{
			"api_key": {
				MarkdownDescription: "A free OMDb API key can be quickly generated [here](https://www.omdbapi.com/apikey.aspx).",
				Type:                types.StringType,
				Required:            true,
			},
		},
	}, diag.Diagnostics{}
}

// Provider configuration struct
type providerConfig struct {
	apiKey types.String `tfsdk:"api_key"`
}

// Configure is supposed to run before any DataSource.Configure() or
// Resource.Configure(), but I'm not sure it's happening.
func (p *Provider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config providerConfig
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// data available to data source Configure() method
	resp.DataSourceData = dataSourceProviderData{
		apiKey:     config.apiKey.Value, // not checking null/unknown because required by schema
		configured: true,
	}

	// data available to resource Configure() method
	resp.ResourceData = dataSourceProviderData{
		apiKey:     config.apiKey.Value, // not checking null/unknown because required by schema
		configured: true,
	}

	resp.Diagnostics.AddWarning("provider configure",
		fmt.Sprintf("provider is: '%v'", resp.DataSourceData))
}

// Resources defines provider resources by returning a slice of functions
// which return resources.
func (p *Provider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		// no resources
	}
}

// DataSources defines provider data sources by returning a slice of functions
// which return data sources.
func (p *Provider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		func() datasource.DataSource { return &FilmByIdDataSource{} },
	}
}
