# Architecture

One strict JSON service model renders a deterministic Compose application
service. An independent inspect validator reads the effective container state
and checks user, root filesystem, privilege, capabilities, security options,
and port bindings. The CLI exposes only `render` and `verify-inspect`.

Application and PostgreSQL lifecycles remain separate. This module does not
invent a universal database topology or copy live state, credentials, volumes,
addresses, or backup destinations into source control.
