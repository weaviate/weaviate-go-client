//nolint:errcheck
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/oasdiff/yaml"
)

const (
	EnvGithubToken = "GITHUB_TOKEN"

	WeaviateRoot        = "https://api.github.com/repos/weaviate/weaviate/contents"
	WeaviateProtobufDir = "grpc/proto/v1"
	WeaviateOpenAPIDir  = "openapi-specs/schema.json"

	LocalOpenAPIDir  = "./api/rest"
	LocalProtobufDir = "./api/proto/v1"

	SchemaCheck         = "schema.check.json"
	SchemaV3            = "schema.v3.yaml"
	SwaggerConverterURL = "https://converter.swagger.io/api/convert"

	ContractsHelp = `Usage: go run ./cmd/build contracts [--check]

Options:
	--check	Exit with non-zero exit code if checksums do not match.
	--help	Print this message.

Environment variables:
	GITHUB_TOKEN	GitHub token for authenticated requests.
`
)

// Contracts validates API contracts in ./api against their upstream
// versions in github.com/weaviate/weaviate.
//
// It works by fetching hashes for these files from the main repo
// and comparing them to hashes in the local .git object storage.
// If these differ, Contracts downloads the latest version of the
// contract and writes it to ./api.
//
// If --check flag is set to true, Contracts will return an error
// if the checksum fails. Useful in the context of a CI pipeline.
//
// Weaviate server defines its REST API schema in OpenAPI v2, which
// oapi-codegen does not support. After an update, Contracts will
// generate a v3 schema from the one we pull from the main repo.
//
// Models and protobuf stubs are not re-generated automatically.
func Contracts(ctx context.Context, args []string) error {
	var opt struct {
		Help  bool
		Check bool
	}
	cmd := flag.NewFlagSet("contracts", flag.ExitOnError)
	cmd.BoolVar(&opt.Help, "help", false, "Print help message.")
	cmd.BoolVar(&opt.Check, "check", false, "Exit with non-zero exit code if checksums do not match.")

	if err := cmd.Parse(args); err != nil {
		return err
	}

	if opt.Help {
		fmt.Print(ContractsHelp)
		return nil
	}

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

	var contracts []contract
	{
		c, err := newContract(ctx, LocalOpenAPIDir, SchemaCheck, *openapi)
		if err != nil {
			return err
		}
		contracts = append(contracts, *c)
	}

	for _, file := range protobufs {
		c, err := newContract(ctx, LocalProtobufDir, file.Name, file)
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
		if opt.Check {
			ok = false
			continue
		}

		log.Print("Downloading latest ", file.Upstream.DownloadURL)
		if err := updateContract(ctx, file); err != nil {
			printError(err)
			ok = false
		}
	}

	// Regenerate schema.v3.json if files were updated successfully.
	if !opt.Check && ok {
		log.Printf("Converting %s to OpenAPI v3", SchemaCheck)
		if err := convertSchemaToV3(ctx); err != nil {
			printError(err)
		}
	}

	// Prune api/proto and api/rest directories from stale files.
	log.Printf("Pruning local contracts")
	if err := pruneDir(LocalOpenAPIDir, contracts, opt.Check); err != nil {
		printError(err)
		ok = false
	}

	if err := pruneDir(LocalProtobufDir, contracts, opt.Check); err != nil {
		printError(err)
		ok = false
	}

	if opt.Check && !ok {
		return fmt.Errorf(`
Contracts in weaviate-go-client are out-of-sync with weaviate/weaviate repository.
Update them to the latest version by running this command:
	go run ./cmd/build contracts
			`) // nolint:staticcheck
	}
	if updated {
		log.Print(`
Contracts were successfully updated, run:
	go generate ./...

to re-generate REST and gRPC stubs.
`)
	}
	log.Print("Done")
	return nil
}

type contract struct {
	Upstream githubFile // Upstream file metadata.
	Exists   bool       // False if the file does not exist locally.
	Path     string     // Local filepath.
	SHA      string     // Local file SHA.
}

type githubFile struct {
	Name        string `json:"name"`
	SHA         string `json:"sha"`
	DownloadURL string `json:"download_url"`
}

// headProtobufs fetches metadata for schema.json.
func headOpenAPISpecs(ctx context.Context) (*githubFile, error) {
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
func headProtobufs(ctx context.Context) ([]githubFile, error) {
	return ghFiles(ctx, WeaviateProtobufDir)
}

// updateContract fetches the latest version of the [c.Upstream] and writes it to [c.Path].
// If w is not nil, it will receive the file's contents via an io.TeeReader.
func updateContract(ctx context.Context, c contract) error {
	rc, err := ghGet(ctx, c.Upstream.DownloadURL)
	if err != nil {
		return err
	}
	defer rc.Close()

	if _, err := writeAtomic(c.Path, rc, false); err != nil {
		return err
	}
	return nil
}

func newContract(ctx context.Context, dir, local string, upstream githubFile) (*contract, error) {
	path := filepath.Join(dir, local)
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return &contract{
			Upstream: upstream,
			Path:     path,
			SHA:      "<file not found>",
		}, nil
	}
	// os.Stat might've failed for a different reason, still try to get the hash.
	sha, err := gitSHA(ctx, path)
	if err != nil {
		return nil, err
	}
	return &contract{
		Upstream: upstream,
		Path:     path,
		SHA:      sha,
		Exists:   true,
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
func ghFiles(ctx context.Context, dir string) ([]githubFile, error) {
	rc, err := ghGet(ctx, WeaviateRoot+"/"+dir)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	body, err := io.ReadAll(rc)
	if err != nil {
		return nil, err
	}

	var files []githubFile
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
		return nil, fmt.Errorf("HTTP %d: %s", res.StatusCode, body)
	}

	return res.Body, nil
}

const headerGenerated = `# File generated by ./cmd/build. DO NOT EDIT MANUALLY.
# Source: %s.

`

// convertToV3 uses Swagger's public OpenAPI converter service to convert v2 docs to v3 docs.
// Swagger's conversion is best-effort, and won't be able to fix some of the invalid syntax:
//   - "Vector" in "components" -> "parameters" -> "Vector" must defined "x-go-type": "interface{}",
//     otherwise it is generated as map[string]any.
//   - "indexFilterable" should not have "x-nullable": true, so it would get a 'omitempty' tag.
//   - "Id" -> "Object" should have type "*openapi_types.UUID"
//   - Values with "format": "int64" or "uint64" => "type": "integer" (not "number") and "format": "int64".
//   - Values with "format": "int" or "int32" => "type": "integer" (not "number") and "format": "int32".
//   - Values with "type": "number" => "format": "float", not "float32" or "float64".
func convertToV3(ctx context.Context, r io.Reader) (map[string]any, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, SwaggerConverterURL, r)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(res.Body)
	res.Body.Close()

	if res.StatusCode > 299 {
		body := make(map[string]any)
		_ = json.Unmarshal(data, &body)
		return nil, fmt.Errorf("HTTP %d: %s", res.StatusCode, body["message"])
	} else if err != nil {
		return nil, err
	}

	var schema map[string]any
	if err := json.Unmarshal(data, &schema); err != nil {
		return nil, err
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

	walk("", schema, func(k string, m map[string]any) {
		switch k {
		case "Vector":
			if _, ok := m["x-go-type"]; !ok {
				m["x-go-type"] = "interface{}"
			}
			return
		case "indexInverted":
			// x-nullable is Swagger v2, nullable is Swagger v3
			delete(m, "nullable")
			return
		case "id":
			if format, ok := m["format"]; ok && format == "uuid" {
				m["x-go-type"] = "*openapi_types.UUID"
			}
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

	return schema, nil
}

func convertSchemaToV3(ctx context.Context) error {
	_, err := os.Stat(filepath.Join(LocalOpenAPIDir, SchemaV3))
	notExists := errors.Is(err, os.ErrNotExist)
	if _, err := os.Stat(filepath.Join(LocalOpenAPIDir, SchemaV3)); err != nil && !notExists {
		return err
	}

	f, err := os.Open(filepath.Join(LocalOpenAPIDir, SchemaCheck))
	if err != nil {
		return err
	}
	defer f.Close()

	schema, err := convertToV3(ctx, f)
	if err != nil {
		return err
	}

	v3, err := yaml.Marshal(schema)
	if err != nil {
		return err
	}

	path := filepath.Join(LocalOpenAPIDir, SchemaV3)
	header := fmt.Sprintf(headerGenerated, path)
	v3 = append([]byte(header), v3...)

	if _, err := writeAtomic(path, bytes.NewReader(v3), false); err != nil {
		return err
	}

	return nil
}

// writeAtomic writes the contents of Reader r to a .tmp file,
// performs a rename, and cleans up the .tmp file.
func writeAtomic(file string, r io.Reader, silent bool) (int64, error) {
	// Check if the file existed before the write to log
	// a more informative message on exit.
	_, err := os.Stat(filepath.Join(LocalOpenAPIDir, SchemaV3))
	existed := errors.Is(err, os.ErrNotExist)

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

	if !silent {
		if existed {
			log.Printf("Updated %s [written %dB]", file, written)
		} else {
			log.Printf("Added new file at %s [written %dB]", file, written)
		}
	}
	return written, nil
}

func pruneDir(root string, contracts []contract, check bool) error {
	// keep returns a boolean indicating if a contract file should be kept.
	keep := func(de fs.DirEntry) bool {
		// Both api/rest and api/proto/v1 should only contain files.
		if de.IsDir() {
			return false
		}

		// schema.v3.yaml is generated by this script.
		if de.Name() == SchemaV3 {
			return true
		}

		for _, c := range contracts {
			if filepath.Base(c.Path) == de.Name() {
				return true
			}
		}
		return false
	}

	stale := false
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == root {
			return nil
		}
		if keep(entry) {
			return nil
		}
		if check {
			log.Printf("prune %s: stale", path)
			stale = true
			return nil
		}
		log.Printf("prune %s: file is stale, removing", path)
		return os.RemoveAll(path)
	})
	if err != nil {
		return err
	}
	if stale {
		return errors.New("stale contracts must be removed")
	}
	return nil
}

func printError(err error) {
	fmt.Printf("\tERROR: %s\n", err)
}
