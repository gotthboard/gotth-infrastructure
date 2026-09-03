package ops

import (
	"encoding/json"
	"strings"
	"testing"
)

func validInspectedContainer() inspectedContainer {
	service := validService()
	var container inspectedContainer
	container.Config.User = "10001:10001"
	container.Config.Image = service.Image
	container.HostConfig.ReadonlyRootfs = true
	container.HostConfig.CapDrop = []string{"ALL"}
	container.HostConfig.SecurityOpt = []string{"no-new-privileges:true"}
	container.HostConfig.PortBindings = map[string][]inspectedPortBinding{
		"8080/tcp": {{HostIP: "127.0.0.1", HostPort: "18080"}},
	}
	container.HostConfig.Tmpfs = map[string]string{"/tmp": "rw,noexec,nosuid,nodev,size=16777216"}
	container.Mounts = []inspectedMount{{Type: "bind", Destination: "/run/secrets/database-url", RW: false}}
	return container
}

func encodeInspection(t *testing.T, container inspectedContainer) []byte {
	t.Helper()
	raw, err := json.Marshal([]inspectedContainer{container})
	if err != nil {
		t.Fatalf("json.Marshal() returned error: %v", err)
	}
	return raw
}

func TestValidateDockerInspectMatchesDeclaredService(t *testing.T) {
	if err := ValidateDockerInspect(encodeInspection(t, validInspectedContainer()), validService()); err != nil {
		t.Fatalf("valid inspection rejected: %v", err)
	}
	for name, mutate := range map[string]func(*inspectedContainer){
		"image":      func(c *inspectedContainer) { c.Config.Image += "wrong" },
		"identity":   func(c *inspectedContainer) { c.Config.User = "10002:10001" },
		"rootfs":     func(c *inspectedContainer) { c.HostConfig.ReadonlyRootfs = false },
		"privileged": func(c *inspectedContainer) { c.HostConfig.Privileged = true },
		"cap add":    func(c *inspectedContainer) { c.HostConfig.CapAdd = []string{"NET_ADMIN"} },
		"cap drop":   func(c *inspectedContainer) { c.HostConfig.CapDrop = nil },
		"escalation": func(c *inspectedContainer) { c.HostConfig.SecurityOpt = nil },
		"port key": func(c *inspectedContainer) {
			c.HostConfig.PortBindings = map[string][]inspectedPortBinding{"8081/tcp": {{HostIP: "127.0.0.1", HostPort: "18080"}}}
		},
		"host IP":   func(c *inspectedContainer) { c.HostConfig.PortBindings["8080/tcp"][0].HostIP = "127.0.0.2" },
		"host port": func(c *inspectedContainer) { c.HostConfig.PortBindings["8080/tcp"][0].HostPort = "18081" },
		"extra port": func(c *inspectedContainer) {
			c.HostConfig.PortBindings["8081/tcp"] = []inspectedPortBinding{{HostIP: "127.0.0.1", HostPort: "18081"}}
		},
		"tmpfs missing": func(c *inspectedContainer) { c.HostConfig.Tmpfs = nil },
		"tmpfs extra":   func(c *inspectedContainer) { c.HostConfig.Tmpfs["/other"] = "rw" },
		"tmpfs unsafe":  func(c *inspectedContainer) { c.HostConfig.Tmpfs["/tmp"] = "rw,size=16m" },
		"mount missing": func(c *inspectedContainer) { c.Mounts = nil },
		"mount type":    func(c *inspectedContainer) { c.Mounts[0].Type = "volume" },
		"mount writable": func(c *inspectedContainer) {
			c.Mounts[0].RW = true
		},
		"mount target": func(c *inspectedContainer) { c.Mounts[0].Destination = "/run/secrets/other" },
		"undeclared mount": func(c *inspectedContainer) {
			c.Mounts = append(c.Mounts, inspectedMount{Type: "bind", Destination: "/data"})
		},
	} {
		container := validInspectedContainer()
		mutate(&container)
		if err := ValidateDockerInspect(encodeInspection(t, container), validService()); err == nil {
			t.Errorf("invalid %s inspection passed", name)
		}
	}
}

func TestValidateDockerInspectRejectsInvalidDocumentsAndService(t *testing.T) {
	for _, raw := range [][]byte{
		nil,
		[]byte(`[]`),
		[]byte(`{}`),
		[]byte(`[{"Config":{}}] [{"Config":{}}]`),
		[]byte(strings.Repeat(" ", (1<<20)+1)),
	} {
		if err := ValidateDockerInspect(raw, validService()); err == nil {
			t.Errorf("invalid inspection of %d bytes passed", len(raw))
		}
	}
	if err := ValidateDockerInspect(encodeInspection(t, validInspectedContainer()), Service{}); err == nil {
		t.Fatal("invalid declared service passed")
	}
}

func TestValidateDockerInspectRejectsDuplicateEffectiveTargets(t *testing.T) {
	service := validService()
	service.Secrets = append(service.Secrets, Secret{Name: "other", SourceEnv: "OTHER_FILE", Path: "/run/secrets/other"})
	container := validInspectedContainer()
	container.Mounts = append(container.Mounts, container.Mounts[0])
	if err := ValidateDockerInspect(encodeInspection(t, container), service); err == nil {
		t.Fatal("duplicate effective secret target passed")
	}
}

func TestInspectOptionHelpers(t *testing.T) {
	if !hasRequiredTmpfsOptions("nodev,size=16m,nosuid,rw,noexec") ||
		hasRequiredTmpfsOptions("rw,noexec,nosuid,nodev,size=8m") ||
		hasRequiredTmpfsOptions("rw,noexec,nosuid,nodev,size=16m,exec") ||
		hasRequiredTmpfsOptions("rw,noexec,nosuid,nodev,nodev") ||
		hasRequiredTmpfsOptions("rw,noexec,nosuid,unexpected,size=16m") {
		t.Fatal("tmpfs option validation returned the wrong result")
	}
}
