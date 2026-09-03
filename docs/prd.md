# Product requirements

Provide reusable least-privilege container deployment contracts across
projects. Images must be immutable by digest. Processes must be nonroot.
Root filesystems must be read-only. Capabilities and privilege escalation must
be disabled. Exactly one published TCP port must bind the declared loopback
port. Secrets must use the Compose `secrets` contract, originate from
caller-selected environment-variable file paths, and appear as exact read-only
runtime mounts under `/run/secrets`.

The tool must verify effective runtime inspection rather than declaring a
source manifest sufficient. It must never deploy or mutate infrastructure on
its own.

The first reusable contract covers one application container. Database
topology, persistent volumes, arbitrary bind mounts, health policy, network
policy, deployment, and secret creation remain consumer/operator concerns.

Acceptance requires deterministic rendering, strict input decoding, exact
desired-versus-effective inspection, Compose parser validation, a disposable
runtime proof, race repetition, clean-clone verification, and a code-graph
audit.
