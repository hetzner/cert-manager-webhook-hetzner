all: lint snapshots-diff

VALUES = $(wildcard *-values.yaml)
SNAPSHOTS = $(VALUES:%-values.yaml=.snapshots/%.yaml)

snapshots-clean:
	rm -f .snapshots/*

.snapshots/%.yaml: %-values.yaml
	helm template snapshot -n integrations -f "$<" . > "$@"

snapshots: snapshots-clean $(SNAPSHOTS)

snapshots-diff: snapshots
	git diff --exit-code -- .snapshots/

lint:
	helm lint .
