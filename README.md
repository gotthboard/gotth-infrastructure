# gotth-infrastructure

`gotth-infrastructure` is a deployment-contract library and inspection tool for small
containerized services. It renders digest-pinned, nonroot, read-only,
capability-free, no-new-privileges, loopback-only Compose configuration and
verifies the effective runtime boundary from Docker-compatible inspect JSON.

It does not deploy, restart, recreate, restore, or mutate a host. Generated
configuration is desired state; runtime inspection is the evidence.

The accepted GOTTH Board alpha.2 container/PostgreSQL files are retained only
as provenance under `examples/gotth-bb`.
