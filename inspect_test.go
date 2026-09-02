package ops

import "testing"

const validInspection = `[{"Config":{"User":"10001:10001"},"HostConfig":{"ReadonlyRootfs":true,"Privileged":false,"CapAdd":null,"CapDrop":["ALL"],"SecurityOpt":["no-new-privileges:true"],"PortBindings":{"8080/tcp":[{"HostIp":"127.0.0.1"}]}}}]`

func TestValidateDockerInspect(t *testing.T) {
	if err := ValidateDockerInspect([]byte(validInspection)); err != nil {
		t.Fatalf("valid inspection rejected: %v", err)
	}
	for _, raw := range []string{
		``, `[]`, `{}`,
		`[{"Config":{"User":"0"},"HostConfig":{"ReadonlyRootfs":true,"CapDrop":["ALL"],"SecurityOpt":["no-new-privileges"]}}]`,
		`[{"Config":{"User":"10001"},"HostConfig":{"ReadonlyRootfs":false,"CapDrop":["ALL"],"SecurityOpt":["no-new-privileges"]}}]`,
		`[{"Config":{"User":"10001"},"HostConfig":{"ReadonlyRootfs":true,"CapDrop":[],"SecurityOpt":["no-new-privileges"]}}]`,
		`[{"Config":{"User":"10001"},"HostConfig":{"ReadonlyRootfs":true,"CapDrop":["ALL"],"SecurityOpt":[]}}]`,
		`[{"Config":{"User":"10001"},"HostConfig":{"ReadonlyRootfs":true,"CapDrop":["ALL"],"SecurityOpt":["no-new-privileges"],"PortBindings":{"8080/tcp":[{"HostIp":"0.0.0.0"}]}}}]`,
	} {
		if err := ValidateDockerInspect([]byte(raw)); err == nil {
			t.Errorf("invalid inspection %q passed", raw)
		}
	}
}
