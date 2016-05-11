// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 110.
//!+

// Package github provides a Go API for the GitHub issue tracker.
// See https://developer.github.com/v3/search/#search-issues.
package github

import (
    "net/http"
    "fmt"
    "encoding/json"
    "bytes"
)

const RepoUrl = "https://api.github.com/repos/"

type Issue struct {
    Title     string
    Body      string // in Markdown format
    Assignee  string
    State     string
    milestone int
    labels    []string
    Number    int
}

func ListIssues(repo string) ([]Issue, error) {
    resp, err := http.Get(RepoUrl + repo + "/issues")
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("search query failed: %s", resp.Status)
    }

    var result []Issue
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }
    return result, nil
}

func GetIssue(repo string, issueId string) (*Issue, error) {
    resp, err := http.Get(RepoUrl + repo + "/issues/" + issueId)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("search query failed: %s", resp.Status)
    }

    var result Issue
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }
    return &result, nil
}

func UpdateIssue(repo string, issueId string, issue *Issue) (*Issue, error) {

    b := new(bytes.Buffer)
    if err := json.NewEncoder(b).Encode(issue); err != nil {
        return nil, err
    }

    println(RepoUrl + repo + "/issues/" + issueId)
    req, err := http.NewRequest("PATCH", RepoUrl + repo + "/issues/" + issueId, b)
    if err != nil {
        return nil, err
    }
    req.Header.Set("Accept", "application/vnd.github.v3+json")

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("search query failed: %s", resp.Status)
    }

    var result Issue
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }
    return &result, nil
}

