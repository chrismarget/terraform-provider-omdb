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
	Title  types.String `tfsdk:"imdb_id"`
	Year   types.String `tfsdk:"year"`
}

var _ datasource.DataSource = &FilmByIdDataSource{}

// FilmByIdDataSource implements the datasource.DataSourceWithConfigure interface
type FilmByIdDataSource struct {
	apiKey string
}

func (d *FilmByIdDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_film_by_id"
}

func (d *FilmByIdDataSource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"imdb_id": {
				Required: true,
				Type:     types.StringType,
			},
			"title": {
				Computed: true,
				Type:     types.StringType,
			},
			"Year": {
				Computed: true,
				Type:     types.StringType,
			},
		},
	}, nil
}

func (d *FilmByIdDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	providerData, ok := req.ProviderData.(*dataSourceProviderData)
	if !ok {
		resp.Diagnostics.AddError("invalid configure data",
			fmt.Sprintf("unable to type assert ProviderData to '%T'", dataSourceProviderData{}))
		return
	}
	if providerData == nil || !providerData.configured {
		resp.Diagnostics.AddError("provider not configured",
			"either the providerData is nil, or the configured flag is unset")
		return
	}
	d.apiKey = providerData.apiKey
}

func (d *FilmByIdDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config filmByIdData
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResponse, err := http.Get(fmt.Sprintf(urlFilmById, config.ImdbId, d.apiKey))
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
