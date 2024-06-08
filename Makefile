SHELL=bash

include .env

start:
	docker compose up -d

dev:
	$(eval export $(sed 's/#.*//g' .env | xargs))
	(cd cmd/web && go run main.go)

dev-sentry:
	$(eval export $(sed 's/#.*//g' .env | xargs))
	(cd cmd/sentry && go run main.go)

test:
	$(eval export $(sed 's/#.*//g' .env | xargs))
	go test -json -skip /pkg/test -v ./... $(args) 2>&1 | gotestfmt

test-match:
	make test args="-run $(case)"

test-ci:
	set -euo pipefail
	go test -json -skip /pkg/test -v ./... 2>&1 | gotestfmt

vulns-check:
	govulncheck ./...

docs:
	pkgsite -http=:4060