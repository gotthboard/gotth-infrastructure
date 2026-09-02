package ops

import (
	"bytes"
	"fmt"
	"sort"
)

// RenderCompose returns deterministic Compose YAML for one isolated app. It
// binds only loopback, drops every capability, forbids privilege escalation,
// and mounts caller-declared external secrets read-only.
func RenderCompose(service Service) ([]byte, error) {
	if err := service.Validate(); err != nil {
		return nil, err
	}
	secrets := append([]Secret(nil), service.Secrets...)
	sort.Slice(secrets, func(left, right int) bool { return secrets[left].Name < secrets[right].Name })
	var output bytes.Buffer
	_, _ = fmt.Fprintf(&output, "name: %s\nservices:\n  app:\n    image: %s\n    user: \"%d:%d\"\n    read_only: true\n    privileged: false\n    cap_drop: [ALL]\n    security_opt: [no-new-privileges:true]\n    ports: [\"127.0.0.1:%d:%d\"]\n    tmpfs: [/tmp:rw,noexec,nosuid,nodev,size=16m]\n", service.Name, service.Image, service.UID, service.GID, service.HostPort, service.ContainerPort)
	if len(secrets) > 0 {
		output.WriteString("    volumes:\n")
		for _, secret := range secrets {
			_, _ = fmt.Fprintf(&output, "      - %s:%s:ro\n", secret.Name, secret.Path)
		}
		output.WriteString("volumes:\n")
		for _, secret := range secrets {
			_, _ = fmt.Fprintf(&output, "  %s:\n    external: true\n", secret.Name)
		}
	}
	return output.Bytes(), nil
}
