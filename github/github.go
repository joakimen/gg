package github

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

// buildGHCommand builds an appropriate gh command to find repos.
func buildGHCommand(owner string, repo string, includeArchived bool, limit int) ([]string, error) {
	if owner != "" && repo != "" {
		return nil, fmt.Errorf("owner, repo or both must be empty to initiate a search")
	}

	slog.Debug("parsing gh args", "owner", owner, "repo", repo, "includeArchived", includeArchived, "limit", limit)
	limitStr := strconv.Itoa(limit)
	var args []string
	switch {
	case owner != "":
		if includeArchived {
			args = []string{"repo", "list", owner, "--json", "name,owner", "--limit", limitStr}
		} else {
			args = []string{"repo", "list", owner, "--json", "name,owner", "--no-archived", "--limit", limitStr}
		}
	case repo != "":
		if includeArchived {
			args = []string{"search", "repos", repo, "--match", "name", "--json", "name,owner", "--limit", limitStr}
		} else {
			args = []string{
				"search", "repos", repo, "--match", "name", "--json", "name,owner", "--limit", limitStr,
				"--archived=false",
			}
		}
	default:
		if includeArchived {
			args = []string{"repo", "list", "--json", "name,owner", "--limit", limitStr}
		} else {
			args = []string{"repo", "list", "--json", "name,owner", "--no-archived", "--limit", limitStr}
		}
	}

	slog.Debug("done building gh args", "args", args)
	return args, nil
}

// listRepos searches for repos using Exec based on the provided search arguments.
func listRepos(repoSearchArgs []string) ([]Repo, error) {
	repoJSONData, err := gh(repoSearchArgs...)
	if err != nil {
		return nil, err
	}

	var searchResults []RepoSearchResult
	if err = json.Unmarshal([]byte(repoJSONData), &searchResults); err != nil {
		return nil, err
	}

	var repos []Repo
	for _, repoResp := range searchResults {
		repos = append(repos, Repo{Owner: repoResp.Owner.Login, Name: repoResp.Name})
	}
	return repos, nil
}

// Clone a single repo from GitHub to the specified cloneDir.
func Clone(cloneDir string, repo Repo, shallowClone bool) error {
	repoAbsPath := filepath.Join(cloneDir, repo.Owner, repo.Name)
	if _, err := os.Stat(repoAbsPath); !os.IsNotExist(err) {
		return fmt.Errorf("repo %s already exists in %s", repo.NameWithOwner(), repoAbsPath)
	}

	cloneArgs := []string{"repo", "clone", repo.NameWithOwner(), repoAbsPath}
	if shallowClone {
		cloneArgs = append(cloneArgs, "--depth", "1")
	}

	slog.Debug("cloning repo", "cloneArgs", cloneArgs)

	_, err := gh(cloneArgs...)
	if err != nil {
		return fmt.Errorf("failed to clone repo %s: %w", repo.NameWithOwner(), err)
	}
	return nil
}

func gh(args ...string) (string, error) {
	path, err := exec.LookPath("gh")
	if err != nil {
		return "", fmt.Errorf("could not find gh executable in PATH. error: %w", err)
	}

	stdout := bytes.Buffer{}
	cmd := exec.Command(path, args...)
	cmd.Stdout = &stdout

	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("gh failed: %w", err)
	}
	return stdout.String(), nil
}

func Search(owner string, repo string, includeArchived bool, limit int) ([]Repo, error) {
	var (
		repos []Repo
		err   error
	)

	repoSearchArgs, err := buildGHCommand(owner, repo, includeArchived, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to build gh search args: %w", err)
	}

	githubRepos, err := listRepos(repoSearchArgs)
	if err != nil {
		return nil, fmt.Errorf("failed to list github repos: %w", err)
	}

	if len(githubRepos) == 0 {
		return nil, errors.New("no github repos found with the provided search criteria")
	}

	repos, err = FuzzySelect(githubRepos)
	if err != nil {
		return nil, fmt.Errorf("failed to filter repos: %w", err)
	}
	return repos, nil
}
