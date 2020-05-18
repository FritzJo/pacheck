# Pacheck
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/0b059fc4954b406ea5c9543a73ecb234)](https://www.codacy.com/manual/fritzjo-git/pacheck?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=FritzJo/pacheck&amp;utm_campaign=Badge_Grade)

## Description
The name of this tool is a combination of the words _pacman_ and _check_ :-)

This tool checks installed Arch packages for known vulnerabilities. The data is collected from the amazing [Arch security dashboard](https://security.archlinux.org/) and matched against all currently installed packages.

My goal is to provide an alternative to the existing tool [arch-audit](https://github.com/ilpianista/arch-audit)

## How-To
### Build
```bash
git clone https://github.com/FritzJo/pacheck.git
cd pacheck
go build -o pacheck main.go
./pacheck
```
### Commandline options
|Parameter|Description|
|---|---|
|-q| quiet: Only prints the name and version of vulnerable packages|
|-c| cache: Use the last cached json (required if you want to use this tool offline)|
|-u| uodate: Fetch the latest json without scanning any packages|

## Example output
```bash
> ./pacheck
High: inetutils 1.9.4-7 CVE-2019-0053
Low: libmp4v2 2.0.0-5 CVE-2018-14054
Medium: libtiff 4.0.10-1 CVE-2019-7663 CVE-2019-6128
Low: openjpeg2 2.3.1-1 CVE-2019-6988
High: pacman 5.1.3-1 CVE-2019-18183 CVE-2019-18182
Low: unzip 6.0-13 CVE-2018-1000035
```
