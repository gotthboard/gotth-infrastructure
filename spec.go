// Package ops defines explicit, inspectable deployment contracts for small
// containerized services. It never performs a deployment.
package ops

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	slugPattern        = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
	environmentPattern = regexp.MustCompile(`^[A-Z][A-Z0-9_]{0,79}$`)
	digestPattern      = regexp.MustCompile(`^[a-z0-9]+(?:[._:/-][a-z0-9]+)*@sha256:[0-9a-f]{64}$`)
)

// Secret is one externally managed file exposed through the Compose secrets
// contract. SourceEnv names the environment variable containing the host file
// path. Path is the exact read-only destination in the container.
type Secret struct {
	Name      string `json:"name"`
	SourceEnv string `json:"source_env"`
	Path      string `json:"path"`
}

// Service is a least-privilege runtime contract for one application container.
type Service struct {
	Name          string   `json:"name"`
	Image         string   `json:"image"`
	UID           int      `json:"uid"`
	GID           int      `json:"gid"`
	HostPort      int      `json:"host_port"`
	ContainerPort int      `json:"container_port"`
	Secrets       []Secret `json:"secrets"`
}

// Validate rejects mutable images, root identity, broad port exposure, and
// ambiguous secret mounts before rendering deployment configuration.
func (service Service) Validate() error {
	if len(service.Name) > 80 || !slugPattern.MatchString(service.Name) {
		return fmt.Errorf("service name is invalid")
	}
	if !digestPattern.MatchString(service.Image) {
		return fmt.Errorf("service image must be pinned by lowercase SHA-256 digest")
	}
	if service.UID <= 0 || service.GID <= 0 {
		return fmt.Errorf("service identity must be nonroot")
	}
	if service.HostPort < 1 || service.HostPort > 65535 || service.ContainerPort < 1 || service.ContainerPort > 65535 {
		return fmt.Errorf("service ports are invalid")
	}
	seenNames := make(map[string]struct{}, len(service.Secrets))
	seenSources := make(map[string]struct{}, len(service.Secrets))
	seenPaths := make(map[string]struct{}, len(service.Secrets))
	for _, secret := range service.Secrets {
		targetName := strings.TrimPrefix(secret.Path, "/run/secrets/")
		if len(secret.Name) > 80 || !slugPattern.MatchString(secret.Name) || !environmentPattern.MatchString(secret.SourceEnv) ||
			targetName == secret.Path || !slugPattern.MatchString(targetName) {
			return fmt.Errorf("secret mount is invalid")
		}
		if _, exists := seenNames[secret.Name]; exists {
			return fmt.Errorf("secret name is duplicated")
		}
		if _, exists := seenPaths[secret.Path]; exists {
			return fmt.Errorf("secret target is duplicated")
		}
		if _, exists := seenSources[secret.SourceEnv]; exists {
			return fmt.Errorf("secret source environment variable is duplicated")
		}
		seenNames[secret.Name] = struct{}{}
		seenSources[secret.SourceEnv] = struct{}{}
		seenPaths[secret.Path] = struct{}{}
	}
	return nil
}
