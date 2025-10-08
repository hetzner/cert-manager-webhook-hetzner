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
            Always set this value to <code>acme.hetzner.com</code>,
            unless you configured a different <code>groupName</code> for the Helm chart.
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
