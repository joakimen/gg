package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"sync"

	"github.com/joakimen/gg/github"
)

func clone(cloneDir string, reposToClone []github.Repo, shallow bool) []github.CloneResult {
	slog.Debug("cloning repos", "cloneDir", cloneDir, "repos", reposToClone)

	var wg sync.WaitGroup
	resultChan := make(chan github.CloneResult, len(reposToClone))
	for _, repo := range reposToClone {
		wg.Add(1)
		go func(r github.Repo) {
			defer wg.Done()
			cloneError := github.Clone(cloneDir, repo, shallow)
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
