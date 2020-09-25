package test

import (
	"os"
	"os/exec"
)

func SetupWeavaite() error {
	app := "docker-compose"
	arguments := []string{
		"up",
		"-d",
	}
	return command(app, arguments)
}

func TearDownWeavaite() error {
	app := "docker-compose"
	arguments := []string{
		"down",
	}
	return command(app, arguments)
}

func command(app string, arguments []string) error {
	mydir, err := os.Getwd()
	if err != nil {
		return err
	}

	cmd := exec.Command(app, arguments...)
	cmd.Dir = mydir + "/../test/"
	err = cmd.Start()
	return err
}