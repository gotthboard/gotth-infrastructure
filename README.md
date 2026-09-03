# gotth-infrastructure

> **Distribution:** GitHub is the public clone, Go import, and future release endpoint.
> Forgejo remains canonical development and the issue/contribution location.
> See [the distribution contract](docs/distribution.md).


`gotth-infrastructure` is a deployment-contract library and inspection tool for
small containerized services. It renders digest-pinned, nonroot, read-only,
capability-free, no-new-privileges, loopback-only Compose configuration with
actual file-backed Compose secrets. It verifies the declared service against
effective Docker-compatible inspect JSON, including the image, identity,
published port, temporary filesystem, and read-only secret mounts.

Canonical Go package:
`github.com/gotthboard/gotth-infrastructure/pkg/ops`. The module root
contains repository governance only; the CLI imports the canonical package.

Repository layout:

- `pkg/ops/` — public desired-state and inspection implementation;
- `cmd/gotth-infrastructure/` — bounded render/verify CLI;
- module root — module metadata and repository governance;
- `examples/`, `docs/`, and `workflow/` — provenance, contracts, and canonical
  workflow state.

It does not deploy, restart, recreate, restore, or mutate a host. Generated
configuration is desired state; runtime inspection is the evidence.

The accepted GOTTH Board alpha.2 container/PostgreSQL files are retained only
as provenance under `examples/gotth-bb`.

The contract is deliberately narrow: one application container, one TCP port,
one bounded `/tmp`, and zero or more externally supplied secret files. It does
not model databases, networks, health checks, restart policy, orchestration,
deployment, or arbitrary volume mounts. Consumers own those decisions.

## Installation, compatibility, and support

Unreleased. The API and CLI are pre-1.0 and may change until the first
tagged compatibility contract.

No post-migration version has been tagged. To inspect the current source
before the first admitted release:

```sh
go get github.com/gotthboard/gotth-infrastructure@main
go install github.com/gotthboard/gotth-infrastructure/cmd/gotth-infrastructure@main
```

The repository has no selected license and no long-term support promise.
Versioning, release admission, security reporting, and contribution details are
in [the release policy](docs/RELEASING.md), [security policy](SECURITY.md), and
[contribution guide](CONTRIBUTING.md).
