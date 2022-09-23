package omdb

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"math/rand"
	"os"
	"time"
)

const (
	defaultBaseUrl  = "https://www.omdbapi.com"
	urlFilmById     = "/?i=%s&apikey=%s"
	defaultLocalDir = "/tmp/.omdb"
)

var _ provider.ProviderWithMetadata = &Provider{}

// Provider fulfils the provider.Provider interface
type Provider struct {
	Version string // populated in main() using value set by the linker
	Commit  string // populated in main() using value set by the linker
}

// providerDataSourceData gets instantiated in the provider.Provider's
// Configure() method and is made available to the Configure() method of
// implementations of datasource.DataSource
type providerDataSourceData struct {
	apiBaseUrl string
	apiKey     string
}

// providerResourceData gets instantiated in the provider.Provider's
// Configure() method and is made available to the Configure() method of
// implementations of resource.Resource
type providerResourceData struct {
	localDir string
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
				Validators:          []tfsdk.AttributeValidator{stringvalidator.LengthAtLeast(1)},
			},
			"api_url": {
				MarkdownDescription: "URL of the OMDb service, defaults to " + defaultBaseUrl,
				Type:                types.StringType,
				Optional:            true,
				Validators:          []tfsdk.AttributeValidator{stringvalidator.LengthAtLeast(1)},
			},
			"local_dir": {
				MarkdownDescription: "The local directory where film \"resources\" are created, defaults to" + defaultLocalDir,
				Type:                types.StringType,
				Optional:            true,
				Validators:          []tfsdk.AttributeValidator{stringvalidator.LengthAtLeast(1)},
			},
		},
	}, diag.Diagnostics{}
}

// Provider configuration struct
type providerConfig struct {
	ApiKey   types.String `tfsdk:"api_key"`
	ApiUrl   types.String `tfsdk:"api_url"`
	LocalDir types.String `tfsdk:"local_dir"`
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

	if config.ApiUrl.Null {
		config.ApiUrl = types.String{Value: defaultBaseUrl}
	}

	if config.LocalDir.Null {
		config.LocalDir = types.String{Value: defaultLocalDir}
	}

	err := os.MkdirAll(config.LocalDir.Value, 0755)
	if err != nil {
		resp.Diagnostics.AddError("error creating local directory", err.Error())
	}

	// data we intend to make available to the Configure() method of
	// implementations of datasource.DataSource
	resp.DataSourceData = &providerDataSourceData{
		apiBaseUrl: config.ApiUrl.Value,
		apiKey:     config.ApiKey.Value,
	}

	// data we intend to make available to the Configure() method of
	// implementations of resource.Resource.
	resp.ResourceData = &providerResourceData{
		localDir: config.LocalDir.Value,
	}

	rand.Seed(time.Now().UnixNano())
}

// DataSources defines provider data sources by returning a slice of functions
// which return data sources.
func (p *Provider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		func() datasource.DataSource { return &DataSourceFilmById{} },
	}
}

// Resources defines provider resources by returning a slice of functions
// which return resources.
func (p *Provider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		func() resource.Resource { return &ResourceFilm{} },
	}
}
