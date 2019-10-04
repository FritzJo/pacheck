package main

import (
    "bytes"
    "fmt"
    "os"
    "os/exec"
    "strings"
    "encoding/json"
    "bufio"
)


type Vulnerability struct {
    Name        string   `json:"name"`
    Packages    []string `json:"packages"`
    Severity    string   `json:"severity"`
    Type        string   `json:"type"`
    Affected    string   `json:"affected"`
    Fixed       string   `json:"fixed"`
    Ticket      string   `json:"ticket"`
    Issues      []string `json:"issues"`
    Advisories  []string `json:"advisories"`
}


func main() {
    vulnerabilities := getVulnerabilities()
    cmd := exec.Command("apt", "list", "--installed")
    //cmd := exec.Command("pacman", "-Q")
    cmdOutput := &bytes.Buffer{}
    cmd.Stdout = cmdOutput
    err := cmd.Run()

    if err != nil {
        fmt.Println(err)
    }

    if len(cmdOutput.Bytes())>0 {
        x := string(cmdOutput.Bytes())
        scanner := bufio.NewScanner(strings.NewReader(x))
        for scanner.Scan() {
            for _, vuln := range vulnerabilities {
                for _, pack := range vuln.Packages {
                    if strings.Contains(scanner.Text(), pack) {
                        fmt.Println(vuln.Severity + ": " + scanner.Text())
                    }
                }
            }
        }
    }
}


func getVulnerabilities() []Vulnerability {
    jsonFile, err := os.Open("v.json")
    defer jsonFile.Close()
    if err != nil {
        fmt.Println(err)
    }
    result := make([]Vulnerability, 0)
    decoder := json.NewDecoder(jsonFile)
    error_ := decoder.Decode(&result)
    if error_ != nil {
        fmt.Println(error_)
    }

    return result
}
