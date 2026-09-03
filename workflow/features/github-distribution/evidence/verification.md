# GitHub Distribution Verification

Status: complete

## Identity and scope

- Pinned base: `9bcfd0f3ff654a3280394e4d659bcfb74a2f68a5`
- Publicly verified candidate: `ea351124da9de13ed9a13efbb50916fb8b4484c6`
- Declared module: `github.com/gotthboard/gotth-infrastructure`
- Runtime/API behavior: unchanged; this is a module-identity and distribution
  contract migration.

Exact stale-prefix searches found no old module or import identity in Go source,
`go.mod`, examples, or fixtures. Canonical Forgejo URLs remain only where the
development, issue, contribution, and security-reporting endpoints are stated.

## Verification

- Local `go mod tidy` produced no dependency drift.
- Local `go vet -mod=readonly ./...` passed.
- Local `go test -mod=readonly ./...` passed.
- On `development`, `make verify` passed with race coverage package 100%; CLI 93.2%.
- On `development`, `go test -mod=readonly -race -count=50 ./...` passed.
- Docker Compose parser verification passed.
- A fresh public GitHub clone of `feature/github-distribution` resolved exact
  commit `ea351124da9de13ed9a13efbb50916fb8b4484c6` and passed `go test -mod=readonly ./...`.
- A fresh external consumer compiled the public package through both direct VCS
  resolution and `https://proxy.golang.org,direct` at
  `v0.0.0-20260903060720-ea351124da9d`.
- Complete Forgejo and GitHub advertised head/tag ref sets matched after the
  candidate push.
- A fresh public GitHub `main` clone resolved
  `5cec87847838d053f336d001ae3b04808a135dc3`, produced no `go mod tidy` drift, and passed
  `go test -mod=readonly ./...`.
- Fresh external consumers resolved `@main` through direct VCS and
  `https://proxy.golang.org,direct`, then compiled at
  `v0.0.0-20260903062630-5cec87847838`.

The slash-containing feature ref is accepted by direct VCS resolution but is
not a legal version query for the module proxy. The pre-promotion proxy lane
therefore used the exact candidate pseudo-version above; both final `@main`
lanes passed after promotion.

## Impact graph

Graphify recorded 50 nodes / 118 edges at implementation commit. Graph SHA-256:
`413750492eb3482cc731f8a0d32faa3f55c3b6597ca8f036b6c97c14ba5f33fb`.
Subsequent commits before this record changed documentation only. The merged
suite graph had 4,116 nodes and 11,415 edges, with no
cross-repository module dependency edge.

## Admission and residual gates

Two cold Judge passes reviewed the completed candidate before promotion. This
completion update changes evidence and workflow state only and receives two
fresh cold passes before commit. No performance benchmark applies because
executable paths and data flow are unchanged.

No license was selected. Release tags remain blocked until Danny closes that
decision gate. GitHub metadata mutation lacks authentication. Forgejo is still
private, so unauthenticated public contribution and private vulnerability
reporting remain unresolved. Account conversion and ownership changes were not
performed.
