package modules

import (
	"bufio"
	"bytes"
	"os/exec"
	"strings"
)

type Packageinfo struct {
	Name    string
	Version string
}

func GetInstalledPackages() []Packageinfo {
	cmd := exec.Command("pacman", "-Q")
	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	err := cmd.Run()
	Check(err, "[ERROR] Can't find pacman executable!")

	info := []Packageinfo{}
	if len(cmdOutput.Bytes()) > 0 {
		x := string(cmdOutput.Bytes())

		scanner := bufio.NewScanner(strings.NewReader(x))
		for scanner.Scan() {
			text := scanner.Text()
			packagename := strings.Split(text, " ")[0]
			packageversion := strings.Split(text, " ")[1]
			info = append(info, Packageinfo{packagename, packageversion})
		}
	}
	return info
}
