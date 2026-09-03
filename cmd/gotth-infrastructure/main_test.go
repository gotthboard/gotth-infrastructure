package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const testImage = "registry.example/app@sha256:0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

func writeFixture(t *testing.T, directory, name, contents string) string {
	t.Helper()
	path := filepath.Join(directory, name)
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatalf("os.WriteFile(%q) returned error: %v", name, err)
	}
	return path
}

func TestRunRenderAndVerifyInspect(t *testing.T) {
	directory := t.TempDir()
	spec := writeFixture(t, directory, "service.json", `{"name":"example-app","image":"`+testImage+`","uid":10001,"gid":10001,"host_port":18080,"container_port":8080}`)
	var output bytes.Buffer
	if err := run([]string{"render", spec}, &output); err != nil || output.Len() == 0 {
		t.Fatalf("render = (%d bytes, %v)", output.Len(), err)
	}
	inspect := writeFixture(t, directory, "inspect.json", `[{"Config":{"User":"10001:10001","Image":"`+testImage+`"},"HostConfig":{"ReadonlyRootfs":true,"CapDrop":["ALL"],"SecurityOpt":["no-new-privileges"],"PortBindings":{"8080/tcp":[{"HostIp":"127.0.0.1","HostPort":"18080"}]},"Tmpfs":{"/tmp":"rw,noexec,nosuid,nodev,size=16777216"}},"Mounts":[]}]`)
	if err := run([]string{"verify-inspect", spec, inspect}, &bytes.Buffer{}); err != nil {
		t.Fatalf("verify-inspect returned error: %v", err)
	}
}

func TestRunRejectsUsageFilesAndMalformedService(t *testing.T) {
	directory := t.TempDir()
	valid := writeFixture(t, directory, "valid.json", `{"name":"example-app","image":"`+testImage+`","uid":10001,"gid":10001,"host_port":18080,"container_port":8080}`)
	unknown := writeFixture(t, directory, "unknown.json", `{"name":"example-app","unknown":true}`)
	invalid := writeFixture(t, directory, "invalid.json", `{"name":"example-app"}`)
	trailing := writeFixture(t, directory, "trailing.json", `{"name":"example-app"} {}`)
	oversize := writeFixture(t, directory, "oversize.json", strings.Repeat(" ", maxInputBytes+1))
	missing := filepath.Join(directory, "missing")
	for _, test := range []struct {
		arguments []string
		output    io.Writer
	}{
		{output: &bytes.Buffer{}},
		{arguments: []string{"unknown", valid}, output: &bytes.Buffer{}},
		{arguments: []string{"render", valid}, output: nil},
		{arguments: []string{"render", missing}, output: &bytes.Buffer{}},
		{arguments: []string{"verify-inspect", valid, missing}, output: &bytes.Buffer{}},
		{arguments: []string{"render", unknown}, output: &bytes.Buffer{}},
		{arguments: []string{"render", invalid}, output: &bytes.Buffer{}},
		{arguments: []string{"render", trailing}, output: &bytes.Buffer{}},
		{arguments: []string{"render", oversize}, output: &bytes.Buffer{}},
	} {
		if err := run(test.arguments, test.output); err == nil {
			t.Errorf("run(%v) returned no error", test.arguments)
		}
	}
	if raw, err := readBoundedFile(directory); err == nil || raw != nil {
		t.Fatalf("readBoundedFile(directory) = (%q, %v), want nil/error", raw, err)
	}
}

func TestRunReturnsOutputFailures(t *testing.T) {
	directory := t.TempDir()
	spec := writeFixture(t, directory, "service.json", `{"name":"example-app","image":"`+testImage+`","uid":10001,"gid":10001,"host_port":18080,"container_port":8080}`)
	cause := errors.New("write failed")
	if err := run([]string{"render", spec}, errorWriter{cause: cause}); !errors.Is(err, cause) {
		t.Fatalf("write failure = %v, want %v", err, cause)
	}
	if err := run([]string{"render", spec}, shortWriter{}); !errors.Is(err, io.ErrShortWrite) {
		t.Fatalf("short write = %v, want %v", err, io.ErrShortWrite)
	}
}

type errorWriter struct{ cause error }

func (writer errorWriter) Write([]byte) (int, error) { return 0, writer.cause }

type shortWriter struct{}

func (shortWriter) Write(raw []byte) (int, error) { return len(raw) - 1, nil }
