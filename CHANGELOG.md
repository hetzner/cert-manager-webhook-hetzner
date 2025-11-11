# Changelog

## [v0.6.0](https://github.com/hetzner/cert-manager-webhook-hetzner/releases/tag/v0.6.0)

### Features

- **helm**: default to restricted security context (#59)

### Bug Fixes

- **ko**: use non-root user in container image (#57)

## [v0.5.0](https://github.com/hetzner/cert-manager-webhook-hetzner/releases/tag/v0.5.0)

### Features

- **chart**: add possibility to configure deployments .spec.strategy (#52)
- **chart**: configure deployments metadata annotations and labels (#51)
- **chart**: add pod disruption budget (#53)

## [v0.4.0](https://github.com/hetzner/cert-manager-webhook-hetzner/releases/tag/v0.4.0)

### New container image namespace

With this release, we moved our container image to [docker.io/hetzner/cert-manager-webhook-hetzner](https://hub.docker.com/r/hetzner/cert-manager-webhook-hetzner).

Users of the Helm chart need to take action if they have manually set the `image.repository` value.

```diff
 # ...
 image:
-  repository: docker.io/hetznercloud/cert-manager-webhook-hetzner
+  repository: docker.io/hetzner/cert-manager-webhook-hetzner
 # ...
```

Existing images in the old `hetznercloud` namespace will remain available, but new images will only be pushed to the new `hetzner` namespace.

### Features

- push image to hetzner docker namespace (#45)

## [v0.3.0](https://github.com/hetzner/cert-manager-webhook-hetzner/releases/tag/v0.3.0)

### Features

- **helm**: configurable podSecurityContext and securityContext (#37)

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
