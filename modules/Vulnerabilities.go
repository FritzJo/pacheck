package modules

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type Vulnerability struct {
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

type VulnerablePackage struct {
	Vuln     Vulnerability
	Packagei Packageinfo
}

func GetVulnerabilities(quiet bool, cache bool) []Vulnerability {
	// Determine datasource (local cache, or web)
	if !cache {
		return FetchJson()
	}
	if !quiet {
		log.Print("[INFO] Using cached json!")
	}

	data, err := ioutil.ReadFile("./vulnerable.json")
	Check(err, "[ERROR] Can't find the cached json file!\n")

	var cachedvuln []Vulnerability
	err = json.Unmarshal(data, &cachedvuln)
	Check(err, "[ERROR] Error while reading cached json. The file may be damaged.")

	return cachedvuln
}

func IsVulnerable(vulnerabilities []Vulnerability, packageInfo Packageinfo) (bool, []Vulnerability) {
	var affectedVulnerabilities []Vulnerability
	// Iterate over all vulnerabilities of the loaded json
	for _, vuln := range vulnerabilities {
		for _, pack := range vuln.Packages {
			// A package is affected, when name and version match
			if !strings.Contains(packageInfo.Name, pack) || !strings.Contains(packageInfo.Version, vuln.Affected) {
				continue
			}
			affectedVulnerabilities = append(affectedVulnerabilities, vuln)
		}
	}
	resultFound := len(affectedVulnerabilities) > 0
	return resultFound, affectedVulnerabilities
}

func FetchJson() []Vulnerability {
	// Fetch newest json
	url := "https://security.archlinux.org/vulnerable.json"

	vulnClient := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	Check(err, "[Error] Can't create new request.")

	res, err := vulnClient.Do(req)
	Check(err, "[Error] Can't execute http request")

	result := make([]Vulnerability, 0)
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&result)
	Check(err, "[ERROR] Can't decode the json response from "+url)

	file, _ := json.MarshalIndent(result, "", " ")
	err = ioutil.WriteFile("vulnerable.json", file, 0644)
	Check(err, "[ERROR] Can't write vulnerable.json file")

	return result
}
