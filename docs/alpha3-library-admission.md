# Alpha.3 library admission

## Scope and authority

This pass admits `pkg/ops` as the canonical package for GOTTH Board alpha.3.
It performs no deployment and cannot invoke Docker, read secret files, restart
services, restore data, or mutate a host.

## Requirement traceability

| Requirement | Design/specification | Code | Verification |
|---|---|---|---|
| `INF-A3-01` | architecture and README layout | `pkg/ops/` | canonical outside-package test |
| `INF-A3-02` | implementation specification | `pkg/ops/` and `cmd/` | build and import inspection |
| `INF-A3-03` | architecture trust boundary | renderer/validator only | CLI and negative tests |
| `INF-A3-04` | implementation and backup contracts | inspect/compose code | Compose and disposable runtime gates |
| `INF-A3-05` | verification contract | tests/workflow evidence | clean clone, graph, two Judge passes |

## Runtime boundary

- Go 1.26.6, Docker Engine 29.7.1 admission evidence, and the installed Docker
  Compose implementation used by `make verify-compose`.
- Digest image identity, UID:GID, loopback port, read-only root, capabilities,
  no-new-privileges, bounded tmpfs, and exact secret mounts are completeness
  oracles; valid YAML alone is insufficient.
- Inputs are capped at one MiB. Disposable runtime verification uses unique
  names and empty temporary secret files and must remove its state afterward.
- Other engines require their own normalization and effective-runtime evidence.

## Performance admission

No rendering, parsing, inspection, Docker round-trip, or allocation mechanism
changes. The canonical package is the original implementation; the root owns
governance only. No speedup is claimed; benchmark/Amdahl evidence is N/A for
this structural admission.

## Failure and rollback

Rollback is a revert before the first consumer pin. Generated configuration
remains desired state only. No live container, secret, network, volume,
database, or deployment is touched.
