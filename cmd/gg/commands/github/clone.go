package github

import (
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"sync"

	"github.com/joakimen/gg"
	"github.com/joakimen/gg/fuzzy"
	"github.com/joakimen/gg/git"
	"github.com/joakimen/gg/github"
	"github.com/joakimen/gg/keyring"
	"github.com/spf13/cobra"
)

type cloneFlags struct {
	Owner           string
	Repo            string
	RepoFile        string
	IncludeArchived bool
	OutDir          string
	Shallow         bool
}

type cloneEnvs struct {
	DefaultGitHubUser string
	OutDir            string
}

func loadCloneEnvs() cloneEnvs {
	return cloneEnvs{
		OutDir:            os.Getenv("GG_CLONE_DIR"),
		DefaultGitHubUser: os.Getenv("GG_GITHUB_USER"),
	}
}

func NewCloneCmd() *cobra.Command {
	var flags cloneFlags
	cloneCmd := &cobra.Command{
		Use:   "clone",
		Short: "Clone GitHub repos interactively",
		RunE: func(cmd *cobra.Command, _ []string) error {
			envs := loadCloneEnvs()
			slog.Debug("cloning repositories", "opts", flags, "envs", envs)
			outDir := cmp.Or(flags.OutDir, envs.OutDir)
			if outDir == "" {
				return fmt.Errorf(
					"must specify clone directory through --clone-dir or by setting the $GG_CLONE_DIR environment variable",
				)
			}
			keyringManager := keyring.NewManager(keyringUser)
			token, err := keyringManager.Get()
			if err != nil {
				return fmt.Errorf("failed to fetch token from keyring: %w", err)
			}

			githubService := github.NewService(token)
			repos, err := getRepos(cmd.Context(), githubService, flags, envs)
			if err != nil {
				return fmt.Errorf("failed to fetch repos: %w", err)
			}

			if len(repos) == 0 {
				return fmt.Errorf("no repos to clone")
			}

			fmt.Printf("cloning %d repos(s)..\n", len(repos))
			cloneErrors := cloneAll(repos, outDir, flags.Shallow)
			if len(cloneErrors) > 0 {
				fmt.Fprintln(os.Stderr, "failed to clone some repos:")
				for _, e := range cloneErrors {
					fmt.Fprintf(os.Stderr, "- %s: %v\n", e.Repo.NameWithOwner(), e.Err)
				}
			} else {
				fmt.Println("all repos cloned successfully!")
			}

			return nil
		},
	}

	cloneCmd.Flags().StringVarP(&flags.Owner, "owner", "o", "", "owner of the repo(s) to clone")
	cloneCmd.Flags().StringVarP(&flags.Repo, "repo", "r", "", "owner of the repo(s) to clone")
	cloneCmd.Flags().StringVarP(&flags.RepoFile, "file", "f", "", "name of file containing list of repos to clone")
	cloneCmd.Flags().StringVarP(&flags.OutDir, "out-dir", "d", "", "the output directory of cloned repos")
	cloneCmd.Flags().BoolVarP(&flags.IncludeArchived, "include-archived", "a", false, "owner of the repo(s) to clone")
	cloneCmd.Flags().BoolVarP(&flags.Shallow, "shallow", "", false, "perform a shallow clone")

	return cloneCmd
}

func getRepos(ctx context.Context, svc github.Service, flags cloneFlags, envs cloneEnvs) ([]gg.Repo, error) {
	fuzzyService := fuzzy.NewService()

	var repos []gg.Repo
	switch {
	case flags.RepoFile != "":
		repoJSONData, err := os.ReadFile(flags.RepoFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read repo file: %w", err)
		}
		if err = json.Unmarshal(repoJSONData, &repos); err != nil {
			return nil, fmt.Errorf("failed to unmarshal repos from file: %w", err)
		}
		return repos, nil
	case flags.Owner != "" && flags.Repo != "":
		repos = []gg.Repo{
			{
				Owner: flags.Owner,
				Name:  flags.Repo,
			},
		}
		return repos, nil
	case flags.Owner != "":
		slog.DebugContext(ctx, "listing repos by owner", "owner", flags.Owner)
		repos, err := svc.ListRepositoriesByUser(ctx, flags.Owner)
		if err != nil {
			return nil, fmt.Errorf("failed to list repos by owner: %w", err)
		}
		return fuzzyService.Select(repos)
	case flags.Repo != "":
		slog.DebugContext(ctx, "searching repos by name", "name", flags.Repo)
		repos, err := svc.SearchRepositoriesByName(ctx, flags.Repo)
		if err != nil {
			return nil, fmt.Errorf("failed to search repos by name: %w", err)
		}
		return fuzzyService.Select(repos)
	default:
		slog.DebugContext(ctx, "listing repos by default github user", "user", envs.DefaultGitHubUser)
		if envs.DefaultGitHubUser == "" {
			return nil, fmt.Errorf("no default GitHub user configured")
		}
		repos, err := svc.ListRepositoriesByUser(ctx, envs.DefaultGitHubUser)
		if err != nil {
			return nil, fmt.Errorf("failed to list repos by authenticated user: %w", err)
		}
		return fuzzyService.Select(repos)
	}
}

func cloneAll(repos []gg.Repo, outDir string, shallow bool) []gg.CloneResult {
	slog.Debug("cloning repos", "outDir", outDir, "repos", repos)

	var wg sync.WaitGroup
	gitService := git.NewService()
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
