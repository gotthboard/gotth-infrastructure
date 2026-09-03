# Coverage

The canonical behavior map is `workflow/artifacts/global-coverage-map.md`.
`make verify` covers `pkg/ops`, the CLI, and Compose normalization. The public
library is at 100.0% statement coverage and the CLI at 93.2%; the remaining CLI
branch is the process-level `main`/`os.Exit` path while the executable boundary
is exercised by integration.
