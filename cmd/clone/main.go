package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/joakimen/clone/github"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "failed to run:", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	info := func(args ...interface{}) {
		fmt.Fprintln(os.Stderr, args...)
	}

	debug := func(args ...interface{}) {
		if cfg.Verbose {
			fmt.Fprintln(os.Stderr, args...)
		}
	}

	debug("config loaded:", cfg)
	if err = DirExists(cfg.CloneDir); err != nil {
		return fmt.Errorf("clone directory does not exist: %w", err)
	}

	var repos []github.Repo
	switch {
	case cfg.RepoFile != "":
		debug("reading repos from file:", cfg.RepoFile)
		repos, err = readReposFromFile(cfg.RepoFile)
	case cfg.Owner != "" && cfg.Repo != "":
		debug("using repo specified by owner and repo flags")
		repos = []github.Repo{
			{
				Owner: cfg.Owner,
				Name:  cfg.Repo,
			},
		}
	default:
		debug("querying github for repos")
		repos, err = github.Search(cfg.Owner, cfg.Repo, cfg.IncludeArchived, cfg.Limit)
	}

	if len(repos) == 0 {
		info("no repos selected, exiting")
		return nil
	}

	info("cloning repos to:", cfg.CloneDir)
	for _, repo := range repos {
		info("-", repo.NameWithOwner())
	}

	cloneResultChan := make(chan github.RepoCloneResult, len(repos))

	for _, repo := range repos {
		go func(r github.Repo) {
			err = github.Clone(cfg.CloneDir, r)
			if err != nil {
				cloneResultChan <- github.RepoCloneResult{Repo: r, Err: err}
				return
			}
			cloneResultChan <- github.RepoCloneResult{Repo: r, Err: nil}
		}(repo)
	}

	var errorResults []github.RepoCloneResult
	for range repos {
		result := <-cloneResultChan
		if result.Err != nil {
			errorResults = append(errorResults, result)
		}
	}
	close(cloneResultChan)

	if len(errorResults) > 0 {
		for _, result := range errorResults {
			info(fmt.Sprintf("%s: %v", result.Repo.NameWithOwner(), result.Err))
		}
	} else {
		info("all repos cloned successfully!")
	}
	return nil
}
func readReposFromFile(filepath string) ([]github.Repo, error) {
	var repos []github.Repo
	repoJSONData, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read repo file: %w", err)
	}

	if err = json.Unmarshal(repoJSONData, &repos); err != nil {
		return nil, fmt.Errorf("failed to unmarshal repo file: %w", err)
	}
	return repos, nil
}
