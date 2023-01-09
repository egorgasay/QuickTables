package createdb

import (
	"bytes"
	"errors"
	"os/exec"
)

func RunContainer(path string) error {
	cmd := exec.Command("docker-compose", "-f", path, "up", "-d")

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return err
	}

	output := out.String()
	if output != "" {
		return errors.New(output)
	}

	return nil
}
