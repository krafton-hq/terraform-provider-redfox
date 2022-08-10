terraform {
  required_providers {
    redfox = {
      source  = "krafton-hq/redfox"
      version = "0.0.2"
    }
  }
}

provider "redfox" {
  config_path = "~/.kube/config"
}

resource "redfox_natip" "test" {
  metadata {
    name      = "my-first-nat"
    namespace = "redfox-metadata"
    labels = {
      "foo" = "bar"
    }
  }
  spec {
    ip_type = "Ipv4"
    cidrs = ["1.1.1.1/32"]
  }
}

data "redfox_natips" "test2" {
  namespace = "redfox-metadata"
  selector {

  }
}

output "a" {
  value = data.redfox_natips.test2
}
