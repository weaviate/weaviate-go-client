//go:build mage

package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/magefile/mage/mg"
)

const (
	EnvGithubToken = "GITHUB_TOKEN"

	WeaviateRoot         = "https://api.github.com/repos/weaviate/weaviate/contents"
	WeaviateProtobufs    = "grpc/proto/v1"
	WeaviateOpenAPISpecs = "openapi-specs/schema.json"

	LocalOpenAPISpecs = "./api/rest"
	LocalProtobufs    = "./api/proto/v1"

	OpenAPISchemaCheck = "schema.check.json"
	OpenAPISchemaV3    = "schema.v3.json"
)

type Contracts mg.Namespace

// Checksums "api/proto" and "api/rest" specs against their latest versions in weaviate/weaviate.
// Best used in CI scripts to ensure the sources of our generated stubs stay up-to-date.
func (Contracts) Check(ctx context.Context) error { return validate(ctx, false) }

// Updates "api/proto" and "api/rest" specs to their latest versions.
// Use locally to pull new changes or fix broken / missing / out-of-date specs.
func (Contracts) Update(ctx context.Context) error { return validate(ctx, true) }

// validate calculates hashes for the OpenAPI schema.json and the protobufs in "api/"
// and compares them against their latest versions in the weaviate/weaviate repo.
// If update returns
func validate(ctx context.Context, update bool) error {
	log.Printf("Fetching metadata for %s", WeaviateOpenAPISpecs)
	openapi, err := headOpenAPISpecs(ctx)
	if err != nil {
		return err
	}

	log.Printf("Fetching metadata for %s/*.proto", WeaviateProtobufs)
	protobufs, err := headProtobufs(ctx)
	if err != nil {
		return err
	}

	var contracts []Contract
	{
		path := filepath.Join(LocalOpenAPISpecs, OpenAPISchemaCheck)
		c, err := newContract(ctx, path, *openapi)
		if err != nil {
			return err
		}
		contracts = append(contracts, *c)
	}

	for _, file := range protobufs {
		path := filepath.Join(LocalProtobufs, file.Name)
		c, err := newContract(ctx, path, file)
		if err != nil {
			return err
		}
		contracts = append(contracts, *c)
	}

	if len(contracts) == 0 {
		return errors.New("no specs in weaviate/weaviate")
	}

	ok := true
	updated := false
	for _, file := range contracts {
		if file.SHA == file.Upstream.SHA {
			log.Printf("check %s: ok", file.Path)
			continue
		}

		log.Printf("check %s:\n\twant:\t%s\n\tgot:\t%s", file.Path, file.Upstream.SHA, file.SHA)
		if update {
			log.Print("Downloading latest ", file.Upstream.DownloadURL)
			if err := updateContract(ctx, file); err != nil {
				log.Printf("\tERROR: %s", err)
				ok = false
			} else {
				updated = true
			}
		} else {
			ok = false
		}
	}
	if !ok {
		return fmt.Errorf(`
Contracts in weaviate-go-client are out-of-sync with weaviate/weaviate repository.
Update them to the latest version by running this command:
	./bin/mage -v contracts:update
`)
	}
	if updated {
		log.Println(`
Contracts were successfully updated, run:
	go generate ./...

to re-generate REST and gRPC stubs.
`)
	}
	log.Print("Done")
	return nil
}

type Contract struct {
	Upstream GithubFile // Upstream file metadata.
	Exists   bool       // False if the file does not exist locally.
	Path     string     // Local filepath.
	SHA      string     // Local file SHA.
}

type GithubFile struct {
	Name        string `json:"name"`
	SHA         string `json:"sha"`
	DownloadURL string `json:"download_url"`
}

// headProtobufs fetches metadata for schema.json.
func headOpenAPISpecs(ctx context.Context) (*GithubFile, error) {
	dir, basename := filepath.Split(WeaviateOpenAPISpecs)
	specs, err := ghFiles(ctx, dir)
	if err != nil {
		return nil, err
	}
	for _, file := range specs {
		if file.Name == basename {
			return &file, nil
		}
	}
	return nil, fmt.Errorf("%s/%s not found", WeaviateRoot, WeaviateOpenAPISpecs)
}

// headProtobufs fetches metadata for protobuf files.
func headProtobufs(ctx context.Context) ([]GithubFile, error) {
	return ghFiles(ctx, WeaviateProtobufs)
}

// updateContract fetches the latest version of the [c.Upstream] and writes it to [c.Path].
func updateContract(ctx context.Context, c Contract) error {
	rc, err := ghGet(ctx, c.Upstream.DownloadURL)
	if err != nil {
		return err
	}
	defer rc.Close()

	if err := os.MkdirAll(filepath.Dir(c.Path), 0o775); err != nil {
		log.Print(err)
	}

	f, err := os.Create(c.Path + ".tmp")
	if err != nil {
		return err
	}
	defer f.Close()

	written, err := io.Copy(f, rc)
	if err != nil {
		return err
	}

	if err := os.Rename(f.Name(), c.Path); err != nil {
		return err
	}

	if c.Exists {
		log.Printf("Updated %s [written %dB]", c.Path, written)
	} else {
		log.Printf("Added new file at %s [written %dB]", c.Path, written)
	}

	if os.Remove(f.Name()); err != nil {
		log.Print(err)
	}
	return nil
}

func newContract(ctx context.Context, local string, upstream GithubFile) (*Contract, error) {
	if _, err := os.Stat(local); errors.Is(err, os.ErrNotExist) {
		return &Contract{
			Upstream: upstream,
			Path:     local,
			SHA:      "<file not found>",
		}, nil
	}
	// os.Stat might've failed for a different reason, still try to get the hash.
	sha, err := gitSHA(ctx, local)
	if err != nil {
		return nil, err
	}
	return &Contract{
		Upstream: upstream,
		Path:     local,
		SHA:      sha,
		Exists:   true,
	}, nil
}

// gitSHA returns SHA-1 hash of a file from the local Git storage.
// This SHA is comparable to SHAs returned in Github file metadata.
func gitSHA(ctx context.Context, file string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "hash-object", file)
	stdout, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git hash-object: %w", err)
	}
	return strings.TrimSpace(string(stdout)), nil
}

// ghFiles fetches metadata for files in a given weaviate/weaviate project directory.
func ghFiles(ctx context.Context, dir string) ([]GithubFile, error) {
	rc, err := ghGet(ctx, WeaviateRoot+"/"+dir)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	body, err := io.ReadAll(rc)
	if err != nil {
		return nil, err
	}

	var files []GithubFile
	if err := json.Unmarshal(body, &files); err != nil {
		return nil, err
	}
	return files, nil
}

// ghGet sends a GET request with User-Agent headers to the specified Github URL.
func ghGet(ctx context.Context, url string) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("gh get: %w", err)
	}

	req.Header.Add("User-Agent", "weaviate-go-client")
	if tok, ok := os.LookupEnv(EnvGithubToken); ok {
		req.Header.Add("Authorization", "Bearer "+tok)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("gh get: %w", err)
	}
	return res.Body, nil
}
