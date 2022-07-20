terraform {
  required_providers {
    redfox = {
      source  = "krafton-hq/redfox"
      version = "0.0.2"
    }
  }
}

provider "redfox" {
  address = "localhost:8081"
  use_tls = false
}

resource "redfox_crd" "test" {
  metadata {
    name      = "baikal.sbx-central.io"
  }
  spec {
    gvk {
      group   = "sbx-central.io"
      kind    = "Baikal"
      version = "v1alpha1"
    }
  }
}

resource "redfox_namespace" "test" {
  metadata {
    name      = "test-ns"
    labels    = {
      "key" = "value2242323"
    }
    annotations = {
      "key" = "value"
    }
  }
  spec {
    api_objects {
      group   = "sbx-central.io"
      kind    = "Baikal"
      version = "v1alpha1"
    }
    api_objects {
      group   = "core"
      kind    = "NatIp"
      version = "v1"
    }
    api_objects {
      group   = "core"
      kind    = "Endpoint"
      version = "v1"
    }
  }

  depends_on = [redfox_crd.test]
}

resource "redfox_natip" "testtt" {
  metadata {
    name = "nat1"
    namespace = redfox_namespace.test.metadata[0].name
    labels = {
      "key" = "discover"
      "key2" = "non-discover/sss"
    }
  }
  spec {
    ip_type = "Ipv4"
    cidrs = ["1.1.1.1/32"]
  }
}

resource "redfox_natip" "test222" {
  metadata {
    name = "nat2"
    namespace = redfox_namespace.test.metadata[0].name
    labels = {
      "key" = "discover222"
    }
  }
  spec {
    ip_type = "Ipv4"
    cidrs = ["2.3.2.3/32", "182.168.0.0/24"]
  }
}

resource "redfox_customdocument" "tttt" {
  api_version = "sbx-central.io/v1alpha1"
  kind = "Baikal"
  metadata {
    name = "doc1"
    namespace = redfox_namespace.test.metadata[0].name
  }
  spec {
    raw_json = jsonencode({
      "foo" = "bar"
      "kkkk" = "kkkkk"
    })
  }

  depends_on = [redfox_crd.test]
}
