package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/transcend-io/terraform-provider-transcend-io/transcend"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: transcend.Provider,
	})
}
