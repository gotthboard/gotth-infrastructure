# Product requirements

Provide reusable least-privilege container deployment contracts across
projects. Images must be immutable by digest. Processes must be nonroot.
Root filesystems must be read-only. Capabilities and privilege escalation must
be disabled. Published ports must bind loopback. Secrets must arrive through
explicit read-only external mounts.

The tool must verify effective runtime inspection rather than declaring a
source manifest sufficient. It must never deploy or mutate infrastructure on
its own.
