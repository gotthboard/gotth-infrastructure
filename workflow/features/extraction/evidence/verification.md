# Extraction verification

Verified on 2026-09-02 with Go `go1.26.6-X:nodwarf5`:

- `go vet -mod=readonly ./...`
- `go test -mod=readonly -race -cover ./...`
- library statement coverage: 98.4%
- command statement coverage: 75.0%
- deterministic GOTTH Board example rendering succeeds
- effective runtime checks accept unrelated inspect fields but reject every
  missing containment control
- no live container, database, proxy, secret, volume, backup, or host touched

Command-process exit behavior is not invoked from unit tests; the testable run
boundary covers render, inspect verification, usage, and read failures.
