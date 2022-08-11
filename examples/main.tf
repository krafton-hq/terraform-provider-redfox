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

resource "redfox_cluster" "aa" {
  metadata {
    name      = "my-second-cluster"
    namespace = "redfox-metadata"
    labels = {
      "foo" = "bar"
    }
  }
  spec {
    cluster_name     = "dev-meta-aws"
    cluster_group    = "dev-meta"
    cluster_engine   = "EKS"
    cluster_region   = "ap-northeast-2"
    infra_account_id = "1234567890"
    infra_vendor     = "AWS"
    service_phase    = "dev"
    service_tag      = "meta"
  }
}

resource "redfox_cluster_status" "aa" {
  metadata {
    name      = "my-second-cluster"
    namespace = "redfox-metadata"
  }
  status {
    apiserver {
      endpoint = "https://example.com"
      ca_cert = base64encode("MYCERT")
    }
    service_account_issuer = "https://example.com"
    aws_iam_idps = {
      "foo" = "bar"
      "fizz" = "buzz"
    }
  }

  depends_on = [redfox_cluster.aa]
}

data "redfox_clusters" "a" {
  namespace = "redfox-metadata"
  selector {
    match_labels = {
      foo = "bar"
    }
  }
}
