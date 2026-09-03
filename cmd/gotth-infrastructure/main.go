package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"git.dannyhunn.com/agents/gotth-infrastructure"
)

const maxInputBytes = 1 << 20

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(arguments []string, output io.Writer) error {
	if output == nil {
		return fmt.Errorf("output writer is required")
	}
	if len(arguments) < 1 ||
		(arguments[0] == "render" && len(arguments) != 2) ||
		(arguments[0] == "verify-inspect" && len(arguments) != 3) ||
		(arguments[0] != "render" && arguments[0] != "verify-inspect") {
		return fmt.Errorf("usage: gotth-infrastructure render <service.json> | verify-inspect <service.json> <inspect.json>")
	}
	serviceRaw, err := readBoundedFile(arguments[1])
	if err != nil {
		return fmt.Errorf("read service: %w", err)
	}
	service, err := decodeService(serviceRaw)
	if err != nil {
		return err
	}
	if arguments[0] == "verify-inspect" {
		inspectRaw, readErr := readBoundedFile(arguments[2])
		if readErr != nil {
			return fmt.Errorf("read inspection: %w", readErr)
		}
		return ops.ValidateDockerInspect(inspectRaw, service)
	}
	compose, err := ops.RenderCompose(service)
	if err != nil {
		return err
	}
	written, err := output.Write(compose)
	if err == nil && written != len(compose) {
		return io.ErrShortWrite
	}
	return err
}

func decodeService(raw []byte) (ops.Service, error) {
	var service ops.Service
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&service); err != nil {
		return ops.Service{}, fmt.Errorf("decode service: %w", err)
	}
	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		return ops.Service{}, fmt.Errorf("decode service: trailing JSON value")
	}
	return service, nil
}

func readBoundedFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	raw, err := io.ReadAll(io.LimitReader(file, maxInputBytes+1))
	if err != nil {
		return nil, err
	}
	if len(raw) > maxInputBytes {
		return nil, fmt.Errorf("input exceeds %d bytes", maxInputBytes)
	}
	return raw, nil
}
