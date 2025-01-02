package main

import (
	"fmt"
	"github.com/joakimen/clone"
	"github.com/joakimen/clone/config"
	"github.com/joakimen/clone/exec"
	"github.com/joakimen/clone/filter"
	"github.com/joakimen/clone/github"
	"github.com/joakimen/clone/log"
	"io"
	"os"
)

func main() {
	if err := run(os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

}

func run(w io.Writer) error {

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	logger := log.ConfigureLogger(w, cfg)
	logger.Info("done merging config", "config", cfg)

	if err := exec.DirExists(cfg.CloneDir); err != nil {
		return fmt.Errorf("clone directory does not exist: %w", err)
	}

	var repos []clone.Repo
	if cfg.Owner != "" && cfg.Repo != "" {
		logger.Debug("both owner and repo was specified, not searching for repos")
		repos = []clone.Repo{
			{
				Owner: cfg.Owner,
				Name:  cfg.Repo,
			},
		}
	} else {
		logger.Info("searching for repos", "owner", cfg.Owner, "repo", cfg.Repo, "includeArchived", cfg.IncludeArchived, "limit", cfg.Limit)
		repoSearchArgs, err := github.BuildGhSearchArgs(cfg.Owner, cfg.Repo, cfg.IncludeArchived, cfg.Limit)
		if err != nil {
			return fmt.Errorf("failed to build gh search args: %w", err)
		}

		logger.Info("repo search args", "args", append([]string{"gh"}, repoSearchArgs...))
		repos, err = github.ListRepos(repoSearchArgs)
		if err != nil {
			return fmt.Errorf("failed to list repos: %w", err)
		}
		logger.Info("repo search complete", "count", (len(repos)))
		logger.Debug("search returned repos", "repos", repos)
	}

	var selectedRepos []clone.Repo
	if len(repos) == 1 {
		selectedRepos = repos
	} else {
		logger.Info("filtering repos")
		var err error
		selectedRepos, err = filter.Select(repos)
		if err != nil {
			return fmt.Errorf("failed to select repos: %w", err)
		}
	}

	fmt.Fprintf(w, "cloning repos to %s:\n", cfg.CloneDir)
	for _, repo := range selectedRepos {
		fmt.Fprintf(w, "- %s\n", repo.NameWithOwner())
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

	var errorResults []clone.RepoCloneResult
	for i := 0; i < len(selectedRepos); i++ {
		result := <-cloneResultChan
		if result.Err != nil {
			errorResults = append(errorResults, result)
		}
	}
	close(cloneResultChan)

	fmt.Fprintf(w, "\n")
	if len(errorResults) > 0 {
		fmt.Fprintf(w, "some repos failed to clone:\n")
		for _, result := range errorResults {
			fmt.Fprintf(w, "- %s: %v\n", result.Repo.NameWithOwner(), result.Err)
		}
	} else {
		fmt.Fprintf(w, "all repos cloned successfully!\n")
	}
	return nil
}
