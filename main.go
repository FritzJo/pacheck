package main

import (
	"flag"
	m "pacheck/modules"
)

func main() {
	quietflag := flag.Bool("q", false, "quiet: Only prints the vulnerable package name and version")
	cacheflag := flag.Bool("c", false, "cache: use the cached json from the last request")
	updateflag := flag.Bool("u", false, "update: fetch the latest json, but don't scan packages")
	flag.Parse()

	if !m.PacmanInstalled() {
		panic("[ERROR] Pacman not installed or not available!\nPlease make sure that everything is correctly set up.")
	}

	if *updateflag {
		m.FetchJson()
		return
	}
	vulnerabilities := m.GetVulnerabilities(*quietflag, *cacheflag)
	packages := m.GetInstalledPackages()
	for _, info := range packages {
		isVuln, vulns := m.IsVulnerable(vulnerabilities, info)
		if isVuln {
			m.PrintVulnerablePackage(info, vulns, *quietflag)
		}
	}
}
