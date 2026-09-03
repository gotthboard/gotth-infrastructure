package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestRunRenderAndVerifyInspect(t *testing.T) {
	directory := t.TempDir()
	spec := filepath.Join(directory, "service.json")
	if err := os.WriteFile(spec, []byte(`{"name":"example-app","image":"registry.example/app@sha256:0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef","uid":10001,"gid":10001,"host_port":18080,"container_port":8080}`), 0o600); err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	if err := run([]string{"render", spec}, &output); err != nil || output.Len() == 0 {
		t.Fatalf("render = (%d bytes, %v)", output.Len(), err)
	}
	inspect := filepath.Join(directory, "inspect.json")
	if err := os.WriteFile(inspect, []byte(`[{"Id":"ignored","Config":{"User":"10001"},"HostConfig":{"ReadonlyRootfs":true,"CapDrop":["ALL"],"SecurityOpt":["no-new-privileges"]}}]`), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := run([]string{"verify-inspect", inspect}, &bytes.Buffer{}); err != nil {
		t.Fatalf("verify-inspect returned error: %v", err)
	}
	for _, arguments := range [][]string{nil, {"render", filepath.Join(directory, "missing")}, {"unknown", spec}} {
		if err := run(arguments, &bytes.Buffer{}); err == nil {
			t.Errorf("run(%v) returned no error", arguments)
		}
	}
}
