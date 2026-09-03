# Coverage map

| Surface | Evidence |
|---|---|
| Service and real Compose-secret contract | `spec_test.go` |
| Deterministic least-privilege Compose | `compose_test.go`, `make verify-compose` |
| Exact desired/effective container inspection | `inspect_test.go`, disposable runtime evidence |
| Operator command boundary | `cmd/gotth-infrastructure/main_test.go` |
| External consumer API | `public_api_test.go` |
| Alpha.2 provenance | `examples/gotth-bb/` |
