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
