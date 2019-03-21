package main

import (
  "github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
	
	"github.com/peertransfer/terraform-provider-dkron/dkron"
)

func main() {
  plugin.Serve(&plugin.ServeOpts{
    ProviderFunc: func() terraform.ResourceProvider {
      return dkron.Provider()
    },
  })
}
