package cli

import (
	"cmp"
	"os"

	"github.com/joakimen/gg/github"
	"github.com/spf13/cobra"
)

func newGitHubCmd() *cobra.Command {
	githubCmd := &cobra.Command{
		Use:   "github",
		Short: "Convenience wrapper for github stuff",
	}

	subcommands := []*cobra.Command{
		{
			Use:   "login",
			Short: "Authenticate to GitHub",
			RunE: func(cmd *cobra.Command, _ []string) error {
				return github.Login(cmd.Context())
			},
		},
		{
			Use:   "logout",
			Short: "Clear stored GitHub credentials",
			RunE: func(_ *cobra.Command, _ []string) error {
				return github.Logout()
			},
		},
		{
			Use:   "show",
			Short: "Show stored GitHub credentials",
			RunE: func(_ *cobra.Command, _ []string) error {
				return github.Show()
			},
		},
		newGitHubCloneCmd(),
	}

	githubCmd.AddCommand(subcommands...)

	return githubCmd
}

func newGitHubCloneCmd() *cobra.Command {
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
			var envs struct {
				defaultGitHubUser string
				outDir            string
			}

			envs.defaultGitHubUser = os.Getenv("GG_GITHUB_USER")
			envs.outDir = os.Getenv("GG_CLONE_DIR")
			outDir := cmp.Or(flags.outDir, envs.outDir)

			return github.Clone(
				cmd.Context(),
				outDir,
				flags.owner,
				flags.repo,
				envs.defaultGitHubUser,
				flags.shallow,
				flags.repoFile,
			)
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
