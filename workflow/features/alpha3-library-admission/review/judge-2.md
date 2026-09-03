# Judge pass 2 — clean

Reviewed revision: `7da5ddc7fe8852feb8c1fa1bde5e376f277276a5`.

All library and white-box test files are exact 100% renames into `pkg/ops`.
The CLI has one semantic edit: its import now names the canonical package.
The root has no Go package, the public package is unique, and desired-state
rendering remains separate from operator-owned execution.

Verdict: CLEAN.
