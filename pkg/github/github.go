package github

import (
	"encoding/json"
	"fmt"
	"github.com/cli/go-gh"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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

func ListRepos(owner string, repo string, includeArchived bool, limit int) ([]Repo, error) {

	ghSearchArgs, err := BuildGhSearchArgs(owner, repo, includeArchived, limit)

	if err != nil {
		return nil, err
	}

	fmt.Println("searchArgs:", strings.Join(append([]string{"gh"}, ghSearchArgs...)[:], " "))

	ghRepoResp, _, err := gh.Exec(ghSearchArgs...)
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

func Clone(cloneDir string, repo Repo) error {
	repoAbsPath := filepath.Join(cloneDir, repo.Owner, repo.Name)
	if _, err := os.Stat(repoAbsPath); !os.IsNotExist(err) {
		return fmt.Errorf("repo already exists")
	}

	_, _, err := gh.Exec("repo", "clone", repo.NameWithOwner(), repoAbsPath)
	if err != nil {
		return err
	}
	return nil
}
