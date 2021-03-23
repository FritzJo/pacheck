package modules

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
)

func PacmanInstalled() bool {
	cmd := exec.Command("pacman")
	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	err := cmd.Run()
	if err != nil {
		return false
	}
	return true
}

func Check(err error, message string) {
	if err != nil {
		log.Fatal(message)
	}
}
