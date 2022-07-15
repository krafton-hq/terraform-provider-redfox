package provider

import (
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestResourceNatIp(t *testing.T) {
	randName := gofakeit.UUID()
	randNs := gofakeit.UUID()

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigResourceNatIp(randNs, randName),
			},
		},
	})
}

func testConfigResourceNatIp(namespace string, name string) string {
	return fmt.Sprintf(`
resource "redfox_namespace" "nstest" {
  metadata {
    name      = "%s"
  }
  spec {
	api_objects {
		kind = "NatIp"
		group = "red-fox.sbx-central.io"
		version = "v1alpha1"
	}
  }
}

resource "redfox_natip" "varname" {
  metadata {
    name      = "%s"
	namespace = "%s"
    labels    = {
      "key" = "value224"
    }
  }
  spec {
	ip_type = "Ipv4"
	cidrs = ["1.1.1.1/32", "2.2.2.2/32", "1.3.3.2/24"]
  }

  depends_on = [redfox_namespace.nstest]
}
`, namespace, name, namespace)
}
