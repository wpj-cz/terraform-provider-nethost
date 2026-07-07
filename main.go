package main

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"gitlab.wpj.cz/terraform-provider-nethost/internal/provider"
)

func main() {
	err := providerserver.Serve(context.Background(), provider.NewProvider, providerserver.ServeOpts{
		Address: "registry.terraform.io/wpj-cz/nethost",
	})

	if err != nil {
		log.Fatal(err)
	}
}
