locals {
  kubeconfig_path     = abspath("${path.root}/files/kubeconfig.yaml")
  pebble_config_path  = abspath("${path.module}/pebble-config.json")
  unbound_config_path = abspath("${path.module}/unbound.conf")

  pebble_args_base = ["-config=/pebble-config.json"]
  pebble_args      = var.use_unbound ? concat(local.pebble_args_base, ["-dnsserver=unbound.unbound.svc.cluster.local:53"]) : local.pebble_args_base
}

data "local_sensitive_file" "kubeconfig" {
  depends_on = [module.dev]
  filename   = local.kubeconfig_path
}

data "local_file" "pebble_config" {
  filename = local.pebble_config_path
}

provider "helm" {
  kubernetes = {
    config_path = data.local_sensitive_file.kubeconfig.filename
  }
}

provider "kubernetes" {
  config_path = data.local_sensitive_file.kubeconfig.filename
}

module "dev" {
  source = "github.com/hetznercloud/kubernetes-dev-env?ref=v0.10.0"

  name         = "cert-manager-webhook-${replace(var.name, "/[^a-zA-Z0-9-_]/", "-")}"
  hcloud_token = var.hetzner_token
  worker_count = 0

  k3s_channel = var.k3s_channel
}

resource "kubernetes_namespace_v1" "cert-manager" {
  depends_on = [module.dev]
  metadata {
    name = "cert-manager"
  }
}

resource "helm_release" "cert_manager" {
  depends_on = [kubernetes_namespace_v1.cert-manager]
  name       = "cert-manager"
  chart      = "cert-manager"
  repository = "https://charts.jetstack.io"
  version    = "v1.19.3"
  namespace  = "cert-manager"

  set = [
    {
      name  = "crds.enabled"
      value = "true"
    },
    {
      name  = "extraArgs"
      value = var.use_unbound ? "{--dns01-recursive-nameservers-only,--dns01-recursive-nameservers=unbound.unbound.svc.cluster.local:53}" : "{}"
    }
  ]

  provisioner "local-exec" {
    when    = destroy
    command = ". files/env.sh && kubectl delete apiservices.apiregistration.k8s.io v1alpha1.acme.hetzner.com || true"
  }
}

resource "kubernetes_secret_v1" "hetzner_token" {
  depends_on = [kubernetes_namespace_v1.cert-manager]
  metadata {
    name      = "hetzner"
    namespace = "cert-manager"
  }

  data = {
    token = var.hetzner_token
  }
}

resource "kubernetes_config_map_v1" "pebble" {
  depends_on = [kubernetes_namespace_v1.cert-manager]
  metadata {
    name      = "pebble"
    namespace = "cert-manager"
  }

  data = {
    "pebble-config.json" = data.local_file.pebble_config.content
  }
}

resource "kubernetes_service_v1" "pebble" {
  depends_on = [kubernetes_namespace_v1.cert-manager]
  metadata {
    name      = "pebble"
    namespace = "cert-manager"
  }

  spec {
    port {
      name        = "http"
      protocol    = "TCP"
      port        = 14000
      target_port = "http"
    }

    selector = {
      "app.kubernetes.io/name" = "pebble"
    }

    type = "ClusterIP"
  }
}

resource "kubernetes_deployment_v1" "pebble" {
  depends_on = [kubernetes_config_map_v1.pebble]
  metadata {
    name      = "pebble"
    namespace = "cert-manager"
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        "app.kubernetes.io/name" = "pebble"
      }
    }

    template {
      metadata {
        labels = {
          "app.kubernetes.io/name" = "pebble"
        }
      }

      spec {
        volume {
          name = "config-volume"

          config_map {
            name = "pebble"

            items {
              key  = "pebble-config.json"
              path = "pebble-config.json"
            }
          }
        }

        container {
          name  = "pebble"
          image = "ghcr.io/letsencrypt/pebble:2.9.0" # renovate: datasource=docker depName=ghcr.io/letsencrypt/pebble
          args  = local.pebble_args

          port {
            name           = "http"
            container_port = 14000
            protocol       = "TCP"
          }

          volume_mount {
            name       = "config-volume"
            read_only  = true
            mount_path = "/pebble-config.json"
            sub_path   = "pebble-config.json"
          }

          image_pull_policy = "Always"
        }
      }
    }
  }
}

resource "terraform_data" "pebble-issuer" {
  depends_on = [helm_release.cert_manager, kubernetes_deployment_v1.pebble]
  provisioner "local-exec" {
    when    = create
    command = ". ${path.module}/files/env.sh && kubectl apply -f ${path.module}/pebble-issuer.yaml"
  }
}

resource "kubernetes_namespace_v1" "unbound" {
  count      = var.use_unbound ? 1 : 0
  depends_on = [module.dev]
  metadata {
    name = "unbound"
  }
}

resource "kubernetes_config_map_v1" "unbound" {
  count      = var.use_unbound ? 1 : 0
  depends_on = [kubernetes_namespace_v1.unbound]
  metadata {
    name      = "unbound"
    namespace = "unbound"
  }

  data = {
    "unbound.conf" = fileexists(local.unbound_config_path) ? file(local.unbound_config_path) : ""
  }
}

resource "kubernetes_service_v1" "unbound" {
  count      = var.use_unbound ? 1 : 0
  depends_on = [kubernetes_namespace_v1.unbound]
  metadata {
    name      = "unbound"
    namespace = "unbound"
  }

  spec {
    port {
      name        = "dns"
      protocol    = "UDP"
      port        = 53
      target_port = "dns"
    }

    selector = {
      "app.kubernetes.io/name" = "unbound"
    }

    type = "ClusterIP"
  }
}

resource "kubernetes_deployment_v1" "unbound" {
  count      = var.use_unbound ? 1 : 0
  depends_on = [kubernetes_config_map_v1.unbound]
  metadata {
    name      = "unbound"
    namespace = "unbound"
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        "app.kubernetes.io/name" = "unbound"
      }
    }

    template {
      metadata {
        labels = {
          "app.kubernetes.io/name" = "unbound"
        }
      }

      spec {
        volume {
          name = "config-volume"

          config_map {
            name = "unbound"

            items {
              key  = "unbound.conf"
              path = "unbound.conf"
            }
          }
        }

        container {
          name  = "unbound"
          image = "ghcr.io/crazy-max/unbound:1.24.0"

          port {
            name           = "dns"
            container_port = 5053
            protocol       = "UDP"
          }

          volume_mount {
            name       = "config-volume"
            read_only  = true
            mount_path = "/config/unbound.conf"
            sub_path   = "unbound.conf"
          }

          image_pull_policy = "Always"
        }
      }
    }
  }
}
