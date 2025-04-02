package main

import (
	"context"
	"flag"
	"log"

	"github.com/bdwyertech/terraform-provider-json2dynamodb/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// Run "go generate" to format example terraform files and generate the docs for the registry/website

// If you do not have terraform installed, you can remove the formatting command, but its suggested to
// ensure the documentation is formatted properly.
//go:generate terraform fmt -recursive ./examples/

// Run the docs generation tool, check its repository for more information on how it works and how docs
// can be customized.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate

var (
	version string = "dev"
	// commit  string = ""
)

func main() {
	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/bdwyertech/json2dynamodb",
	}

	flag.BoolVar(&opts.Debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	if err := providerserver.Serve(context.Background(), provider.New(version), opts); err != nil {
		log.Fatal(err.Error())
	}
}
