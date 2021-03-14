package main

import (
	"flag"
	"pacheck/modules"
)

func main() {
	quietflag := flag.Bool("q", false, "quiet: Only prints the vulnerable package name and version")
	cacheflag := flag.Bool("c", false, "cache: use the cached json from the last request")
	updateflag := flag.Bool("u", false, "update: fetch the latest json, but don't scan packages")
	flag.Parse()

	if !modules.PacmanInstalled() {
		panic("[ERROR] Pacman not installed or not available!\nPlease make sure that everything is correctly set up.")
	}

	if *updateflag {
		modules.FetchJson()
		return
	}
	vulnerabilities := modules.GetVulnerabilities(*quietflag, *cacheflag)
	packages := modules.GetInstalledPackages()
	for _, info := range packages {
		modules.IsVulnerable(vulnerabilities, info, *quietflag)
	}
}
