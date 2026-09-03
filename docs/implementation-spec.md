# Implementation specification

- Images require an explicit lowercase SHA-256 digest.
- UID and GID must both be positive.
- Ports must be within 1..65535 and render on `127.0.0.1` only.
- Root filesystem is read-only with a bounded noexec/nosuid/nodev `/tmp`.
- Privileged mode is false, all capabilities are dropped, and new privileges
  are disabled.
- Secret names are bounded slugs. Source environment variables are bounded
  uppercase identifiers. Targets are unique canonical single-file paths
  directly under `/run/secrets`.
- Service-level and top-level Compose `secrets` entries use long syntax and a
  required environment-variable file source. Named volumes are never used for
  secret material.
- Inspect input is bounded to one MiB and exactly one container record.
- Inspect validation requires the declared service model and compares the
  exact image, UID:GID, one TCP port mapping, `/tmp` options, and the complete
  read-only secret mount set. Missing and undeclared mounts fail closed.
- CLI service and inspect inputs are each bounded to one MiB, reject trailing
  JSON values, and reject unknown service fields.
- Unknown inspect fields are ignored because real engines emit many unrelated
  fields; every security-relevant field in the contract is checked explicitly.
- File-backed Compose secrets cannot portably remap ownership. Operators must
  make each source file readable by the declared container UID without making
  it broadly readable; this library does not inspect or chmod secret sources.
