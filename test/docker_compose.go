package test

import (
	"os"
	"os/exec"
)

// SetupWeaviate run docker compose up
func SetupWeaviate(authEnabled bool) error {
	app := "docker"
	arguments := []string{"compose"}

	if authEnabled {
		arguments = append(arguments, "-f", "docker-compose-wcs.yml")
	}
	arguments = append(arguments, "up", "-d")

	return command(app, arguments)
}

// TearDownWeaviate run docker-compose down
func TearDownWeaviate() error {
	app := "docker"
	arguments := []string{"compose", "down", "--remove-orphans"}
	return command(app, arguments)
}

func command(app string, arguments []string) error {
	mydir, err := os.Getwd()
	if err != nil {
		return err
	}

	cmd := exec.Command(app, arguments...)
	execDir := mydir + "/../"
	cmd.Dir = execDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	return err
}

type Container struct {
	DockerComposeFile string

	// APISecret must be configured via `AUTHENTICATION_APIKEY_ALLOWED_KEYS`.
	APISecret string

	host     string
	httpPort string // Should match the *exposed* port in docker-compose file.
	gRPCPort string // Specify to enable gRPC in weaviate.Client.
}

// HTTPAddress of the Weaviate container.
func (c Container) HTTPAddress() string {
	return c.host + ":" + c.httpPort
}

// GRPCAddress of the Weaviate container.
func (c Container) GRPCAddress() string {
	return c.host + ":" + c.gRPCPort
}

// EnableGRPS returns true if gRPC port is specified for this container.
func (c Container) EnableGRPC() bool {
	return c.gRPCPort != ""
}

// Preset is a group of Docker containers defined in one of the ../test/docker-compose*.yml files.
type Preset int

const (
	// Containers from ../test/docker-compose.yml
	Basic Preset = iota
	// Containers from ../test/docker-compose-rbac.yml
	RBAC
	Cluster
)

// TODO: add other presets from ../test/docker-compose*.yaml files
// and unify container management.
var containers map[Preset]Container = map[Preset]Container{
	Basic: {
		DockerComposeFile: "docker-compose.yml",
		host:              "localhost",
		httpPort:          "8080",
		gRPCPort:          "50051",
	},
	RBAC: {
		DockerComposeFile: "docker-compose-rbac.yml",
		APISecret:         "my-secret-key",
		host:              "localhost",
		httpPort:          "8089",
	},
	Cluster: {
		DockerComposeFile: "docker-compose-cluster.yml",
		host:              "localhost",
		httpPort:          "8087",
	},
}

// Get selected docker-compose configuration and function to start/stop the containers.
// The "down" command includes `--remove-orphans` flag, which will remove the containers
// in the current docker-compose file and *ANY* other containers (effectively: all).
func GetContainer(preset Preset) (container Container, up, down func() error) {
	container = containers[preset]
	compose := func(sub string, args ...string) error {
		return command("docker", append([]string{
			"compose", "-f", container.DockerComposeFile, sub,
		}, args...))
	}
	return container, func() error {
			return compose("up", "-d")
		}, func() error {
			return compose("down", "--remove-orphans")
		}
}
