# Quick start

Before deploying the webhook you need to install cert-manager. You can find the official install instructions [here](https://cert-manager.io/docs/installation/).

1. Create a read+write API token in the [Hetzner Cloud Console](https://console.hetzner.cloud/) as described in the [official guide](https://docs.hetzner.com/cloud/api/getting-started/generating-api-token).
2. Create a secret containing the API token. If you plan to create a namespace scoped `Issuer`, place the secret in the Issuers namespace. If you plan to configure a `ClusterIssuer`, place the secret in the configured [cluster resource namespace](https://cert-manager.io/docs/configuration/#cluster-resource-namespace), defaulting to `cert-manager`.

```yaml
# secret.yml
apiVersion: v1
kind: Secret
metadata:
  name: hetzner
  namespace: cert-manager
stringData:
  token: <HETZNER_TOKEN>
```

Apply the secret:

```bash
kubectl apply -f secret.yml
```

3. Settle on a group name. The group name is used to uniquely identify your company or business unit that created this webhook. It is referenced in each Issuer's `webhook` stanza to inform cert-manager of where to send ChallengePayload resources in order to solve the DNS01 challenge. A suitable choice is you own company's domain. For the quick start guide we settle for `acme.mycompany.com`.

4. Add the Helm repository:

```bash
helm repo add hcloud https://charts.hetzner.cloud
helm repo update hcloud
```

5. Install the webhook:

```bash
helm install hetzner-cert-manager-webhook hcloud/hetzner-cert-manager-webhook -n cert-manager --set groupName="acme.mycompany.com"
```

6. Configure an `Issuer` or `ClusterIssuer` to suit your needs, as described in the official cert-manager [documentation](https://cert-manager.io/docs/configuration/acme/).

7. Configure the Hetzner webhook as a DNS01 solver at your `Issuer` or `ClusterIssuer` ([reference](https://cert-manager.io/docs/configuration/acme/dns01/webhook/)):

```yaml
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: cluster-issuer
spec:
  acme:
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

8. You can now start [requesting certificates](https://cert-manager.io/docs/usage/).
