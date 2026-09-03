# Verification

`make verify` checks formatting, vet, race, coverage, and example rendering.
Tests reject mutable images, root users, invalid ports, unsafe/duplicate secret
definitions, image/identity/port mismatches, missing, writable, wrongly typed,
or undeclared mounts, unsafe `/tmp`, writable/privileged/capability-bearing
containers, privilege escalation, malformed inspect data, and non-loopback
port publication.

Reusable admission also requires external-package compilation, strict and
bounded CLI regressions, 50 race-enabled repetitions, `docker compose config`
against rendered output, a disposable effective-runtime inspection, a clean
clone, and a fresh code-graph audit.

No live container, database, proxy, secret, volume, backup, or host is touched
by this repository's test gate.

The disposable runtime gate uses unique temporary names and paths, contains no
real secret, and removes its container, network, files, and temporary checkout
after verification. It must never target an existing Compose project.
