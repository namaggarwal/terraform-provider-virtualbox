package virtualbox

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/namaggarwal/go-virtualbox"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"virtualbox_vm": resourceVM(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	return virtualbox.NewVBoxManage(), nil
}
