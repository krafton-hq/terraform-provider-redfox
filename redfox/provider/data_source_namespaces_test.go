package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestDataSourceNamespaces(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigDataSourcesNamespace(""),
				Check: func(state *terraform.State) error {
					return nil
				},
			},
		},
	})
}

func testConfigDataSourcesNamespace(name string) string {
	return fmt.Sprintf(`
data "redfox_namespaces" "varname" {
}
`)
}
