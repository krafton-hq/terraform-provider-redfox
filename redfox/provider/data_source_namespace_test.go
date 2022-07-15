package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestDataSourceNamespace(t *testing.T) {
	randName := "test2"

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigDataSourceNamespace(randName),
			},
		},
	})
}

func testConfigDataSourceNamespace(name string) string {
	return fmt.Sprintf(`
data "redfox_namespace" "varname" {
  metadata {
    name = "%s"
  }
}
`, name)
}
