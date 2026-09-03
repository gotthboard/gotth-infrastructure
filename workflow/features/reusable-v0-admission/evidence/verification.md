# Reusable v0 infrastructure admission

Verified implementation commit
`8063cdd31c4c431a7423f6498ef4bbf20b3ffea5` on 2026-09-03 UTC.

## Contract evidence

- `Secret` now declares a bounded source environment variable and exact target
  directly under `/run/secrets`. Rendering emits service-level and top-level
  Compose `secrets`; it emits no named volume masquerading as a secret.
- `ValidateDockerInspect` accepts the declared `Service` and one bounded inspect
  document. It compares the exact image, UID:GID, one `127.0.0.1` TCP mapping,
  read-only root, privilege/capability boundary, no-new-privileges option,
  exact bounded `/tmp`, and complete read-only bind-mount set.
- Missing, extra, duplicate, writable, wrong-type, and wrong-target mounts fail
  closed. Unknown inspect fields remain accepted because the engine emits many
  fields outside this contract.
- CLI service and inspect files are independently bounded to one MiB. Service
  JSON rejects unknown fields and both inputs reject trailing JSON values.
- The library and CLI never invoke Docker, resolve source environment
  variables, read secret content, deploy, or mutate a host. The operator owns
  every external action.

## Verification

- Go toolchain: `go1.26.6-X:nodwarf5`.
- `make verify`: pass; formatting, vet, race, public-library 100.0% statement
  coverage, CLI 93.2% statement coverage, and normalized Compose semantics.
- The only uncovered CLI statements are `main`'s process exit path. The
  subprocess Compose gate exercises the executable; forcing `os.Exit` inside a
  unit process would terminate the test runner and adds no contract evidence.
- `go test -mod=readonly -race -count=50 ./...`: pass.
- Clean local clone of the committed feature branch followed by `make verify`:
  pass with no generated worktree changes.
- Docker Compose 5.4.0 normalized the rendered example to exactly one `/tmp`
  entry, two service secrets, and two top-level file-backed secrets.

## Disposable runtime evidence

- Host: `development`; Docker Engine 29.7.1; Docker Compose 5.4.0.
- Image: locally present
  `alpine@sha256:25109184c71bdad752c8312a8623239686a9a2071e8825f20acb8f2198c3f659`.
- A uniquely named Compose project launched one container as `10001:10001`
  with read-only root, no privilege, no added capabilities, `cap_drop=ALL`,
  no-new-privileges, one loopback port, one bounded `/tmp`, and one empty
  disposable file-backed secret. The committed CLI accepted the actual Docker
  inspect document against the declared service.
- The first runtime attempt exposed malformed inline tmpfs YAML: commas became
  separate list entries and Docker rejected `nodev` as a mount path. Rendering
  now uses block-list syntax, and the local Compose gate asserts normalized
  tmpfs semantics rather than merely accepting parse success.
- A second attempt proved Docker records `/tmp` in `HostConfig.Tmpfs` but not
  `Mounts`, and records a read-only Compose secret with `RW=false` and empty
  `Mode`. The validator now checks those documented effective fields instead
  of requiring nonexistent evidence.
- The final runtime verification passed. Its container, network, temporary
  source file, binary, service document, Compose document, and inspect document
  were removed. A separate post-cleanup query found no matching container,
  network, or temporary directory. Existing containers were never addressed.

## Graph evidence

- Graphify 0.9.32 code-only graph: 49 nodes, 83 directed edges, 8 communities,
  no self-loops, exact duplicate edges, or same-endpoint relation groups.
- Graph SHA-256:
  `4a9fc0108ee8da6aa196e7c865ba6cc94e1a40729ad697876dcf157aa1c5992a`.
- Graph cache:
  `/home/linus/.cache/openclaw-graphify/gotth-infrastructure-reusable/graphify-out/graph.json`.
- No sensitive file was read. Graphify skipped non-code manifests and lacked
  the optional SQL parser for the provenance-only runtime-grants example; the
  consequential Go source, tests, Compose output, and actual runtime state were
  verified directly. Extraction changed no repository file.

No tag or consumer pin was created. The first compatibility promise belongs to
a real consumer integration.
