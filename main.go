package main

import (
    "bytes"
    "fmt"
    "os"
    "os/exec"
    "strings"
    "encoding/json"
)

func main() {

    issues := getIssues()
    for _, issue := range issues {
        fmt.Println(issue.Severity + ": " +issue.Package[0])
    }

    cmd := exec.Command("pacman", "-Q")
    cmdOutput := &bytes.Buffer{}
    cmd.Stdout = cmdOutput
    printCommand(cmd)
    err := cmd.Run()
    printError(err)
    printOutput(cmdOutput.Bytes())

}


func printCommand(cmd *exec.Cmd) {
    fmt.Printf("==> Executing: %s\n", strings.Join(cmd.Args, " "))
}


func printError(err error) {
    if err != nil {
        os.Stderr.WriteString(fmt.Sprintf("==> Error: %s\n", err.Error()))
    }
}


func printOutput(outs []byte) {
    if len(outs) > 0 {
        fmt.Printf("==> Output: %s\n", string(outs))
    }
}

type Issue struct {
    Name          string `json:"name"`
    Package       []string `json:"packages"`
    Severity     string    `json:"severity"`
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
