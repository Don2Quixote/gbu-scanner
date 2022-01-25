.PHONY: run test lint stat

.SILENT:

run:
	./scripts/run.sh

test:
	./scripts/test.sh

lint:
	./scripts/lint.sh

stat:
	./scripts/stat.sh