# gotth-infrastructure

`gotth-infrastructure` is a deployment-contract library and inspection tool for
small containerized services. It renders digest-pinned, nonroot, read-only,
capability-free, no-new-privileges, loopback-only Compose configuration with
actual file-backed Compose secrets. It verifies the declared service against
effective Docker-compatible inspect JSON, including the image, identity,
published port, temporary filesystem, and read-only secret mounts.

It does not deploy, restart, recreate, restore, or mutate a host. Generated
configuration is desired state; runtime inspection is the evidence.

The accepted GOTTH Board alpha.2 container/PostgreSQL files are retained only
as provenance under `examples/gotth-bb`.

The contract is deliberately narrow: one application container, one TCP port,
one bounded `/tmp`, and zero or more externally supplied secret files. It does
not model databases, networks, health checks, restart policy, orchestration,
deployment, or arbitrary volume mounts. Consumers own those decisions.
