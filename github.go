// Package github provides functions to interact with GitHub using the gh CLI tool.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

type Repo struct {
	Owner string
	Name  string
}

func (r Repo) NameWithOwner() string {
	return fmt.Sprintf("%s/%s", r.Owner, r.Name)
}

type SearchResp struct {
	Name  string `json:"name"`
	Owner struct {
		Login string `json:"login"`
	} `json:"owner"`
}

func gh(args ...string) (stdOut, stdErr bytes.Buffer, err error) {

	path, err := exec.LookPath("gh")
	if err != nil {
		err = fmt.Errorf("could not find gh executable in PATH. error: %w", err)
		return
	}

	cmd := exec.Command(path, args...)
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf("failed to run gh: %s. error: %w", stdErr.String(), err)
		return
	}

	return
}

// BuildGhSearchArgs builds an appropriate gh command to search for repos based on the provided parameters
func BuildGhSearchArgs(owner string, repo string, includeArchived bool, limit int) ([]string, error) {
	if owner != "" && repo != "" {
		return nil, fmt.Errorf("owner, repo or both must be empty to initiate a search")
	}

	limitStr := strconv.Itoa(limit)
	var args []string
	if owner != "" {
		if includeArchived {
			args = []string{"repo", "list", owner, "--json", "name,owner", "--limit", limitStr}
		} else {
			args = []string{"repo", "list", owner, "--json", "name,owner", "--no-archived", "--limit", limitStr}
		}
	} else if repo != "" {
		if includeArchived {
			args = []string{"search", "repos", repo, "--match", "name", "--json", "name,owner", "--limit", limitStr}
		} else {
			args = []string{"search", "repos", repo, "--match", "name", "--json", "name,owner", "--archived", "false", "--limit", limitStr}
		}
	} else {
		if includeArchived {
			args = []string{"repo", "list", "--json", "name,owner", "--limit", limitStr}
		} else {
			args = []string{"repo", "list", "--json", "name,owner", "--no-archived", "--limit", limitStr}
		}
	}
	return args, nil
}

// ListRepos searches for repos using gh based on the provided search arguments
func ListRepos(repoSearchArgs []string) ([]Repo, error) {

	ghRepoResp, _, err := gh(repoSearchArgs...)
	if err != nil {
		return nil, err
	}

	repoJsonData := ghRepoResp.String()

	var searchResp []SearchResp
	if err = json.Unmarshal([]byte(repoJsonData), &searchResp); err != nil {
		return nil, err
	}

	var repos []Repo
	for _, repoResp := range searchResp {
		repos = append(repos, Repo{Owner: repoResp.Owner.Login, Name: repoResp.Name})
	}
	return repos, nil
}

func clone(cloneDir string, repo Repo) error {
	repoAbsPath := filepath.Join(cloneDir, repo.Owner, repo.Name)
	if _, err := os.Stat(repoAbsPath); !os.IsNotExist(err) {
		return fmt.Errorf("repo already exists")
	}

	_, _, err := gh("repo", "clone", repo.NameWithOwner(), repoAbsPath)
	if err != nil {
		return err
	}
	return nil
}
