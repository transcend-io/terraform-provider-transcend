package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/transcend-io/terraform-provider-transcend/transcend"
)

// Generate the Terraform provider documentation using `tfplugindocs`:
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: transcend.Provider,
	})
}
