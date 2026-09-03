package ops_test

import (
	"strings"
	"testing"

	ops "github.com/gotthboard/gotth-infrastructure"
)

func TestPublicAPIIsUsableOutsidePackage(t *testing.T) {
	var _ func([]byte, ops.Service) error = ops.ValidateDockerInspect
	service := ops.Service{
		Name: "consumer", Image: "registry.example/consumer@sha256:0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
		UID: 10001, GID: 10001, HostPort: 18080, ContainerPort: 8080,
		Secrets: []ops.Secret{{Name: "token", SourceEnv: "TOKEN_FILE", Path: "/run/secrets/token"}},
	}
	compose, err := ops.RenderCompose(service)
	if err != nil || !strings.Contains(string(compose), "file: ${TOKEN_FILE:?TOKEN_FILE is required}") {
		t.Fatalf("RenderCompose() = (%q, %v)", compose, err)
	}
}
