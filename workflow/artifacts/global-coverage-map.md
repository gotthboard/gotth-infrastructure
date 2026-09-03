# Coverage map

| Surface | Evidence |
|---|---|
| Service and real Compose-secret contract | `spec_test.go` |
| Deterministic least-privilege Compose | `compose_test.go`, `make verify-compose` |
| Exact desired/effective container inspection | `inspect_test.go`, disposable runtime evidence |
| Operator command boundary | `cmd/gotth-infrastructure/main_test.go` |
| External consumer API | `pkg/ops/public_api_test.go` |
| Alpha.2 provenance | `examples/gotth-bb/` |

Implementation and outside-package tests above now live under `pkg/ops/`;
`pkg/ops/public_api_test.go` and the CLI import the canonical package.
