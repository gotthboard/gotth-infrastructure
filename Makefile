.PHONY: verify verify-compose

verify:
	version="$$(go env GOVERSION)"; test "$${version%%-*}" = "go1.26.6"
	test -z "$$(gofmt -l $$(git ls-files --cached --others --exclude-standard -- '*.go'))"
	go vet -mod=readonly ./...
	go test -mod=readonly -race -cover ./...
	$(MAKE) verify-compose

verify-compose:
	command -v docker >/dev/null
	docker compose version >/dev/null
	compose_file="$$(mktemp /tmp/gotth-infrastructure-compose.XXXXXX)"; normalized_file="$$(mktemp /tmp/gotth-infrastructure-normalized.XXXXXX)"; database_secret="$$(mktemp /tmp/gotth-infrastructure-database-secret.XXXXXX)"; oidc_secret="$$(mktemp /tmp/gotth-infrastructure-oidc-secret.XXXXXX)"; trap 'rm -f -- "$$compose_file" "$$normalized_file" "$$database_secret" "$$oidc_secret"' EXIT; go run ./cmd/gotth-infrastructure render examples/gotth-bb/service.json >"$$compose_file"; GOTTH_BB_DATABASE_URL_FILE="$$database_secret" GOTTH_BB_OIDC_CLIENT_SECRET_FILE="$$oidc_secret" docker compose -f "$$compose_file" config --format json >"$$normalized_file"; jq -e '.services.app.tmpfs == ["/tmp:rw,noexec,nosuid,nodev,size=16m"] and (.services.app.secrets | length) == 2 and (.secrets | length) == 2' "$$normalized_file" >/dev/null
