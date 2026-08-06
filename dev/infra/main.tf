module "infra" {
  source = "github.com/hetznercloud/kubernetes-dev-env//modules/infra?ref=v0.11.0"

  name         = "cert-manager-webhook-${replace(var.name, "/[^a-zA-Z0-9-_]/", "-")}"
  hcloud_token = var.hetzner_token
  worker_count = 0

  k3s_channel = var.k3s_channel

  # Share the generated files with the k8s state
  output_dir = abspath("${path.root}/../files")
}
