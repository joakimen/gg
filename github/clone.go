package github

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"sync"

	"github.com/joakimen/gg"
	"github.com/joakimen/gg/fuzzy"
	"github.com/joakimen/gg/git"
	"github.com/joakimen/gg/keyring"
)

func Clone(ctx context.Context,
	outDir string,
	owner string,
	repo string,
	defaultGitHubUser string,
	shallow bool,
	repoFile string,
) error {
	slog.DebugContext(ctx, "cloning repositories")
	if outDir == "" {
		return fmt.Errorf(
			"must specify clone directory through --clone-dir or by setting the $GG_CLONE_DIR environment variable",
		)
	}
	keyringManager := keyring.NewManager(KeyringUser)
	token, err := keyringManager.Get()
	if err != nil {
		return fmt.Errorf("failed to fetch token from keyring: %w", err)
	}

	githubClient := NewClient(token)
	repos, err := getRepos(ctx,
		githubClient,
		owner,
		repo,
		repoFile,
		defaultGitHubUser,
	)
	if err != nil {
		return fmt.Errorf("failed to fetch repos: %w", err)
	}

	if len(repos) == 0 {
		return fmt.Errorf("no repos to clone")
	}

	fmt.Printf("cloning %d repos(s)..\n", len(repos))
	cloneErrors := cloneAll(repos, outDir, shallow)
	if len(cloneErrors) > 0 {
		fmt.Fprintln(os.Stderr, "failed to clone some repos:")
		for _, e := range cloneErrors {
			fmt.Fprintf(os.Stderr, "- %s/%s: %v\n", e.Repo.Owner, e.Repo.Name, e.Err)
		}
	} else {
		fmt.Println("all repos cloned successfully!")
	}
	return nil
}

func getRepos(
	ctx context.Context,
	clent Client,
	owner string,
	repo string,
	repoFile string,
	defaultGitHubUser string,
) ([]gg.Repo, error) {
	fuzzyProvider := fuzzy.NewProvider()

	var repos []gg.Repo
	switch {
	case repoFile != "":
		repoJSONData, err := os.ReadFile(repoFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read repo file: %w", err)
		}
		if err = json.Unmarshal(repoJSONData, &repos); err != nil {
			return nil, fmt.Errorf("failed to unmarshal repos from file: %w", err)
		}
		return repos, nil
	case owner != "" && repo != "":
		repos = []gg.Repo{
			{
				Owner: owner,
				Name:  repo,
			},
		}
		return repos, nil
	case owner != "":
		slog.DebugContext(ctx, "listing repos by owner", "owner", owner)
		repos, err := clent.ListRepositoriesByUser(ctx, owner)
		if err != nil {
			return nil, fmt.Errorf("failed to list repos by owner: %w", err)
		}
		return fuzzyProvider.Select(repos)
	case repo != "":
		slog.DebugContext(ctx, "searching repos by name", "name", repo)
		repos, err := clent.SearchRepositoriesByName(ctx, repo)
		if err != nil {
			return nil, fmt.Errorf("failed to search repos by name: %w", err)
		}
		return fuzzyProvider.Select(repos)
	default:
		slog.DebugContext(ctx, "listing repos by default github user", "user", defaultGitHubUser)
		if defaultGitHubUser == "" {
			return nil, fmt.Errorf("no default GitHub user configured")
		}
		repos, err := clent.ListRepositoriesByUser(ctx, defaultGitHubUser)
		if err != nil {
			return nil, fmt.Errorf("failed to list repos by authenticated user: %w", err)
		}
		return fuzzyProvider.Select(repos)
	}
}

func cloneAll(repos []gg.Repo, outDir string, shallow bool) []gg.CloneResult {
	slog.Debug("cloning repos", "outDir", outDir, "repos", repos)

	var wg sync.WaitGroup
	gitService := git.NewClient()
	resultChan := make(chan gg.CloneResult, len(repos))
	for _, repo := range repos {
		wg.Add(1)
		go func(r gg.Repo) {
			defer wg.Done()
			cloneError := gitService.Clone(repo, outDir, shallow)
			resultChan <- gg.CloneResult{Repo: r, Err: cloneError}
		}(repo)
	}
	wg.Wait()
	close(resultChan)
	var errs []gg.CloneResult
	for res := range resultChan {
		if res.Err != nil {
			errs = append(errs, res)
		}
	}
	return errs
}
