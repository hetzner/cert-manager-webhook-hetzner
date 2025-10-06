variable "name" {
  type = string
}

variable "hetzner_token" {
  type      = string
  sensitive = true
}

variable "k3s_channel" {
  type    = string
  default = "stable"
}
