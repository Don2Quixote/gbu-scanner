.PHONY: run test

.SILENT:

run:
	./scripts/run.sh

test:
	./scripts/test.sh

lint:
	./scripts/lint.sh