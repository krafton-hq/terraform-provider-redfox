package main

import (
	"flag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/krafton-hq/terraform-provider-redfox/redfox"
)

var version = "dev-build"
var commit = "devcommit"

//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs
func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := &plugin.ServeOpts{
		Debug: debugMode,

		ProviderAddr: "registry.terraform.io/krafton-hq/redfox",

		ProviderFunc: redfox.Provider,
	}

	plugin.Serve(opts)
}
