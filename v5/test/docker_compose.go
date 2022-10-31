package test

import (
	"os"
	"os/exec"
)

// SetupWeaviate run docker compose up
func SetupWeaviate() error {
	app := "docker-compose"
	arguments := []string{
		"up",
		"-d",
	}
	return command(app, arguments)
}

// TearDownWeaviate run docker-compose down
func TearDownWeaviate() error {
	app := "docker-compose"
	arguments := []string{
		"down",
		"--remove-orphans",
	}
	return command(app, arguments)
}

// SetupWeaviateDeprecated run docker compose up
// for pre-v1.14 backwards compatibility tests
func SetupWeaviateDeprecated() error {
	app := "docker-compose"
	arguments := []string{
		"-f",
		"docker-compose-deprecated-api-test.yml",
		"up",
		"-d",
	}
	return command(app, arguments)
}

// TearDownWeaviateDeprecated run docker-compose down
// for pre-v1.14 backwards compatibility tests
func TearDownWeaviateDeprecated() error {
	app := "docker-compose"
	arguments := []string{
		"-f",
		"docker-compose-deprecated-api-test.yml",
		"down",
		"--remove-orphans",
	}
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
	err = cmd.Start()
	return err
}
