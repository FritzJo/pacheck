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
	updateflag := flag.Bool("u", false, "update: fetch the latest json, but don't scan packages")
	flag.Parse()

	if !pacmanInstalled() {
		panic("[ERROR] Pacman not installed or not available!\nPlease make sure that everything is correctly set up.")
	}

	if *updateflag {
		fetchJson()
		return
	}
	vulnerabilities := getVulnerabilities(*quietflag, *cacheflag)
	packages := getInstalledPackages()
	for _, info := range packages {
		isVulnerable(vulnerabilities, info, *quietflag)
	}
}

func getInstalledPackages() []packageinfo {
	cmd := exec.Command("pacman", "-Q")
	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	err := cmd.Run()
	check(err, "[ERROR] Can't find pacman executable!")

	info := []packageinfo{}
	if len(cmdOutput.Bytes()) > 0 {
		x := string(cmdOutput.Bytes())

		scanner := bufio.NewScanner(strings.NewReader(x))
		for scanner.Scan() {
			text := scanner.Text()
			packagename := strings.Split(text, " ")[0]
			packageversion := strings.Split(text, " ")[1]
			info = append(info, packageinfo{packagename, packageversion})
		}
	}
	return info
}

func getVulnerabilities(quiet bool, cache bool) []vulnerability {
	// Determine datasource (local cache, or web)
	if !cache {
		return fetchJson()
	}
	if !quiet {
		log.Print("[INFO] Using cached json!")
	}

	data, err := ioutil.ReadFile("./vulnerable.json")
	check(err, "[ERROR] Can't find the cached json file!\n")

	var cachedvuln []vulnerability
	err = json.Unmarshal(data, &cachedvuln)
	check(err, "[ERROR] Error while reading cached json. The file may be damaged.")

	return cachedvuln
}

func isVulnerable(vulnerabilities []vulnerability, packagei packageinfo, quiet bool) {
	// Iterate over all vulnerabilities of the loaded json
	for _, vuln := range vulnerabilities {
		for _, pack := range vuln.Packages {
			// A package is affected, when name and version match
			if !strings.Contains(packagei.Name, pack) || !strings.Contains(packagei.Version, vuln.Affected) {
				continue
			}
			if quiet == true {
				fmt.Println(packagei.Name + " " + vuln.Affected)
				continue
			}
			fmt.Print(vuln.Severity + ": " + packagei.Name + " " + packagei.Version + " ")
			for _, cve := range vuln.Issues {
				fmt.Print(cve + " ")
			}
			fmt.Println()
		}
	}
}

func fetchJson() []vulnerability {
	// Fetch newest json
	url := "https://security.archlinux.org/vulnerable.json"

	vulnClient := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	check(err, "[Error] Can't create new request.")

	res, err := vulnClient.Do(req)
	check(err, "[Error] Can't execute http request")

	result := make([]vulnerability, 0)
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&result)
	check(err, "[ERROR] Can't decode the json response from "+url)

	file, _ := json.MarshalIndent(result, "", " ")
	err = ioutil.WriteFile("vulnerable.json", file, 0644)
	check(err, "[ERROR] Can't write vulnerable.json file")

	return result
}

func pacmanInstalled() bool {
	fmt.Println("Test")
	cmd := exec.Command("pacman")
	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	err := cmd.Run()
	if err != nil {
		return false
	}
	return true
}

func check(err error, message string) {
	if err != nil {
		log.Fatal(message)
	}
}
