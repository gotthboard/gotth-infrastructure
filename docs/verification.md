# Verification

`make verify` checks formatting, vet, race, coverage, and example rendering.
Tests reject mutable images, root users, invalid ports, unsafe/duplicate secret
mounts, writable/privileged/capability-bearing containers, privilege
escalation, malformed inspect data, and non-loopback port publication.

No live container, database, proxy, secret, volume, backup, or host is touched
by this repository's test gate.
