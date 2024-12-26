package main

import (
	"fmt"
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

	cfg, err := LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	logger := ConfigureLogger(w, cfg)
	logger.Info("done merging config", "config", cfg)

	if err := dirExists(err, cfg); err != nil {
		return err
	}

	var repos []Repo
	if cfg.Owner != "" && cfg.Repo != "" {
		logger.Debug("both owner and repo was specified, not searching for repos")
		repos = []Repo{
			{
				Owner: cfg.Owner,
				Name:  cfg.Repo,
			},
		}
	} else {
		logger.Info("searching for repos", "owner", cfg.Owner, "repo", cfg.Repo, "includeArchived", cfg.IncludeArchived, "limit", cfg.Limit)
		repoSearchArgs, err := BuildGhSearchArgs(cfg.Owner, cfg.Repo, cfg.IncludeArchived, cfg.Limit)
		if err != nil {
			return err
		}

		logger.Info("repo search args", "args", append([]string{"gh"}, repoSearchArgs...))
		repos, err = ListRepos(repoSearchArgs)
		if err != nil {
			return err
		}
		logger.Info("repo search complete", "count", (len(repos)))
		logger.Debug("search returned repos", "repos", repos)
	}

	var selectedRepos []Repo
	if len(repos) == 1 {
		selectedRepos = repos
	} else {
		logger.Info("filtering repos")
		var err error
		selectedRepos, err = SelectRepos(repos)
		if err != nil {
			return err
		}
	}

	fmt.Fprintf(w, "cloning repos to %s:\n", cfg.CloneDir)
	for _, repo := range selectedRepos {
		fmt.Fprintf(w, "- %s\n", repo.NameWithOwner())
	}

	type CloneResult struct {
		repo Repo
		err  error
	}

	cloneResultChan := make(chan CloneResult, len(selectedRepos))
	for _, repo := range selectedRepos {
		go func(r Repo) {
			err := clone(cfg.CloneDir, r)
			if err != nil {
				cloneResultChan <- CloneResult{repo: r, err: err}
				return
			}
			cloneResultChan <- CloneResult{repo: r, err: nil}
		}(repo)
	}

	var errorResults []CloneResult
	for i := 0; i < len(selectedRepos); i++ {
		result := <-cloneResultChan
		if result.err != nil {
			errorResults = append(errorResults, result)
		}
	}
	close(cloneResultChan)

	fmt.Fprintf(w, "\n")
	if len(errorResults) > 0 {
		fmt.Fprintf(w, "some repos failed to clone:\n")
		for _, result := range errorResults {
			fmt.Fprintf(w, "- %s: %v\n", result.repo.NameWithOwner(), result.err)
		}
	} else {
		fmt.Fprintf(w, "all repos cloned successfully!\n")
	}
	return nil
}

func dirExists(err error, cfg Config) error {
	info, err := os.Stat(cfg.CloneDir)
	if os.IsNotExist(err) || !info.IsDir() {
		return fmt.Errorf("the specified clone directory does not exist: %s", cfg.CloneDir)
	}
	return nil
}
