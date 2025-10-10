# Changelog

## [v0.2.0](https://github.com/hetzner/cert-manager-webhook-hetzner/releases/tag/v0.2.0)

### Default GroupName `acme.hetzner.com`

With this version we have changed the default `groupName` of the Helm chart to `acme.hetzner.com`. If you have prevously deployed the webhook according to our quickstart guide, you will not run into any issues when upgrading. If you want to update your existing webhook to the `acme.hetzner.com` `groupName` you have to specify the Helm value `groupName=acme.hetzner.com` during the upgrade and update your existing Issuers/ClusterIssuers accordingly.

### Features

- use default groupName acme.hetzner.com (#28)

## [v0.1.1](https://github.com/hetzner/cert-manager-webhook-hetzner/releases/tag/v0.1.1)

### Bug Fixes

- do not use the latest container image tag (#18)

## [v0.1.0](https://github.com/hetzner/cert-manager-webhook-hetzner/releases/tag/v0.1.0)

This release introduces the new [cert-manager](https://cert-manager.io) webhook for Hetzner.

The webhook relies on the new [DNS API](https://docs.hetzner.cloud/reference/cloud#dns).

The DNS API is currently in **beta**, which will likely end on 10 November 2025. See the
[DNS Beta FAQ](https://docs.hetzner.com/networking/dns/faq/beta) for more details.

The webhook is currently experimental, breaking changes may occur within minor releases.

To get started, head to the [webhook documentation](https://github.com/hetzner/cert-manager-webhook-hetzner/blob/main/docs).

### Features

- new cert-manager webhook for Hetzner
