SHELL = bash
KUBERNETES_VERSION ?= 1.34.1 # renovate: datasource=github-releases depName=kubernetes-sigs/controller-tools extractVersion=^envtest-v(?<version>.+)$

# Override PATH so that the cert-manager test environment detects the required binaries.
testdata/env:
	setup-envtest use -p env $(KUBERNETES_VERSION) > $@
	echo 'export PATH=$$KUBEBUILDER_ASSETS:$$PATH' >> $@

testdata/hetzner/secret.yaml: testdata/hetzner/secret.yaml.tpl
	@sed \
	-e 's/$$HETZNER_TOKEN_BASE64/$(shell echo -n "$(HETZNER_TOKEN)" | base64 -w0)/' \
	$< > $@

TESTDATA = testdata/env testdata/hetzner/secret.yaml

.PHONY: test
test:
	go test -v ./...

.PHONY: e2e-setup
e2e-setup: $(TESTDATA)

.PHONY: e2e
e2e: $(TESTDATA)
	source testdata/env; go test -v -tags e2e -timeout=30m ./...

.PHONY: e2e-coverage
e2e-coverage: $(TESTDATA)
	source testdata/env; go test -v -tags e2e -timeout=30m -coverpkg=./... -coverprofile=coverage.txt -covermode count ./...

.PHONY: helm-package
helm-package:
ifndef VERSION
	$(error "VERSION is not defined")
endif
	sed -e "s/version: .*/version: $(VERSION)/" --in-place chart/Chart.yaml
	helm package chart

.PHONY: clean
clean:
	rm -Rf $(TESTDATA)
