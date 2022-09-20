package omdb

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"net/http"
)

const (
	urlFilmById = "https://www.omdbapi.com/?i=%s&apikey=%s"
)

// filmByIdApiResponse defines what we expect from urlFilmById
type filmByIdApiResponse struct {
	ImdbID string `json:"imdbID"`
	Title  string `json:"Title"`
	Year   string `json:"Year"`
}

// filmByIdData is a terraform config/plan/state style object
type filmByIdData struct {
	ImdbId types.String `tfsdk:"imdb_id"`
	Title  types.String `tfsdk:"title"`
	Year   types.String `tfsdk:"year"`
}

var _ datasource.DataSource = &FilmByIdDataSource{}

// FilmByIdDataSource implements the datasource.DataSourceWithConfigure interface
type FilmByIdDataSource struct {
	apiKey     string
	configured bool
}

func (d *FilmByIdDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_film_by_id"
}

func (d *FilmByIdDataSource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "This Data Source returns details about a film by its IMDb ID.",
		Attributes: map[string]tfsdk.Attribute{
			"imdb_id": {
				MarkdownDescription: "Unique ID used by both OMDb and IMDb.",
				Required:            true,
				Type:                types.StringType,
			},
			"title": {
				MarkdownDescription: "Film title.",
				Computed:            true,
				Type:                types.StringType,
			},
			"year": {
				MarkdownDescription: "Release year.",
				Computed:            true,
				Type:                types.StringType,
			},
		},
	}, diag.Diagnostics{}
}

func (d *FilmByIdDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if providerData, ok := req.ProviderData.(*providerData); ok {
		if providerData.configured {
			d.apiKey = providerData.apiKey
			d.configured = true
		}
	}
}

func (d *FilmByIdDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if !d.configured {
		resp.Diagnostics.AddError("data source Read() method called prior to Configure()", "don't do that")
		return
	}
	var config filmByIdData
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResponse, err := http.Get(fmt.Sprintf(urlFilmById, config.ImdbId.Value, d.apiKey))
	if err != nil {
		resp.Diagnostics.AddError("error making http request", err.Error())
		return
	}

	var apiResponse filmByIdApiResponse
	err = json.NewDecoder(httpResponse.Body).Decode(&apiResponse)
	if err != nil {
		resp.Diagnostics.AddError("error decoding API response", err.Error())
	}

	state := filmByIdData{
		ImdbId: types.String{Value: config.ImdbId.Value},
		Title:  types.String{Value: apiResponse.Title},
		Year:   types.String{Value: apiResponse.Year},
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
