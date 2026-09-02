# Implementation specification

- Images require an explicit lowercase SHA-256 digest.
- UID and GID must both be positive.
- Ports must be within 1..65535 and render on `127.0.0.1` only.
- Root filesystem is read-only with a bounded noexec/nosuid/nodev `/tmp`.
- Privileged mode is false, all capabilities are dropped, and new privileges
  are disabled.
- Secret names are bounded slugs; targets are unique paths under `/run/secrets`.
- Inspect input is bounded to one MiB and exactly one container record.
- Unknown inspect fields are ignored because real engines emit many unrelated
  fields; every security-relevant field in the contract is checked explicitly.
