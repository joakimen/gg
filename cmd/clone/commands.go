package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/joakimen/clone/github"
)

func (cmd VersionCmd) Run(_ *Globals) error {
	fmt.Println(version)
	return nil
}

func (cmd *CloneCmd) Run(glb *Globals) error {
	reposToClone, err := getReposToClone(
		cmd.Owner,
		cmd.Repo,
		cmd.RepoFile,
		cmd.IncludeArchived,
		cmd.Limit,
	)
	if err != nil {
		return fmt.Errorf("failed to get repos to clone: %w", err)
	}

	if len(reposToClone) == 0 {
		fmt.Println("no repos selected, exiting")
		return nil
	}

	cloneErrors := clone(cmd.CloneDir, reposToClone, cmd.Shallow)
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

func getReposToClone(
	owner string,
	repo string,
	repoFile string,
	includeArchived bool,
	repoListLimit int,
) ([]github.Repo, error) {
	var (
		reposToClone []github.Repo
		err          error
	)
	switch {
	case repoFile != "":
		slog.Debug("reading repos from file", "file", repoFile)
		reposToClone, err = readReposFromFile(repoFile)
		if err != nil {
			return []github.Repo{}, fmt.Errorf("couldn't read repos from file: %w", err)
		}
	case owner != "" && repo != "":
		slog.Debug("both owner and repo were provided, cloning single repo", "owner", owner, "repo", repo)
		reposToClone = []github.Repo{
			{
				Owner: owner,
				Name:  repo,
			},
		}
	default:
		slog.Debug("searching github for repos", "owner", owner, "repo", repo)
		reposToClone, err = github.Search(owner, repo, includeArchived, repoListLimit)
		if err != nil {
			return []github.Repo{}, fmt.Errorf("failed to search for repos: %w", err)
		}
	}
	return reposToClone, nil
}
