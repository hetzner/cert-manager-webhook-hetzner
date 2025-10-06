# cert-manager-webhook-hetzner

![Maturity](https://img.shields.io/badge/maturity-experiment-orange)

This webhook creates the necessary DNS entries in the [Hetzner DNS API](https://docs.hetzner.cloud/reference/cloud#zones) to solve a [DNS01 challenge](https://letsencrypt.org/docs/challenge-types/#dns-01-challenge) for a cert-manager [`Issuer`](https://cert-manager.io/docs/concepts/issuer/) of the [ACME](https://cert-manager.io/docs/configuration/acme/) type.

## Docs

- :rocket: See the [quick start guide](docs/guides/quickstart.md) to get you started.
- :book: See the [configuration reference](docs/reference/issuer-configuration.md) for the available configuration.

For more information, see the [documentation](docs/).

## Development

### Start a development environment

1. Configure a `HETZNER_TOKEN` in your shell session.
2. Deploy the development cluster.

```bash
make -C dev up
```

3. Load the generated configuration to access the development cluster:

```bash
source dev/files/env.sh
```

4. Start developing cert-manager-webhook-hetzner in the development cluster:

```bash
skaffold dev
```

5. Test your deployment by placing your zone name into `commonName` and `dnsName` of `dev/example-cert.yaml`:

```bash
kubectl apply -f dev/example-cert.yaml
```

6. Wait for your certificate to be issued. This can take up to two minutes:

```bash
kubectl -n cert-manager get certificates example-cert -w
```

⚠️ Do not forget to clean up the development cluster once are finished:

```bash
make -C dev down
```

### Run the unit tests

```bash
go test ./internal/... -v
```

### Run the cert-manager conformance test suite

All DNS providers **must** run the DNS01 provider conformance testing suite,
else they will have undetermined behaviour when used with cert-manager.

You can run the test suite by:

1. Placing your base64 encoded hcloud API token in `testdata/hetzner/secret.yaml`
2. Run the test suite:

```bash
make test
```
