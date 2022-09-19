package main

import (
	"context"
	"github.com/chrismarget/terraform-provider-examples/terraform-provider-omdb/omdb"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"log"
)

var commit, version string // populated by goreleaser

// NewOmdbProvider instantiates the provider in main
func NewOmdbProvider() provider.Provider {
	return &omdb.Provider{
		Version: version,
		Commit:  commit,
	}
}

func main() {
	err := providerserver.Serve(context.Background(), NewOmdbProvider, providerserver.ServeOpts{
		Address: "github.com/chrismarget/omdb",
	})
	if err != nil {
		log.Fatal(err)
	}
}
