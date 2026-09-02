# Backup and restore contract

Every consumer must define, test, and record its own state boundary. A valid
procedure identifies the exact database/storage version, produces a checksum,
verifies restoration into an isolated target, records application/schema
compatibility, states whether later writes are lost, and requires fresh
confirmation before destructive restore.

This repository deliberately contains no generic restore command. Storage
engines, credentials, retention, and loss boundaries are product facts. Hiding
them behind one magical command would be operational fraud.
