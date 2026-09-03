# Changelog

This repository records user-visible and compatibility-relevant changes here.
Released sections use Semantic Versioning; unreleased work remains under
`Unreleased` and does not imply a tag.

## Unreleased

### 2026-09-03 01:04 CDT — Structure and formally admit the alpha.3 library

Commit: `bf1a9827c29f4f737d50b13fd1e3c36335648153`

Affected files:

- `pkg/ops/`
- canonical API test and CLI imports
- `README.md`, `docs/`, `workflow.toml`, and admission evidence

Explanation:

Move the containment implementation and tests out of the repository root,
point the CLI at the canonical public package, leave the root for governance,
and add formal coding-setup admission records.

Verification:

- preliminary `go test ./...` passed after the move
- final race, Compose, runtime, clean-clone, graph, and Judge evidence is in the
  admission workflow evidence

Risks / non-goals:

- no deployment, Docker invocation by library/CLI, secret read, or host change

### 2026-09-03 00:42 CDT — Establish GitHub public distribution

Commit: `4b69d81`

Affected files:

- `README.md`
- `CONTRIBUTING.md`
- `SECURITY.md`
- `docs/distribution.md`
- `docs/RELEASING.md`
- `go.mod` and repository-owned Go import references
- repository-owned Go source, tests, fixtures, and package documentation
- `workflow.toml` and `workflow/features/github-distribution/`

Explanation:

Declare GitHub as the public distribution endpoint while retaining Forgejo as
canonical development, define maturity and support honestly, and document the
independent release process. The Go module identity and exact self-imports move to the public GitHub path.

Verification:

- exact old-import search
- documentation contract audit
- `go mod tidy` drift check
- `go vet -mod=readonly ./...`
- `go test -mod=readonly ./...`

Risks / non-goals:

- No license is selected.
- No existing tag is changed and no new release is created.
- Mirror direction, repository ownership, and account type are unchanged.
