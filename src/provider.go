package main

import (
	"github.com/hashicorp/terraform/helper/schema"
)

//Provider for dkron
func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"dkronjob": DkronJob(),
		},
	}
}
