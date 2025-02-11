package github

import (
	"cmp"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"sync"

	"github.com/joakimen/gg"
	"github.com/joakimen/gg/github"
	"github.com/spf13/cobra"
)

const repoSearchLimit = 100

type flags struct {
	Owner           string
	Repo            string
	RepoFile        string
	IncludeArchived bool
	Limit           int
	CloneDir        string
	Shallow         bool
}

type envs struct {
	CloneDir string
}

func loadEnvs() envs {
	return envs{
		CloneDir: os.Getenv("GG_CLONE_DIR"),
	}
}

func NewCloneCmd() *cobra.Command {
	var opts flags
	cloneCmd := &cobra.Command{
		Use:   "clone",
		Short: "clone a repo interactively",
		RunE: func(_ *cobra.Command, _ []string) error {
			envs := loadEnvs()
			slog.Debug("clone opts", "opts", opts, "envs", envs)
			cloneDir := cmp.Or(opts.CloneDir, envs.CloneDir)
			if cloneDir == "" {
				return fmt.Errorf(
					"must specify clone directory through --clone-dir or by setting the $GG_CLONE_DIR environment variable",
				)
			}

			reposToClone, err := getReposToClone(
				opts.Owner,
				opts.Repo,
				opts.RepoFile,
				opts.IncludeArchived,
				opts.Limit,
			)
			if err != nil {
				return fmt.Errorf("failed to get repos to clone: %w", err)
			}

			if len(reposToClone) == 0 {
				fmt.Println("no repos selected, exiting")
				return nil
			}

			cloneErrors := clone(cloneDir, reposToClone, opts.Shallow)
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

	cloneCmd.Flags().StringVarP(&opts.Owner, "owner", "o", "", "owner of the repo(s) to clone")
	cloneCmd.Flags().StringVarP(&opts.Repo, "repo", "r", "", "owner of the repo(s) to clone")
	cloneCmd.Flags().StringVarP(&opts.RepoFile, "file", "f", "", "name of file containing list of repos to clone")
	cloneCmd.Flags().StringVarP(&opts.CloneDir, "clone-dir", "d", "", "the output directory of cloned repos")
	cloneCmd.Flags().BoolVarP(&opts.IncludeArchived, "include-archived", "a", false, "owner of the repo(s) to clone")
	cloneCmd.Flags().IntVarP(&opts.Limit, "limit", "l", repoSearchLimit, "maximum number of repos to list during search")

	return cloneCmd
}

func clone(cloneDir string, reposToClone []gg.Repo, shallow bool) []gg.CloneResult {
	slog.Debug("cloning repos", "cloneDir", cloneDir, "repos", reposToClone)

	var wg sync.WaitGroup
	resultChan := make(chan gg.CloneResult, len(reposToClone))
	for _, repo := range reposToClone {
		wg.Add(1)
		go func(r gg.Repo) {
			defer wg.Done()
			cloneError := github.Clone(cloneDir, repo, shallow)
			resultChan <- gg.CloneResult{Repo: r, Err: cloneError}
		}(repo)
	}
	wg.Wait()
	close(resultChan)

	slog.Debug("procesing clone results")
	var errs []gg.CloneResult
	for res := range resultChan {
		if res.Err != nil {
			errs = append(errs, res)
		}
	}
	return errs
}

func readReposFromFile(filepath string) ([]gg.Repo, error) {
	var reposFromFile []gg.Repo
	repoJSONData, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read repo file: %w", err)
	}

	if err = json.Unmarshal(repoJSONData, &reposFromFile); err != nil {
		return nil, fmt.Errorf("failed to unmarshal repos from file: %w", err)
	}
	return reposFromFile, nil
}

func getReposToClone(
	owner string,
	repo string,
	repoFile string,
	includeArchived bool,
	repoListLimit int,
) ([]gg.Repo, error) {
	var (
		reposToClone []gg.Repo
		err          error
	)

	switch {
	case repoFile != "":
		slog.Debug("reading repos from file", "file", repoFile)
		reposToClone, err = readReposFromFile(repoFile)
		if err != nil {
			return []gg.Repo{}, fmt.Errorf("couldn't read repos from file: %w", err)
		}
	case owner != "" && repo != "":
		slog.Debug("both owner and repo were provided, cloning single repo", "owner", owner, "repo", repo)
		reposToClone = []gg.Repo{
			{
				Owner: owner,
				Name:  repo,
			},
		}
	default:
		slog.Debug("searching github for repos", "owner", owner, "repo", repo)
		reposToClone, err = github.Search(owner, repo, includeArchived, repoListLimit)
		if err != nil {
			return []gg.Repo{}, fmt.Errorf("failed to search for repos: %w", err)
		}
	}

	return reposToClone, nil
}
