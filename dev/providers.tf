terraform {
  required_providers {
    helm = {
      source  = "hashicorp/helm"
      version = "< 4.0.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "2.38.0"
    }
  }
}
