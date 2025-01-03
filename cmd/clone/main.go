package main

import (
	"encoding/json"
	"fmt"
	"github.com/joakimen/clone"
	"github.com/joakimen/clone/config"
	"github.com/joakimen/clone/exec"
	"github.com/joakimen/clone/filter"
	"github.com/joakimen/clone/github"
	"io"
	"os"
)

func main() {
	writer := os.Stderr
	if err := run(writer); err != nil {
		fmt.Fprintf(writer, "%s\n", err)
		os.Exit(1)
	}

}

func run(writer io.Writer) error {

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %stdout", err)
	}

	info := func(args ...interface{}) {
		fmt.Fprintln(writer, args...)
	}

	debug := func(args ...interface{}) {
		if cfg.Verbose {
			fmt.Fprintln(writer, args...)
		}
	}

	debug("loaded config", cfg)
	if err := exec.DirExists(cfg.CloneDir); err != nil {
		return fmt.Errorf("clone directory does not exist: %stdout", err)
	}

	var repos clone.Repos

	if cfg.RepoFile != "" {
		debug("reading repos from file", "file", cfg.RepoFile)
		repoJsonData, err := os.ReadFile(cfg.RepoFile)
		if err != nil {
			return fmt.Errorf("failed to read repo file: %stdout", err)
		}

		if err = json.Unmarshal(repoJsonData, &repos); err != nil {
			return fmt.Errorf("failed to unmarshal repo file: %stdout", err)
		}
	} else if cfg.Owner != "" && cfg.Repo != "" {
		debug("both owner and repo was specified, not searching for repos")
		repos = clone.Repos{
			clone.Repo{
				Owner: cfg.Owner,
				Name:  cfg.Repo,
			},
		}
	} else {
		debug("searching for repos", "owner", cfg.Owner, "repo", cfg.Repo, "includeArchived", cfg.IncludeArchived, "limit", cfg.Limit)
		repoSearchArgs, err := github.BuildGhSearchArgs(cfg.Owner, cfg.Repo, cfg.IncludeArchived, cfg.Limit)
		if err != nil {
			return fmt.Errorf("failed to build gh search args: %stdout", err)
		}

		debug("repo search args", "args", append([]string{"gh"}, repoSearchArgs...))
		repos, err = github.ListRepos(repoSearchArgs)
		if err != nil {
			return fmt.Errorf("failed to list repos: %stdout", err)
		}
		debug("repo search complete", "count", len(repos))
		debug("search returned repos", "repos", repos)
	}

	var selectedRepos clone.Repos

	debug("repo count:", len(repos))
	switch len(repos) {
	case 0:
		debug("no repos found")
		return nil
	case 1:
		debug("cloning the specified owner/repo combination")
		selectedRepos = repos
	default:
		debug("filtering repos")
		selectedRepos, err = filter.Select(repos)
		if err != nil {
			return fmt.Errorf("failed to select repos: %stdout", err)
		}
	}

	info("cloning repos to:", cfg.CloneDir)
	for _, repo := range selectedRepos {
		info("-", repo.NameWithOwner())
	}

	cloneResultChan := make(chan clone.RepoCloneResult, len(selectedRepos))
	for _, repo := range selectedRepos {
		go func(r clone.Repo) {
			err := github.Clone(cfg.CloneDir, r)
			if err != nil {
				cloneResultChan <- clone.RepoCloneResult{Repo: r, Err: err}
				return
			}
			cloneResultChan <- clone.RepoCloneResult{Repo: r, Err: nil}
		}(repo)
	}

	var errorResults clone.RepoCloneResults
	for i := 0; i < len(selectedRepos); i++ {
		result := <-cloneResultChan
		if result.Err != nil {
			errorResults = append(errorResults, result)
		}
	}
	close(cloneResultChan)

	if len(errorResults) > 0 {
		for _, result := range errorResults {
			info(result.Repo.NameWithOwner()+":", result.Err)
		}
	} else {
		info("all repos cloned successfully!")
	}
	return nil
}
