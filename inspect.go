package ops

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"strings"
)

type inspectedContainer struct {
	Config struct {
		User string `json:"User"`
	} `json:"Config"`
	HostConfig struct {
		ReadonlyRootfs bool     `json:"ReadonlyRootfs"`
		Privileged     bool     `json:"Privileged"`
		CapAdd         []string `json:"CapAdd"`
		CapDrop        []string `json:"CapDrop"`
		SecurityOpt    []string `json:"SecurityOpt"`
		PortBindings   map[string][]struct {
			HostIP string `json:"HostIp"`
		} `json:"PortBindings"`
	} `json:"HostConfig"`
}

// ValidateDockerInspect verifies effective runtime containment from bounded
// Docker-compatible inspect JSON rather than trusting a Compose source file.
func ValidateDockerInspect(raw []byte) error {
	if len(raw) == 0 || len(raw) > 1<<20 {
		return fmt.Errorf("container inspection has an invalid size")
	}
	var containers []inspectedContainer
	decoder := json.NewDecoder(bytes.NewReader(raw))
	if err := decoder.Decode(&containers); err != nil || len(containers) != 1 {
		return fmt.Errorf("container inspection is invalid")
	}
	container := containers[0]
	uid := strings.SplitN(container.Config.User, ":", 2)[0]
	if uid == "" || uid == "0" || uid == "root" {
		return fmt.Errorf("container runs as root or unspecified identity")
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
	for _, bindings := range container.HostConfig.PortBindings {
		for _, binding := range bindings {
			if ip := net.ParseIP(binding.HostIP); ip == nil || !ip.IsLoopback() {
				return fmt.Errorf("container publishes a non-loopback port")
			}
		}
	}
	return nil
}
