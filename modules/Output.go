package modules

import (
	"fmt"
)

func PrintVulnerablePackage(packagei Packageinfo, vulnerabilities []Vulnerability, quiet bool) {
	for _, vuln := range vulnerabilities {
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
