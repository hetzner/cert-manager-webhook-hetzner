variable "hetzner_token" {
  type      = string
  sensitive = true
}

variable "use_unbound" {
  type    = bool
  default = false
}
