# Issuer configuration reference

This page references the `Issuer` and `ClusterIssuer` configurations for the cert-manager webhook Hetzner.

The webhook is responsible for one or more issuers, each with its own configuration. It is also possible for each Issuer to use its own Hetzner Cloud API token. The following options are available:

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
        <td>string</td>
        <td></td>
        <td>
            Name of the Kubernetes secret, which stores the Hetzner Cloud API token.
            Required unless <code>config.tokenFilePath</code> is set.
        </td>
    </tr>
    <tr>
        <td><code>config.tokenSecretKeyRef.key</code></td>
        <td>string</td>
        <td></td>
        <td>
            Key in the Kubernetes secret, which stores the Hetzner Cloud API token.
            Required unless <code>config.tokenFilePath</code> is set.
        </td>
    </tr>
    <tr>
        <td><code>config.tokenFilePath</code></td>
        <td>string</td>
        <td></td>
        <td>
            Path to a file containing the Hetzner Cloud API token, mounted into
            the webhook pod. Mutually exclusive with <code>config.tokenSecretKeyRef</code>;
            setting both is an error. Leading and trailing whitespace in the file
            are trimmed.
        </td>
    </tr>
</table>

### Example: token from Kubernetes secret

```yaml
# issuer.yaml
# [...]
solvers:
  - dns01:
      webhook:
        groupName: acme.hetzner.com
        solverName: hetzner
        config:
          tokenSecretKeyRef:
            name: hetzner
            key: token
```

### Example: token from mounted file

Mount the token into the webhook pod (for example via the chart's
`extraVolumes` / `extraVolumeMounts` values, a projected volume, or a CSI
secret driver) and reference its path in the issuer config:

```yaml
# issuer.yaml
# [...]
solvers:
  - dns01:
      webhook:
        groupName: acme.hetzner.com
        solverName: hetzner
        config:
          tokenFilePath: /var/run/secrets/hetzner/token
```
