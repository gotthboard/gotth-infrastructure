.PHONY: verify

verify:
	version="$$(go env GOVERSION)"; test "$${version%%-*}" = "go1.26.6"
	test -z "$$(gofmt -l $$(git ls-files --cached --others --exclude-standard -- '*.go'))"
	go vet -mod=readonly ./...
	go test -mod=readonly -race -cover ./...
	go run ./cmd/gotth-ops render examples/gotth-bb/service.json >/tmp/gotth-bb-generated-compose.yml
