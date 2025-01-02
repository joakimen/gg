package github

import (
	"encoding/json"
	"fmt"
	"github.com/joakimen/clone"
	"github.com/joakimen/clone/exec"
	"os"
	"path/filepath"
	"strconv"
)

func gh(args ...string) (string, error) {

	stdout, _, err := exec.Exec("gh", args...)
	if err != nil {
		return "", fmt.Errorf("failed to execute gh: %w", err)
	}

	return stdout.String(), nil
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
func ListRepos(repoSearchArgs []string) ([]clone.Repo, error) {

	repoJsonData, err := gh(repoSearchArgs...)
	if err != nil {
		return nil, err
	}

	var searchResults clone.RepoSearchResults
	if err = json.Unmarshal([]byte(repoJsonData), &searchResults); err != nil {
		return nil, err
	}

	var repos []clone.Repo
	for _, repoResp := range searchResults {
		repos = append(repos, clone.Repo{Owner: repoResp.Owner.Login, Name: repoResp.Name})
	}
	return repos, nil
}

func Clone(cloneDir string, repo clone.Repo) error {
	repoAbsPath := filepath.Join(cloneDir, repo.Owner, repo.Name)
	if _, err := os.Stat(repoAbsPath); !os.IsNotExist(err) {
		return fmt.Errorf("repo already exists")
	}

	_, err := gh("repo", "clone", repo.NameWithOwner(), repoAbsPath)
	return err
}
