# cert-manager-webhook-hetzner

cert-manager ACME webhook for Hetzner

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| affinity | object | `{}` | [Kubernetes affinities](https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/#affinity-and-anti-affinity) for the webhook. |
| certManager.namespace | string | `"cert-manager"` | Namespace of your cert-manager deployment. |
| certManager.serviceAccountName | string | `"cert-manager"` | Name of the cert-managers service account. |
| env | object | `{}` | Additional environment variables, where each key represents the name of the variable. The value follows standard Kubernetes environment variable formats. |
| fullnameOverride | string | `""` | Override the full name of the chart. |
| groupName | string | `"my-hetzner-project"` | The GroupName here is used to identify your company or business unit that created this webhook. For example, this may be "acme.mycompany.com". This name will need to be referenced in each Issuer's `webhook` stanza to inform cert-manager of where to send ChallengePayload resources in order to solve the DNS01 challenge. This group name should be **unique**, hence using your own company's domain here is recommended. Each webhook deployment is responsible for one Hetzner Cloud project, where the access token is provided. A Hetzner Cloud project can contain multiple zones. |
| image.pullPolicy | string | `"IfNotPresent"` | Pull policy of the webhooks image. |
| image.repository | string | `"docker.io/hetznercloud/cert-manager-webhook-hetzner"` | Repository of the webhooks image. |
| image.tag | string | Current version | Tag of the webhooks image. |
| imagePullSecrets | list | `[]` | Additional image pull secrets in the [standard Kubernetes format](https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/) |
| metrics.serviceMonitor.enabled | bool | `false` | Deploys a ServiceMonitor to scrape the metrics. **Requires** the ServiceMonitor CRD. |
| nameOverride | string | `""` | Override the name of the chart. |
| nodeSelector | object | `{}` | [Kubernetes node selector](https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/#nodeselector) for the webhook. |
| replicaCount | int | `1` | Number of replicas. |
| resources | object | `{}` | [Kubernetes resource management](https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/) for the webhook |
| service.port | int | `443` | Port of the webhooks service. |
| service.type | string | `"ClusterIP"` | Kubernetes service type of the webhooks service. |
| tolerations | list | `[]` | [Kubernetes tolerations](https://kubernetes.io/docs/concepts/scheduling-eviction/taint-and-toleration) for the webhook. |
