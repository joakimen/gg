package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"sync"

	"github.com/joakimen/clone/github"
)

type CloneCommand struct{}

func (c *CloneCommand) Run(cfg Config) error {

	reposToClone, err := getReposToClone(cfg)
	if err != nil {
		return fmt.Errorf("failed to get repos to clone: %w", err)
	}

	if len(reposToClone) == 0 {
		fmt.Println("no repos selected, exiting")
		return nil
	}

	cloneErrors := clone(cfg, reposToClone)
	if len(cloneErrors) > 0 {
		fmt.Fprintln(os.Stderr, "failed to clone some repos:")
		for _, e := range cloneErrors {
			fmt.Fprintf(os.Stderr, "- %s: %v\n", e.Repo.NameWithOwner(), e.Err)
		}
	} else {
		fmt.Println("all repos cloned successfully!")
	}
	return nil
}

func getReposToClone(cfg Config) ([]github.Repo, error) {
	var (
		reposToClone []github.Repo
		err          error
	)
	switch {
	case cfg.RepoFile != "":
		slog.Debug("reading repos from file", "file", cfg.RepoFile)
		reposToClone, err = readReposFromFile(cfg.RepoFile)
		if err != nil {
			return []github.Repo{}, fmt.Errorf("couldn't read repos from file: %w", err)
		}
	case cfg.Owner != "" && cfg.Repo != "":
		slog.Debug("both owner and repo were provided, cloning single repo", "owner", cfg.Owner, "repo", cfg.Repo)
		reposToClone = []github.Repo{
			{
				Owner: cfg.Owner,
				Name:  cfg.Repo,
			},
		}
	default:
		slog.Debug("searching github for repos", "owner", cfg.Owner, "repo", cfg.Repo)
		reposToClone, err = github.Search(cfg.Owner, cfg.Repo, cfg.IncludeArchived, cfg.Limit)
		if err != nil {
			return []github.Repo{}, fmt.Errorf("failed to search for repos: %w", err)
		}
	}
	return reposToClone, nil
}

func clone(cfg Config, reposToClone []github.Repo) []github.CloneResult {
	slog.Debug("cloning repos", "cloneDir", cfg.CloneDir, "repos", reposToClone)

	var wg sync.WaitGroup
	resultChan := make(chan github.CloneResult, len(reposToClone))
	for _, repo := range reposToClone {
		wg.Add(1)
		go func(r github.Repo) {
			defer wg.Done()
			cloneError := github.Clone(cfg.CloneDir, repo)
			resultChan <- github.CloneResult{Repo: r, Err: cloneError}
		}(repo)
	}
	wg.Wait()
	close(resultChan)

	slog.Debug("procesing clone results")
	var errs []github.CloneResult
	for res := range resultChan {
		if res.Err != nil {
			errs = append(errs, res)
		}
	}
	return errs
}

func readReposFromFile(filepath string) ([]github.Repo, error) {
	var reposFromFile []github.Repo
	repoJSONData, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read repo file: %w", err)
	}

	if err = json.Unmarshal(repoJSONData, &reposFromFile); err != nil {
		return nil, fmt.Errorf("failed to unmarshal repos from file: %w", err)
	}
	return reposFromFile, nil
}
