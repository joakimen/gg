package clone

import (
	"cmp"
	"context"
	"fmt"
	"os"

	"github.com/joakimen/gg"
)

type Command struct {
	GitHubClient gg.GitHubClient
	Fuzzy        gg.RepoFuzzyProvider
	GitClient    gg.GitClient
}

type Flags struct {
	Owner           string
	Repo            string
	OutDir          string
	Shallow         bool
	RepoFile        string
	IncludeArchived bool
}

func (c *Command) Run(ctx context.Context, flags Flags) error {
	var (
		defaultGitHubUser = os.Getenv("GG_GITHUB_USER")
		outDirEnv         = os.Getenv("GG_CLONE_DIR")
		outDir            = cmp.Or(flags.OutDir, outDirEnv)
	)
	if outDir == "" {
		return fmt.Errorf("must specify clone directory")
	}

	repos, err := c.GitHubClient.FindRepos(
		ctx, gg.FindRepoOpts{
			RepoFilter:        c.Fuzzy,
			Owner:             flags.Owner,
			Repo:              flags.Repo,
			OutDir:            outDir,
			Shallow:           flags.Shallow,
			RepoFile:          flags.RepoFile,
			DefaultGitHubUser: defaultGitHubUser,
			IncludeArchived:   flags.IncludeArchived,
		})
	if err != nil {
		return fmt.Errorf("failed to find repos to clone using the provided args: %w", err)
	}

	return c.GitHubClient.Clone(ctx, c.GitClient, repos, outDir, flags.Shallow)
}
