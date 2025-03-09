package cli

import (
	"fmt"

	"github.com/joakimen/gg/fuzzy"
	"github.com/joakimen/gg/git"
	"github.com/joakimen/gg/github"
	"github.com/joakimen/gg/keyring"
	"github.com/joakimen/gg/pkg/commands/github/clone"
	"github.com/joakimen/gg/pkg/commands/github/login"
	"github.com/joakimen/gg/pkg/commands/github/logout"
	"github.com/joakimen/gg/pkg/commands/github/show"
	"github.com/joakimen/gg/tty"
	"github.com/spf13/cobra"
)

const githubKeyring = "github"

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
				loginCmd := &login.Command{
					Keyring:       keyring.New(githubKeyring),
					ClientFactory: github.ClientFactory,
					TTY:           tty.NewProvider(),
				}
				return loginCmd.Run(cmd.Context())
			},
		},
		{
			Use:   "logout",
			Short: "Clear stored GitHub credentials",
			RunE: func(cmd *cobra.Command, _ []string) error {
				logoutCmd := &logout.Command{
					Keyring: keyring.New(githubKeyring),
				}
				return logoutCmd.Run(cmd.Context())
			},
		},
		{
			Use:   "show",
			Short: "Show stored GitHub credentials",

			RunE: func(cmd *cobra.Command, _ []string) error {
				keyring := keyring.New(githubKeyring)
				showCmd := &show.Command{
					Keyring: keyring,
				}
				return showCmd.Run(cmd.Context())
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
			keyring := keyring.New(githubKeyring)
			token, err := keyring.Get()
			if err != nil {
				return fmt.Errorf("failed to get token: %w", err)
			}
			cloneCmd := &clone.Command{
				GitHubClient: github.NewClient(token),
				GitClient:    git.NewClient(),
				Fuzzy:        fuzzy.NewProvider(),
			}

			cloneFlags := clone.Flags{
				Owner:           flags.owner,
				Repo:            flags.repo,
				OutDir:          flags.outDir,
				Shallow:         flags.shallow,
				RepoFile:        flags.repoFile,
				IncludeArchived: flags.includeArchived,
			}
			return cloneCmd.Run(cmd.Context(), cloneFlags)
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
