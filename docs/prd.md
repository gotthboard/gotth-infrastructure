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

## Alpha.3 admission requirements

- `INF-A3-01`: New consumers import the documented `pkg/ops` package.
- `INF-A3-02`: The CLI imports the sole public package and the module root owns
  no Go implementation or runtime mechanism.
- `INF-A3-03`: Layout work does not add deployment, Docker invocation, secret
  reads, host-path disclosure, or host mutation.
- `INF-A3-04`: Compose normalization and disposable Docker inspection remain
  mandatory runtime-boundary evidence.
- `INF-A3-05`: Clean-clone, race, canonical-consumer, graph, and two clean Judge
  passes gate alpha.3 admission.
