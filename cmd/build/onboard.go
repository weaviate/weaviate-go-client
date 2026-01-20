//nolint:errcheck
package main

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	BinDir    = "./bin"
	BinREADME = `Package ./bin contains binary executables and shims for some Go tools.

This directory is git-ignored and meant to stay local to your machine.
Do not commit any of its contents to version control.

- protoc -- Protobuf compiler for generating stubs.
- include/ -- Well-known protobuf type definitions. Required for protoc, see: include/readme.txt
- oapi-codegen -- Go tool for generating client-side models from OpenAPI specs.
- golangci-lint -- Pre-compiled golangci-lint for local use.
`

	OapiCodegen     = "oapi-codegen"
	OapiCodegenShim = `#!/bin/sh

exec go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen "$@"
`
	GolangCILintVersion = "2.8.0"

	ProtocVersion = "33.2"

	OnboardHelp = `Usage: go run ./cmd/build onboard

Options:
	--help Print this message.
`
)

// Onboard bootstraps local development environment.
//   - Installs `protoc`, the protobuf compiler.
//   - Installs `golangci-lint` binary for local lint runs.
//   - Creates a shim for `oapi-codegen`.
func Onboard(ctx context.Context, args []string) error {
	var opt struct {
		Help bool
	}
	cmd := flag.NewFlagSet("onboard", flag.ExitOnError)
	cmd.BoolVar(&opt.Help, "help", false, "Print help message.")

	if err := cmd.Parse(args); err != nil {
		return err
	}

	if opt.Help {
		fmt.Print(OnboardHelp)
		return nil
	}

	if _, err := writeAtomic(filepath.Join(BinDir, "README"), strings.NewReader(BinREADME), false); err != nil {
		return fmt.Errorf("\tERROR: %w", err)
	}

	log.Print("Installing protoc")
	if err := installProtoc(ctx); err != nil {
		return fmt.Errorf("\tERROR: %w", err)
	}

	log.Print("Installing golangci-lint")
	if err := installGolangciLint(ctx); err != nil {
		return fmt.Errorf("\tERROR: %w", err)
	}

	shim := filepath.Join(BinDir, OapiCodegen)
	if _, err := writeAtomic(shim, strings.NewReader(OapiCodegenShim), false); err != nil {
		return fmt.Errorf("\tERROR: %w", err)
	}
	if err := chmodx(shim); err != nil {
		return fmt.Errorf("\tERROR: %w", err)
	}

	log.Print("Done")
	return nil
}

const (
	ProtocReleases   = "https://github.com/protocolbuffers/protobuf/releases/download/v%s/"
	ProtocDarwinZip  = "protoc-%s-osx-universal_binary.zip"
	ProtocLinuxZip   = "protoc-%s-linux-x86_64.zip"
	ProtocWindowsZip = "protoc-%s-win64.zip"
)

func installProtoc(_ context.Context) error {
	releases := fmt.Sprintf(ProtocReleases, ProtocVersion)
	var artifact string
	switch runtime.GOOS {
	case "darwin":
		artifact = fmt.Sprintf(ProtocDarwinZip, ProtocVersion)
	case "linux":
		artifact = fmt.Sprintf(ProtocLinuxZip, ProtocVersion)
	case "windows":
		artifact = fmt.Sprintf(ProtocLinuxZip, ProtocVersion)
	default:
		return fmt.Errorf("%s is not supported", runtime.GOOS)
	}

	log.Printf("Fetching release artifact %s for %q", artifact, runtime.GOOS)

	r, err := http.Get(releases + artifact)
	if err != nil {
		return fmt.Errorf("install protoc: %w", err)
	}
	defer r.Body.Close()

	b, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	log.Printf("Unzipping protoc to %s", BinDir)

	br := bytes.NewReader(b)
	zr, err := zip.NewReader(br, br.Size())
	if err != nil {
		return err
	}

	uncompress := func(f *zip.File) (int64, error) {
		rc, err := f.Open()
		if err != nil {
			return 0, err
		}
		defer rc.Close()

		name := f.Name
		switch f.Name {
		case "bin/protoc":
			name = "protoc"
		case "readme.txt":
			name = "include/" + f.Name
		}
		return writeAtomic(filepath.Join(BinDir, name), rc, true)
	}

	ok := true
	for _, f := range zr.File {
		// uncompress will create all parent directories
		if f.FileInfo().IsDir() {
			continue
		}
		if written, err := uncompress(f); err == nil {
			log.Printf("\tuncompress %s [%dB] ok", f.Name, written)
		} else {
			log.Printf("\tERROR: %v", err)
			ok = false
		}
	}

	if err := chmodx(filepath.Join(BinDir, "protoc")); err != nil {
		return err
	}

	if !ok {
		return errors.New("protoc was not installed")
	}

	return nil
}

// GolangCILintInstallCmd downloads and runs a script to install golangci-lint locally.
//
// Source: https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh
// See: https://golangci-lint.run/docs/welcome/install/local
const GolangCILintInstallCmd = "curl -sSfL https://golangci-lint.run/install.sh | sh -s v%s"

func installGolangciLint(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "sh", "-c", fmt.Sprintf(GolangCILintInstallCmd, GolangCILintVersion))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("install golanci-lint: %w", err)
	}
	return nil
}

// Give file executable permissions.
func chmodx(file string) error { return os.Chmod(file, 0o755) }
