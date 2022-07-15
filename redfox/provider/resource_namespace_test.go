package provider

import (
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestResourceNamespace(t *testing.T) {
	randName := gofakeit.UUID()

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigResourceNamespace(randName),
			},
		},
	})
}

func testConfigResourceNamespace(name string) string {
	return fmt.Sprintf(`
resource "redfox_namespace" "varname" {
  metadata {
    name      = "%s"
    labels    = {
      "key" = "value224"
    }
    annotations = {
      "key" = "value"
    }
  }
  spec {
    api_objects {
      group   = "1"
      kind    = "2"
      version = "3"
    }
    api_objects {
      group   = "4"
      kind    = "5"
      version = "7"
    }
  }
}
`, name)
}
