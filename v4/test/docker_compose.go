package test

import (
	"os"
	"os/exec"
)

// SetupWeavaite run docker compose up
func SetupWeavaite() error {
	app := "docker-compose"
	arguments := []string{
		"up",
		"-d",
	}
	return command(app, arguments)
}

// TearDownWeavaite run docker-compose down
func TearDownWeavaite() error {
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
