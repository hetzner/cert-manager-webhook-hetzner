# Changelog

## [v0.1.0](https://github.com/hetzner/cert-manager-webhook-hetzner/releases/tag/v0.1.0)

This release introduces the new [cert-manager](https://cert-manager.io) webhook for Hetzner.

The webhook relies on the new [DNS API](https://docs.hetzner.cloud/reference/cloud#dns).

The DNS API is currently in **beta**, which will likely end on 10 November 2025. See the
[DNS Beta FAQ](https://docs.hetzner.com/networking/dns/faq/beta) for more details.

The webhook is currently experimental, breaking changes may occur within minor releases.

To get started, head to the [webhook documentation](https://github.com/hetzner/cert-manager-webhook-hetzner/blob/main/docs).

### Features

- new cert-manager webhook for Hetzner
