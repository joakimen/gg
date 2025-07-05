package cli

import (
	"github.com/joakimen/gg"
	"github.com/joakimen/gg/github"
	"github.com/spf13/cobra"
)

func newGitHubLoginCmd(gh *github.Service) *cobra.Command {
	return &cobra.Command{
		Use:   "login",
		Short: "Authenticate to GitHub",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return gh.Login(cmd.Context())
		},
	}
}

func newGitHubLogoutCmd(gh *github.Service) *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Clear stored GitHub credentials",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return gh.Logout(cmd.Context())
		},
	}
}

func newGitHubShowCmd(gh *github.Service) *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Show stored GitHub credentials",

		RunE: func(cmd *cobra.Command, _ []string) error {
			return gh.Show(cmd.Context())
		},
	}
}

// TODO: refactor/split to separate finding and cloning repos
func newGitHubCloneCmd(gh *github.Service) *cobra.Command {
	var flags struct {
		owner           string
		repo            string
		repoFile        string
		outDir          string
		includeArchived bool
		shallow         bool
	}
	cmd := &cobra.Command{
		Use:   "clone",
		Short: "Clone GitHub repos interactively",
		RunE: func(cmd *cobra.Command, _ []string) error {
			cloneFlags := gg.CloneFlags{
				Owner:           flags.owner,
				Repo:            flags.repo,
				OutDir:          flags.outDir,
				Shallow:         flags.shallow,
				RepoFile:        flags.repoFile,
				IncludeArchived: flags.includeArchived,
			}

			return gh.Clone(cmd.Context(), cloneFlags)
		},
	}

	cmd.Flags().StringVarP(&flags.owner, "owner", "o", "", "owner of the repo(s) to clone")
	cmd.Flags().StringVarP(&flags.repo, "repo", "r", "", "owner of the repo(s) to clone")
	cmd.Flags().StringVarP(&flags.repoFile, "file", "f", "", "name of file containing list of repos to clone")
	cmd.Flags().StringVarP(&flags.outDir, "out-dir", "d", "", "the output directory of cloned repos")
	cmd.Flags().BoolVarP(&flags.includeArchived, "include-archived", "a", false, "owner of the repo(s) to clone")
	cmd.Flags().BoolVarP(&flags.shallow, "shallow", "", false, "perform a shallow clone")

	return cmd
}
