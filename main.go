package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

type vulnerability struct {
	Name       string   `json:"name"`
	Packages   []string `json:"packages"`
	Severity   string   `json:"severity"`
	Type       string   `json:"type"`
	Affected   string   `json:"affected"`
	Fixed      string   `json:"fixed"`
	Ticket     string   `json:"ticket"`
	Issues     []string `json:"issues"`
	Advisories []string `json:"advisories"`
}

type packageinfo struct {
	Name    string
	Version string
}

func main() {
	quietflag := flag.Bool("q", false, "quiet: Only prints the vulnerable package name and version")
	flag.Parse()

	vulnerabilities := getVulnerabilities()
	cmd := exec.Command("pacman", "-Q")

	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	err := cmd.Run()

	if err != nil {
		fmt.Println(err)
	}

	if len(cmdOutput.Bytes()) > 0 {
		x := string(cmdOutput.Bytes())

		scanner := bufio.NewScanner(strings.NewReader(x))
		for scanner.Scan() {
			text := scanner.Text()
			packagename := strings.Split(text, " ")[0]
			packageversion := strings.Split(text, " ")[1]
			info := packageinfo{packagename, packageversion}
			isVulnerable(vulnerabilities, info, *quietflag)
		}
	}
}

func getVulnerabilities() []vulnerability {
	url := "https://security.archlinux.org/vulnerable.json"
	vulnClient := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, getErr := vulnClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if err != nil {
		fmt.Println(err)
	}

	result := make([]vulnerability, 0)
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&result)
	if err != nil {
		fmt.Println(err)
	}

	return result
}

func isVulnerable(vulnerabilities []vulnerability, packagei packageinfo, quiet bool) {
	for _, vuln := range vulnerabilities {
		for _, pack := range vuln.Packages {
			if strings.Contains(packagei.Name, pack) && strings.Contains(packagei.Version, vuln.Affected) {
				if quiet == true {
					fmt.Println(packagei.Name + " " + vuln.Affected)
				} else {
					fmt.Println(vuln.Severity + ": " + packagei.Name + " " + packagei.Version)
				}
			}
		}
	}
}
