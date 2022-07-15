package provider

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var providerFactories = map[string]func() (*schema.Provider, error){
	"redfox": func() (*schema.Provider, error) {
		return New("dev-build")(), nil
	},
}
