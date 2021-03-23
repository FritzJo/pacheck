package modules

import (
	"log"
	"os/exec"
)

func PacmanInstalled() bool {
	cmd := exec.Command("pacman", "-h")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

func Check(err error, message string) {
	if err != nil {
		log.Fatal(message)
	}
}
