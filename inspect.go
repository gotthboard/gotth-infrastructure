package ops

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type inspectedPortBinding struct {
	HostIP   string `json:"HostIp"`
	HostPort string `json:"HostPort"`
}

type inspectedMount struct {
	Type        string `json:"Type"`
	Destination string `json:"Destination"`
	RW          bool   `json:"RW"`
}

type inspectedContainer struct {
	Config struct {
		User  string `json:"User"`
		Image string `json:"Image"`
	} `json:"Config"`
	HostConfig struct {
		ReadonlyRootfs bool                              `json:"ReadonlyRootfs"`
		Privileged     bool                              `json:"Privileged"`
		CapAdd         []string                          `json:"CapAdd"`
		CapDrop        []string                          `json:"CapDrop"`
		SecurityOpt    []string                          `json:"SecurityOpt"`
		PortBindings   map[string][]inspectedPortBinding `json:"PortBindings"`
		Tmpfs          map[string]string                 `json:"Tmpfs"`
	} `json:"HostConfig"`
	Mounts []inspectedMount `json:"Mounts"`
}

// ValidateDockerInspect verifies that one effective Docker-compatible runtime
// exactly implements the declared service boundary.
func ValidateDockerInspect(raw []byte, service Service) error {
	if err := service.Validate(); err != nil {
		return fmt.Errorf("declared service is invalid: %w", err)
	}
	if len(raw) == 0 || len(raw) > 1<<20 {
		return fmt.Errorf("container inspection has an invalid size")
	}
	var containers []inspectedContainer
	decoder := json.NewDecoder(bytes.NewReader(raw))
	if err := decoder.Decode(&containers); err != nil || len(containers) != 1 {
		return fmt.Errorf("container inspection is invalid")
	}
	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		return fmt.Errorf("container inspection is invalid")
	}
	container := containers[0]
	wantUser := strconv.Itoa(service.UID) + ":" + strconv.Itoa(service.GID)
	if container.Config.User != wantUser {
		return fmt.Errorf("container identity does not match the declared service")
	}
	if container.Config.Image != service.Image {
		return fmt.Errorf("container image does not match the declared service")
	}
	if !container.HostConfig.ReadonlyRootfs || container.HostConfig.Privileged || len(container.HostConfig.CapAdd) != 0 {
		return fmt.Errorf("container privilege boundary is invalid")
	}
	droppedAll := false
	for _, capability := range container.HostConfig.CapDrop {
		if strings.EqualFold(capability, "ALL") {
			droppedAll = true
		}
	}
	if !droppedAll {
		return fmt.Errorf("container does not drop all capabilities")
	}
	noNewPrivileges := false
	for _, option := range container.HostConfig.SecurityOpt {
		if option == "no-new-privileges" || option == "no-new-privileges:true" {
			noNewPrivileges = true
		}
	}
	if !noNewPrivileges {
		return fmt.Errorf("container permits privilege escalation")
	}
	wantPortKey := strconv.Itoa(service.ContainerPort) + "/tcp"
	bindings, exists := container.HostConfig.PortBindings[wantPortKey]
	if !exists || len(container.HostConfig.PortBindings) != 1 || len(bindings) != 1 {
		return fmt.Errorf("container port publication does not match the declared service")
	}
	binding := bindings[0]
	if binding.HostIP != "127.0.0.1" || binding.HostPort != strconv.Itoa(service.HostPort) {
		return fmt.Errorf("container port publication does not match the declared service")
	}
	if len(container.HostConfig.Tmpfs) != 1 {
		return fmt.Errorf("container temporary filesystem does not match the declared service")
	}
	tmpfsOptions, exists := container.HostConfig.Tmpfs["/tmp"]
	if !exists || !hasRequiredTmpfsOptions(tmpfsOptions) {
		return fmt.Errorf("container temporary filesystem does not match the declared service")
	}
	if len(container.Mounts) != len(service.Secrets) {
		return fmt.Errorf("container secret mounts do not match the declared service")
	}
	wantedSecrets := make(map[string]struct{}, len(service.Secrets))
	for _, secret := range service.Secrets {
		wantedSecrets[secret.Path] = struct{}{}
	}
	seenSecrets := make(map[string]struct{}, len(container.Mounts))
	for _, mount := range container.Mounts {
		if _, wanted := wantedSecrets[mount.Destination]; !wanted || mount.Type != "bind" || mount.RW {
			return fmt.Errorf("container secret mounts do not match the declared service")
		}
		if _, duplicate := seenSecrets[mount.Destination]; duplicate {
			return fmt.Errorf("container secret mounts do not match the declared service")
		}
		seenSecrets[mount.Destination] = struct{}{}
	}
	return nil
}

func hasRequiredTmpfsOptions(raw string) bool {
	parts := strings.Split(raw, ",")
	if len(parts) != 5 {
		return false
	}
	options := make(map[string]struct{})
	for _, option := range parts {
		if _, duplicate := options[option]; duplicate {
			return false
		}
		options[option] = struct{}{}
	}
	for _, required := range []string{"rw", "noexec", "nosuid", "nodev"} {
		if _, exists := options[required]; !exists {
			return false
		}
	}
	_, bytesSize := options["size=16777216"]
	_, humanSize := options["size=16m"]
	return bytesSize || humanSize
}
