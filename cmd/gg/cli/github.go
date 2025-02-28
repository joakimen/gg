package cli

import (
	"cmp"
	"context"
	"fmt"
	"os"

	"github.com/joakimen/gg/fuzzy"
	"github.com/joakimen/gg/git"
	"github.com/joakimen/gg/github"
	"github.com/joakimen/gg/keyring"
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
				return githubLogin(cmd.Context())
			},
		},
		{
			Use:   "logout",
			Short: "Clear stored GitHub credentials",
			RunE: func(_ *cobra.Command, _ []string) error {
				return githubLogout()
			},
		},
		{
			Use:   "show",
			Short: "Show stored GitHub credentials",
			RunE: func(_ *cobra.Command, _ []string) error {
				return githubShow()
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
			return githubClone(
				cmd.Context(),
				flags.owner,
				flags.repo,
				flags.outDir,
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

func githubLogin(ctx context.Context) error {
	keyring := keyring.New(githubKeyring)

	// read user token
	token, err := tty.ReadPassword("Enter your GitHub API token: ")
	if err != nil {
		return err
	}

	if token == "" {
		return fmt.Errorf("the provided token cannot be empty")
	}

	// test user token
	client := github.NewClient(token)
	userLogin, err := client.GetAuthenticatedUser(ctx)
	if err != nil {
		return err
	}

	fmt.Println("Successfully authenticated as user:", userLogin)

	// if token is valid, store it
	err = keyring.Set(token)
	if err != nil {
		return err
	}

	return nil
}

func githubLogout() error {
	fmt.Println("Deleting stored github credentials from system keyring..")
	keyring := keyring.New("github")
	err := keyring.Delete()
	if err != nil {
		return fmt.Errorf("failed to delete github credentals: %w", err)
	}
	fmt.Println("Done.")
	return nil
}

func githubShow() error {
	keyring := keyring.New("github")
	token, err := keyring.Get()
	if err != nil {
		return fmt.Errorf("failed to delete github credentals: %w", err)
	}
	fmt.Println(token)
	return nil
}

func githubClone(
	ctx context.Context,
	owner string,
	repo string,
	outDirFlag string,
	shallow bool,
	repoFile string,
) error {
	var (
		defaultGitHubUser = os.Getenv("GG_GITHUB_USER")
		outDirEnv         = os.Getenv("GG_CLONE_DIR")
		outDir            = cmp.Or(outDirFlag, outDirEnv)
	)
	if outDir == "" {
		return fmt.Errorf("must specify clone directory")
	}

	keyring := keyring.New(githubKeyring)
	token, err := keyring.Get()
	if err != nil {
		return err
	}

	githubClient := github.NewClient(token)
	repoFilter := fuzzy.NewProvider()

	repos, err := githubClient.FindRepos(
		ctx, github.FindRepoOpts{
			RepoFilter:        repoFilter,
			Owner:             owner,
			Repo:              repo,
			OutDir:            outDir,
			Shallow:           shallow,
			RepoFile:          repoFile,
			DefaultGitHubUser: defaultGitHubUser,
		})
	if err != nil {
		return fmt.Errorf("failed to find repos to clone using the provided args: %w", err)
	}

	gitClient := git.NewClient()
	return githubClient.Clone(ctx, gitClient, repos, outDir, shallow)
}
