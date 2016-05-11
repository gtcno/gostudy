package main

import (
    "os"
    "github.com/marcsauter/cli"
    "github.com/gtcno/gostudy/ch4/github"
    "log"
    "io/ioutil"
    "os/exec"
    "encoding/json"
)

const (
    Prompt = ">>> "
    Histsize = 255
    Editor = "vi"
)

var (
    Commands map[string]func(*cli.CLI, []string) error
)

func init() {
    Commands = map[string]func(*cli.CLI, []string) error{
        "help": help,
        "list": list,
        "read": read,
        "update": update,
        "exit": exit,
    }
}

func main() {
    c := cli.NewCLI(Commands, Prompt, Histsize)
    c.Info("Type help ...\n")
    for {
        c.Prompt()
        c.Exec(c.ReadLine())
    }
}

func update(c *cli.CLI, command []string) error {
    if len(command) < 3 {
        c.Info("Usage: update <user>/<repo> <issue#>\n")
        return nil
    }

    issue, err := github.GetIssue(command[1], command[2])
    if (err != nil ) {
        log.Fatal(err)
    }

    if (issue == nil) {
        c.Info("No issues found for issue: %s and id %s\n", command[1], command[1])
    } else {
        file, err := writeFile(issue)
        if (err != nil ) {
            log.Fatal(err)
        }
        path, err := exec.LookPath(Editor)
        if err != nil {
            log.Fatal(err)
        }

        cmd := exec.Command(path, file)
        cmd.Stdin = os.Stdin
        cmd.Stdout = os.Stdout
        err = cmd.Start()

        if err != nil {
            log.Fatal(err)
        }
        err = cmd.Wait()

        issue, err := readFile(file)
        if err != nil {
            log.Fatal(err)
        }

        newIssue, err := github.UpdateIssue(command[1], command[2],issue)
        if err != nil {
            log.Fatal(err)
        }
        json, _ := json.MarshalIndent(newIssue, "", " ")
        c.Info("%s\n", json)

    }
    return nil
}

func writeFile(issue *github.Issue) (string, error) {
    tmpDir := os.TempDir()
    tmpFile, tmpFileErr := ioutil.TempFile(tmpDir, "tmp")
    if tmpFileErr != nil {
        return "", tmpFileErr
    }
    defer tmpFile.Close()

    json, err := json.MarshalIndent(issue, "", " ")
    if err != nil {
        return "", err
    }

    _, err = tmpFile.Write(json)
    if err != nil {
        return "", err
    }

    return tmpFile.Name(), nil
}

func readFile(name string) (*github.Issue, error) {
    file, err := os.Open(name)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var issue github.Issue
    jsonParser := json.NewDecoder(file)
    if err = jsonParser.Decode(&issue); err != nil {
        return nil, err
    }

    return &issue, nil
}

func help(c *cli.CLI, command []string) error {
    helptext := `
help                        this text
list <user>/<repo>          list issues for the given repo
read <user>/<repo> <issue#> Read the given issue
exit                        exit program

`
    c.Info(helptext)
    return nil
}

func read(c *cli.CLI, command[]string) error {
    if len(command) < 3 {
        c.Info("Usage: read <user>/<repo> <issue#>\n")
        return nil
    }

    issue, err := github.GetIssue(command[1], command[2])
    if (err != nil ) {
        log.Fatal(err)
        return nil
    }

    if (issue == nil) {
        c.Info("No issues found for issue: %s and id %s\n", command[1], command[1])
    } else {
        c.Info("%s\n", issue.Body)
    }
    return nil
}

func list(c *cli.CLI, command []string) error {
    if len(command) < 2 {
        c.Info("Usage: list <user>/<repo>\n")
        return nil
    }
    issues, err := github.ListIssues(command[1])
    if (err != nil ) {
        log.Fatal(err)
        return nil
    }

    if (issues == nil) {
        c.Info("No issues found for repo: %s\n", command[1])
    } else {
        for _, item := range issues {
            c.Info("#%-5d %.55s \n",
                item.Number, item.Title)
        }
    }
    return nil
}

func exit(c *cli.CLI, command []string) error {
    if c.YesNo("Are you sure? (y)es/(n)o", []string{"y"}) {
        os.Exit(0)
    }
    return nil
}