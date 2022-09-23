package omdb

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
)

// filmFileData defines what we expect to find in the files in baseDir
type filmFileData struct {
	Title   string `json:"Title"`
	Year    string `json:"Year"`
	Ratings []struct {
		Source string `json:"Source"`
		Value  string `json:"Value"`
	} `json:"Ratings,omitempty"`
}

// filmByIdData is a terraform config/plan/state style object
type filmData struct {
	Id       types.String     `tfsdk:"id"`
	Title    types.String     `tfsdk:"title"`
	Year     types.String     `tfsdk:"year"`
	Ratings0 []filmRatingData `tfsdk:"ratings0"`
	Ratings1 []types.Object   `tfsdk:"ratings1"`
	Ratings2 types.List       `tfsdk:"ratings2"`
}

var _ resource.Resource = &ResourceFilm{}
var _ resource.ResourceWithConfigure = &ResourceFilm{}

// ResourceFilm implements the datasource.DataSourceWithConfigure interface
type ResourceFilm struct {
	localDir string
}

func (r *ResourceFilm) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_film"
}

func (r *ResourceFilm) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	if providerData, ok := req.ProviderData.(*providerResourceData); ok {
		r.localDir = providerData.localDir
	} else {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type",
			fmt.Sprintf("Expected '%T', got: '%T'. Please report this issue to the provider developers", providerData, req.ProviderData))
	}
}

func (r *ResourceFilm) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "This Data Source returns details about a film by its IMDb ID.",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "Unique ID",
				Computed:            true,
				Type:                types.StringType,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.RequiresReplace(),
					resource.UseStateForUnknown(),
				},
			},
			"title": {
				MarkdownDescription: "Film title",
				Required:            true,
				Type:                types.StringType,
			},
			"year": {
				MarkdownDescription: "Release year",
				Required:            true,
				Type:                types.StringType,
			},
			"ratings0": {
				MarkdownDescription: "Ratings0",
				Optional:            true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"source": {
						MarkdownDescription: "Review source",
						Optional:            true,
						Type:                types.StringType,
					},
					"value": {
						MarkdownDescription: "Review value",
						Optional:            true,
						Type:                types.StringType,
					},
				}),
			},
			"ratings1": {
				MarkdownDescription: "Ratings1",
				Optional:            true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"source": {
						MarkdownDescription: "Review source",
						Optional:            true,
						Type:                types.StringType,
					},
					"value": {
						MarkdownDescription: "Review value",
						Optional:            true,
						Type:                types.StringType,
					},
				}),
			},
			"ratings2": {
				MarkdownDescription: "Ratings2",
				Optional:            true,
				Type: types.ListType{
					ElemType: types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"source": types.StringType,
							"value":  types.StringType,
						},
					},
				},
			},
		},
	}, diag.Diagnostics{}
}

func (r *ResourceFilm) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan filmData
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	film := &filmFileData{
		Title: plan.Title.Value,
		Year:  plan.Year.Value,
	}

	film.Ratings = make([]struct {
		Source string `json:"Source"`
		Value  string `json:"Value"`
	}, len(plan.Ratings0))

	for i, rating := range plan.Ratings0 {
		film.Ratings[i] = struct {
			Source string `json:"Source"`
			Value  string `json:"Value"`
		}{
			Source: rating.Source.Value,
			Value:  rating.Value.Value,
		}
	}

	data, err := json.MarshalIndent(film, "", "  ")
	if err != nil {
		resp.Diagnostics.AddError("error marshaling film to JSON", err.Error())
		return
	}

	b := make([]byte, 8)
	rand.Read(b)
	plan.Id = types.String{Value: fmt.Sprintf("%x", b)}

	err = ioutil.WriteFile(filepath.Join(r.localDir, plan.Id.Value), data, 0644)
	if err != nil {
		resp.Diagnostics.AddError("error writing film to file", err.Error())
		return
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (r *ResourceFilm) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state filmData
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	file, err := os.Open(filepath.Join(r.localDir, state.Id.Value))
	if err != nil {
		if os.IsNotExist(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("error opening file", err.Error())
	}
	defer func(file *os.File) { _ = file.Close() }(file)

	var film filmFileData
	err = json.NewDecoder(file).Decode(&film)
	if err != nil {
		resp.Diagnostics.AddError("error reading/parsing file", err.Error())
	}

	newState := filmData{
		Id:    types.String{Value: state.Id.Value},
		Title: types.String{Value: film.Title},
		Year:  types.String{Value: film.Year},
		Ratings2: types.List{
			Elems: make([]attr.Value, len(film.Ratings)),
			ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"source": types.StringType,
					"value":  types.StringType,
				},
			},
		},
	}
	for i, rating := range film.Ratings {
		newState.Ratings2.Elems[i] = types.Object{
			Attrs: map[string]attr.Value{
				"source": types.String{Value: rating.Source},
				"value":  types.String{Value: rating.Value},
			},
			AttrTypes: map[string]attr.Type{
				"source": types.StringType,
				"value":  types.StringType,
			},
		}
	}
	if len(film.Ratings) == 0 {
		newState.Ratings2.Null = true
	}

	//o, _ := json.Marshal(state)
	//n, _ := json.Marshal(newState)
	//resp.Diagnostics.AddWarning("old", string(o))
	//resp.Diagnostics.AddWarning("new", string(n))

	diags = resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
}

func (r *ResourceFilm) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state filmData
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var plan filmData
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.Id = types.String{Value: state.Id.Value}

	film := &filmFileData{
		Title: plan.Title.Value,
		Year:  plan.Year.Value,
	}

	film.Ratings = make([]struct {
		Source string `json:"Source"`
		Value  string `json:"Value"`
	}, len(plan.Ratings0))

	for i, rating := range plan.Ratings0 {
		film.Ratings[i] = struct {
			Source string `json:"Source"`
			Value  string `json:"Value"`
		}{
			Source: rating.Source.Value,
			Value:  rating.Value.Value,
		}
	}

	data, err := json.MarshalIndent(film, "", "  ")
	if err != nil {
		resp.Diagnostics.AddError("error marshaling film to JSON", err.Error())
		return
	}

	err = ioutil.WriteFile(filepath.Join(r.localDir, plan.Id.Value), data, 0644)
	if err != nil {
		resp.Diagnostics.AddError("error writing film to file", err.Error())
		return
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (r *ResourceFilm) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state filmData
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.Id.IsNull() || state.Id.IsUnknown() || state.Id.Value == "" {
		resp.Diagnostics.AddError("delete error", "cannot delete film with unknown ID")
		return
	}

	fileName := filepath.Join(r.localDir, state.Id.Value)
	err := os.Remove(fileName)
	if err != nil {
		resp.Diagnostics.AddError("delete error", err.Error())
	}
}
