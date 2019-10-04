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


type Issue struct {
    Name          string `json:"name"`
    Package       []string `json:"packages"`
    Severity     string    `json:"severity"`
}


func main() {
    issues := getIssues()
    //cmd := exec.Command("apt", "list", "--installed")
    cmd := exec.Command("pacman", "-Q")
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
            for _, issue := range issues {
                if strings.Contains(scanner.Text(), issue.Package[0]) {
                    fmt.Println(issue.Severity + ": " + scanner.Text())
                }
            }
        }
    }
}


func getIssues() []Issue {
    jsonFile, err := os.Open("v.json")
    defer jsonFile.Close()
    if err != nil {
        fmt.Println(err)
    }
    result := make([]Issue, 0)
    decoder := json.NewDecoder(jsonFile)
    error_ := decoder.Decode(&result)
    if error_ != nil {
        fmt.Println(error_)
    }

    return result
}
