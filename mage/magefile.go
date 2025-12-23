//go:build mage

package main

import (
	"bytes"
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

	WeaviateRoot        = "https://api.github.com/repos/weaviate/weaviate/contents"
	WeaviateProtobufDir = "grpc/proto/v1"
	WeaviateOpenAPIDir  = "openapi-specs/schema.json"

	LocalOpenAPIDir  = "./api/rest"
	LocalProtobufDir = "./api/proto/v1"

	SchemaCheck = "schema.check.json"
	SchemaV3    = "schema.v3.json"

	SwaggerConverterURL = "https://converter.swagger.io/api/convert"
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
func validate(ctx context.Context, update bool) error {
	log.Printf("Fetching metadata for %s", WeaviateOpenAPIDir)
	openapi, err := headOpenAPISpecs(ctx)
	if err != nil {
		return err
	}

	log.Printf("Fetching metadata for %s/*.proto", WeaviateProtobufDir)
	protobufs, err := headProtobufs(ctx)
	if err != nil {
		return err
	}

	var contracts []Contract
	{
		c, err := newContract(ctx, LocalOpenAPIDir, SchemaCheck, *openapi, convertOpenAPIToV3)
		if err != nil {
			return err
		}
		contracts = append(contracts, *c)
	}

	for _, file := range protobufs {
		c, err := newContract(ctx, LocalProtobufDir, file.Name, file, nil)
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

			// Clone contents of the updated file to buf to make them
			// available to the callback without re-opening the file.
			var buf io.ReadWriter = nil
			if file.Callback != nil {
				buf = new(bytes.Buffer)
			}
			if err := updateContract(ctx, file, buf); err != nil {
				log.Printf("\tERROR: %s", err)
				ok = false
			} else {
				updated = true
				if file.Callback != nil {
					if err := file.Callback(ctx, buf); err != nil {
						log.Printf("\tERROR: %s", err)
						ok = false
					}
				}
			}
		} else {
			ok = false
		}
	}

	// If schema.check.json was present but schema.v3.json wasn't,
	// the convertOpenAPIToV3 did not run as a callback and we need
	// to force it manually.
	if ok && update && !updated {
		if err := convertExistingToV3(ctx); err != nil {
			log.Printf("\tERROR: convert existing %s to v3: %s", SchemaCheck, err)
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
	Upstream GithubFile   // Upstream file metadata.
	Exists   bool         // False if the file does not exist locally.
	Path     string       // Local filepath.
	SHA      string       // Local file SHA.
	Callback CallbackFunc // Callback is called after the contract's been updated.
}

type CallbackFunc func(context.Context, io.Reader) error

type GithubFile struct {
	Name        string `json:"name"`
	SHA         string `json:"sha"`
	DownloadURL string `json:"download_url"`
}

// headProtobufs fetches metadata for schema.json.
func headOpenAPISpecs(ctx context.Context) (*GithubFile, error) {
	dir, basename := filepath.Split(WeaviateOpenAPIDir)
	specs, err := ghFiles(ctx, dir)
	if err != nil {
		return nil, err
	}
	for _, file := range specs {
		if file.Name == basename {
			return &file, nil
		}
	}
	return nil, fmt.Errorf("%s/%s not found", WeaviateRoot, WeaviateOpenAPIDir)
}

// headProtobufs fetches metadata for protobuf files.
func headProtobufs(ctx context.Context) ([]GithubFile, error) {
	return ghFiles(ctx, WeaviateProtobufDir)
}

// updateContract fetches the latest version of the [c.Upstream] and writes it to [c.Path].
// If w is not nil, it will receive the file's contents via an io.TeeReader.
func updateContract(ctx context.Context, c Contract, w io.Writer) error {
	rc, err := ghGet(ctx, c.Upstream.DownloadURL)
	if err != nil {
		return err
	}
	defer rc.Close()

	var r io.Reader = rc
	if w != nil {
		r = io.TeeReader(r, w)
	}

	written, err := writeAtomic(c.Path, r)
	if err != nil {
		return err
	}

	if c.Exists {
		log.Printf("Updated %s [written %dB]", c.Path, written)
	} else {
		log.Printf("Added new file at %s [written %dB]", c.Path, written)
	}
	return nil
}

func newContract(ctx context.Context, dir, local string, upstream GithubFile, cb CallbackFunc) (*Contract, error) {
	path := filepath.Join(dir, local)
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return &Contract{
			Upstream: upstream,
			Path:     path,
			SHA:      "<file not found>",
			Callback: cb,
		}, nil
	}
	// os.Stat might've failed for a different reason, still try to get the hash.
	sha, err := gitSHA(ctx, path)
	if err != nil {
		return nil, err
	}
	return &Contract{
		Upstream: upstream,
		Path:     path,
		SHA:      sha,
		Exists:   true,
		Callback: cb,
	}, nil
}

// gitSHA returns [SHA-1] hash of a file from the local Git storage.
// This SHA is comparable to SHAs returned in Github file metadata.
//
// [SHA-1]: https://git-scm.com/book/en/v2/Git-Internals-Git-Objects
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

	if res.StatusCode > 299 {
		body, _ := io.ReadAll(res.Body)
		res.Body.Close()
		return nil, fmt.Errorf("HTTP %d: %s", res.StatusCode, string(body))
	}

	return res.Body, nil
}

// convertToV3 uses Swagger's public OpenAPI converter service to convert v2 docs to v3 docs.
// Swagger's conversion is best-effort, and won't be able to fix some of the invalid syntax:
//   - "Vector" in "components" -> "parameters" -> "Vector" must defined "x-go-type": "interface{}",
//     otherwise it is generated as map[string]any.
//   - Values with "format": "int64" or "uint64" => "type": "integer" (not "number") and "format": "int64".
//   - Values with "format": "int" or "int32" => "type": "integer" (not "number") and "format": "int32".
//   - Values with "type": "number" => "format": "float", not "float32" or "float64".
func convertOpenAPIToV3(ctx context.Context, r io.Reader) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, SwaggerConverterURL, r)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	data, err := io.ReadAll(res.Body)
	res.Body.Close()

	if res.StatusCode > 299 {
		body := make(map[string]any)
		_ = json.Unmarshal(data, &body)
		return fmt.Errorf("HTTP %d: %s", res.StatusCode, body["message"])
	} else if err != nil {
		return err
	}

	var v2 map[string]any
	if err := json.Unmarshal(data, &v2); err != nil {
		return err
	}

	var walk func(k string, m map[string]any, f func(string, map[string]any))

	walk = func(kout string, m map[string]any, f func(string, map[string]any)) {
		if kout != "" { // do not call f on the root map
			f(kout, m)
		}
		for k, v := range m {
			if v, ok := v.(map[string]any); ok {
				walk(k, v, f)
			}
		}
	}

	walk("", v2, func(k string, m map[string]any) {
		if k == "Vector" {
			if _, ok := m["x-go-type"]; !ok {
				m["x-go-type"] = "interface{}"
			}
			return
		}

		if format, ok := m["format"]; ok {
			switch format {
			case "int64", "uint64":
				m["type"] = "integer"
				m["format"] = "int64"
			case "int32", "int":
				m["type"] = "integer"
				m["format"] = "int32"
			}
		}

		if t, ok := m["type"]; ok && t == "number" {
			m["format"] = "float"
		}
	})

	v3, err := json.MarshalIndent(v2, "", "  ")
	if err != nil {
		return err
	}

	path := filepath.Join(LocalOpenAPIDir, SchemaV3)
	written, err := writeAtomic(path, bytes.NewReader(v3))
	if err != nil {
		return err
	}

	log.Printf("Converted OpenAPI schema to v3, new file at %s [written %dB]", path, written)
	return nil
}

func convertExistingToV3(ctx context.Context) error {
	if _, err := os.Stat(filepath.Join(LocalOpenAPIDir, SchemaV3)); err == nil {
		return nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	f, err := os.Open(filepath.Join(LocalOpenAPIDir, SchemaCheck))
	if err != nil {
		return err
	}
	defer f.Close()

	if err := convertOpenAPIToV3(ctx, f); err != nil {
		return err
	}
	return nil
}

// writeAtomic writes the contents of Reader r to a .tmp file,
// performs a rename, and cleans up the .tmp file.
func writeAtomic(file string, r io.Reader) (int64, error) {
	if err := os.MkdirAll(filepath.Dir(file), 0o775); err != nil {
		return 0, err
	}

	f, err := os.Create(file + ".tmp")
	if err != nil {
		return 0, err
	}
	defer f.Close()
	defer func() { os.Remove(f.Name()) }()

	written, err := io.Copy(f, r)
	if err != nil {
		return 0, err
	}

	if err := os.Rename(f.Name(), file); err != nil {
		return 0, err
	}
	return written, nil
}
