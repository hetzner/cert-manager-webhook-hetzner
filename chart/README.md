# cert-manager-webhook-hetzner

cert-manager ACME webhook for Hetzner

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| affinity | object | `{}` | [Kubernetes affinities](https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/#affinity-and-anti-affinity) for the webhook. |
| certManager.namespace | string | `"cert-manager"` | Namespace of your cert-manager deployment. |
| certManager.serviceAccountName | string | `"cert-manager"` | Name of the cert-managers service account. |
| containerSecurityContext | object | `{}` | [Kubernetes container security context](https://kubernetes.io/docs/tasks/configure-pod-container/security-context/#set-the-security-context-for-a-container) for the webhook. |
| env | object | `{}` | Additional environment variables, where each key represents the name of the variable. The value follows standard Kubernetes environment variable formats. |
| fullnameOverride | string | `""` | Override the full name of the chart. |
| groupName | string | `"acme.hetzner.com"` | The GroupName here is used to identify your company or business unit that created this webhook. For example, this may be "acme.mycompany.com". This name will need to be referenced in each Issuer's `webhook` stanza to inform cert-manager of where to send ChallengePayload resources in order to solve the DNS01 challenge. This group name should be **unique**, hence using your own company's domain here is recommended. |
| image.pullPolicy | string | `"IfNotPresent"` | Pull policy of the webhook image. |
| image.repository | string | `"docker.io/hetznercloud/cert-manager-webhook-hetzner"` | Repository of the webhook image. |
| image.tag | string | Current version | Tag of the webhook image. |
| imagePullSecrets | list | `[]` | Additional image pull secrets in the [standard Kubernetes format](https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/) |
| metrics.serviceMonitor.enabled | bool | `false` | Deploys a ServiceMonitor to scrape the metrics. **Requires** the ServiceMonitor CRD. |
| nameOverride | string | `""` | Override the name of the chart. |
| nodeSelector | object | `{}` | [Kubernetes node selector](https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/#nodeselector) for the webhook. |
| podSecurityContext | object | `{}` | [Kubernetes pod security context](https://kubernetes.io/docs/tasks/configure-pod-container/security-context/#set-the-security-context-for-a-pod) for the webhook. |
| replicaCount | int | `1` | Number of replicas. |
| resources | object | `{}` | [Kubernetes resource management](https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/) for the webhook |
| service.port | int | `443` | Port of the webhook service. |
| service.type | string | `"ClusterIP"` | Kubernetes service type of the webhook service. |
| tolerations | list | `[]` | [Kubernetes tolerations](https://kubernetes.io/docs/concepts/scheduling-eviction/taint-and-toleration) for the webhook. |
