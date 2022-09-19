package main

import (
	"context"
	"github.com/chrismarget/terraform-provider-examples/terraform-provider-omdb/omdb"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"log"
)

func main() {
	err := providerserver.Serve(context.Background(), omdb.New, providerserver.ServeOpts{
		Address: "github.com/chrismarget/omdb",
	})
	if err != nil {
		log.Fatal(err)
	}
}
