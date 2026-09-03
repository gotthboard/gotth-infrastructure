# Architecture

The canonical public implementation lives in `pkg/ops`; the module root
contains no Go package and the CLI imports `pkg/ops` directly.

One strict JSON service model renders a deterministic Compose application
service. Each declared secret names an environment variable containing its
source-file path and an exact target under `/run/secrets`; rendering uses the
Compose `secrets` model, never a volume mislabeled as a secret.

The inspect validator receives both the service model and the effective
container state. It checks exact image reference, UID:GID, root filesystem,
privilege, capabilities, security options, published port, bounded `/tmp`, and
the complete mount set. Every declared secret must be one read-only bind mount
at its declared target, and undeclared mounts are rejected. This is the only
honest runtime evidence available from Docker inspect for file-backed Compose
secrets; secret contents and host paths are never read or emitted.

The CLI exposes only `render <service.json>` and
`verify-inspect <service.json> <inspect.json>`. Inputs are bounded and strictly
decoded. Neither operation invokes Docker.

Application and PostgreSQL lifecycles remain separate. This module does not
invent a universal database topology or copy live state, credentials, volumes,
addresses, or backup destinations into source control.
