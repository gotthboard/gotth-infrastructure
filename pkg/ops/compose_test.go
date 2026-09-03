package ops

import (
	"bytes"
	"strings"
	"testing"
)

func TestRenderComposeIsDeterministicAndLeastPrivilege(t *testing.T) {
	service := validService()
	service.Secrets = append(service.Secrets, Secret{Name: "client-secret", SourceEnv: "OIDC_CLIENT_SECRET_FILE", Path: "/run/secrets/client-secret"})
	first, err := RenderCompose(service)
	if err != nil {
		t.Fatal(err)
	}
	second, _ := RenderCompose(service)
	if !bytes.Equal(first, second) {
		t.Fatal("rendering is not deterministic")
	}
	for _, required := range []string{"read_only: true", "privileged: false", "cap_drop: [ALL]", "no-new-privileges:true", "127.0.0.1:18080:8080", "- /tmp:rw,noexec,nosuid,nodev,size=16m", "secrets:", "source: client-secret", "target: /run/secrets/client-secret", "file: ${OIDC_CLIENT_SECRET_FILE:?OIDC_CLIENT_SECRET_FILE is required}"} {
		if !strings.Contains(string(first), required) {
			t.Errorf("compose output missing %q", required)
		}
	}
	if strings.Contains(string(first), "volumes:") || strings.Index(string(first), "source: client-secret") > strings.Index(string(first), "source: database-url") {
		t.Fatalf("compose output contains a volume or unsorted secrets:\n%s", first)
	}
	if output, err := RenderCompose(Service{}); err == nil || output != nil {
		t.Fatalf("invalid render = (%q, %v)", output, err)
	}
}
