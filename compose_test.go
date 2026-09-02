package ops

import (
	"bytes"
	"strings"
	"testing"
)

func TestRenderComposeIsDeterministicAndLeastPrivilege(t *testing.T) {
	service := validService()
	service.Secrets = append(service.Secrets, Secret{Name: "client-secret", Path: "/run/secrets/client-secret"})
	first, err := RenderCompose(service)
	if err != nil {
		t.Fatal(err)
	}
	second, _ := RenderCompose(service)
	if !bytes.Equal(first, second) {
		t.Fatal("rendering is not deterministic")
	}
	for _, required := range []string{"read_only: true", "privileged: false", "cap_drop: [ALL]", "no-new-privileges:true", "127.0.0.1:18080:8080", ":ro", "external: true"} {
		if !strings.Contains(string(first), required) {
			t.Errorf("compose output missing %q", required)
		}
	}
	if output, err := RenderCompose(Service{}); err == nil || output != nil {
		t.Fatalf("invalid render = (%q, %v)", output, err)
	}
}
