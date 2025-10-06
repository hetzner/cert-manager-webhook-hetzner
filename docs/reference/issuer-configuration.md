# Issuer configuration reference

This page references the `Issuer` and `ClusterIssuer` configurations for the cert-manager webhook Hetzner.

The webhook is responsible for one or more issuers, each with its own configuration. The following options are available:

<table>
    <tr>
        <th>Option</th>
        <th>Type</th>
        <th>Default</th>
        <th>Description</th>
    </tr>
    <tr>
        <td><code>groupName</code></td>
        <td>string (<strong>Required</strong>)</td>
        <td></td>
        <td>
            The group name is used to uniquely identify your company or business unit that created this webhook.
            It is referenced in each Issuer's `webhook` stanza to inform cert-manager of where to send `ChallengePayload` resources
            in order to solve the DNS01 challenge. A suitable choice is your own company's domain.
        </td>
    </tr>
    <tr>
        <td><code>solverName</code></td>
        <td>string (<strong>Required</strong>)</td>
        <td></td>
        <td>Always set this value to <code>hetzner</code></td>
    </tr>
    <tr>
        <td><code>config.tokenSecretKeyRef.name</code></td>
        <td>string (<strong>Required</strong>)</td>
        <td></td>
        <td>Name of the Kubernetes secret, which stores the Hetzner Cloud API token.</td>
    </tr>
    <tr>
        <td><code>config.tokenSecretKeyRef.key</code></td>
        <td>string (<strong>Required</strong>)</td>
        <td></td>
        <td>Key in the Kubernetes secret, which stores the Hetzner Cloud API token.</td>
    </tr>
</table>

### Example

```yaml
# issuer.yaml
# [...]
solvers:
  - dns01:
      webhook:
        groupName: acme.mycompany.com
        solverName: hetzner
        config:
          tokenSecretKeyRef:
            name: hetzner
            key: token
```
