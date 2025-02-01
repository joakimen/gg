package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/joakimen/clone/github"
)

type CloneCommand struct{}

func (c *CloneCommand) Run(cfg Config) error {
	var err error
	var repos []github.Repo
	switch {
	case cfg.RepoFile != "":
		slog.Debug("reading repos from file:", "repoFile", cfg.RepoFile)
		repos, err = readReposFromFile(cfg.RepoFile)
		if err != nil {
			return fmt.Errorf("failed to read repos from file: %w", err)
		}
	case cfg.Owner != "" && cfg.Repo != "":
		slog.Debug("using r specified by owner and r flags")
		repos = []github.Repo{
			{
				Owner: cfg.Owner,
				Name:  cfg.Repo,
			},
		}
	default:
		slog.Debug("querying github for repos")
		repos, err = github.Search(cfg.Owner, cfg.Repo, cfg.IncludeArchived, cfg.Limit)
		if err != nil {
			return fmt.Errorf("failed to search for repos: %w", err)
		}
	}

	if len(repos) == 0 {
		slog.Debug("no repos selected, exiting")
		return nil
	}

	slog.Debug("cloning repos", "cloneDir", cfg.CloneDir, "repos", repos)
	resultChan := make(chan github.RepoCloneResult, len(repos))
	for _, repo := range repos {
		go func(r github.Repo) {
			err = github.Clone(cfg.CloneDir, r)
			errClone := github.Clone(cfg.CloneDir, repo)
			resultChan <- github.RepoCloneResult{Repo: repo, Err: errClone}
		}(repo)
	}

	var errs []github.RepoCloneResult
	for range repos {
		res := <-resultChan
		if res.Err != nil {
			errs = append(errs, res)
		}
	}
	close(resultChan)

	if len(errs) > 0 {
		for _, e := range errs {
			slog.Debug(fmt.Sprintf("%s: %v", e.Repo.NameWithOwner(), e.Err))
		}
	} else {
		slog.Debug("all repos cloned successfully!")
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
