// provider.go
package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"example_server": resourceServerA(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"example_server": dataSourceCoffees(),
		},
	}
}
