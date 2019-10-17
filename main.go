package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
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
	cacheflag := flag.Bool("c", false, "cache: use the cached json from the last request")
	flag.Parse()

	vulnerabilities := getVulnerabilities(*quietflag, *cacheflag)
	cmd := exec.Command("pacman", "-Q")

	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	err := cmd.Run()

	if err != nil {
		log.Panic("[ERROR] Can't find pacman executable!")
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

func getVulnerabilities(quiet bool, cache bool) []vulnerability {
	if cache == true {
		if quiet == false {
			log.Print("[INFO] Using cached json!")
		}

		data, err := ioutil.ReadFile("./vulnerable.json")
		if err != nil {
			log.Panic("[ERROR] Can't find the cached json file!\n")
		}

		var cachedvuln []vulnerability
		err = json.Unmarshal(data, &cachedvuln)
		if err != nil {
			log.Panic("[ERROR] Error while reading cached json. The file may be damaged.")
		}

		return cachedvuln
	}

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

	result := make([]vulnerability, 0)
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&result)
	if err != nil {
		log.Panic("[ERROR] Can't decode the json response from " + url)
	}

	file, _ := json.MarshalIndent(result, "", " ")
	err = ioutil.WriteFile("vulnerable.json", file, 0644)

	if err != nil {
		log.Fatal("[ERROR] Can't save json response to disk. Cached version might be unavailable!")
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
					fmt.Print(vuln.Severity + ": " + packagei.Name + " " + packagei.Version + " ")
					for _, cve := range vuln.Issues {
						fmt.Print(cve + " ")
					}
					fmt.Println()
				}
			}
		}
	}
}
