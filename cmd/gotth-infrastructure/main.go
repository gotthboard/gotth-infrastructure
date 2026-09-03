package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"git.dannyhunn.com/agents/gotth-infrastructure"
)

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(arguments []string, output io.Writer) error {
	if len(arguments) != 2 || (arguments[0] != "render" && arguments[0] != "verify-inspect") {
		return fmt.Errorf("usage: gotth-infrastructure render <service.json> | verify-inspect <inspect.json>")
	}
	raw, err := os.ReadFile(arguments[1])
	if err != nil {
		return fmt.Errorf("read input: %w", err)
	}
	if arguments[0] == "verify-inspect" {
		return ops.ValidateDockerInspect(raw)
	}
	var service ops.Service
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&service); err != nil {
		return fmt.Errorf("decode service: %w", err)
	}
	compose, err := ops.RenderCompose(service)
	if err != nil {
		return err
	}
	_, err = output.Write(compose)
	return err
}
