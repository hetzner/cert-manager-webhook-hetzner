terraform {
  required_providers {
    local = {
      source  = "hashicorp/local"
      version = ">= 2.5.1, < 3.0.0"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "< 4.0.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "< 4.0.0"
    }
  }
}

# Both providers are configured from the infra state, and never from an
# attribute of a resource managed in this state. This is what allows OpenTofu to
# plan (and destroy) the resources below without having to create the cluster
# first.

provider "kubernetes" {
  config_path = local.kubeconfig_path
}

provider "helm" {
  kubernetes = {
    config_path = local.kubeconfig_path
  }
}
