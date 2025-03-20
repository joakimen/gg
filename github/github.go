package github

import (
	"cmp"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/joakimen/gg"
)

type Service struct {
	Keyring        gg.KeyringProvider
	TTY            gg.TTYProvider
	ClientProvider gg.GitHubClientProvider
	Git            gg.GitClient
	RepoSelector   gg.RepoSelector
}

func NewService(
	keyring gg.KeyringProvider,
	ttyProvider gg.TTYProvider,
	clientProvider gg.GitHubClientProvider,
	gitClient gg.GitClient,
	repoFilter gg.RepoSelector,
) *Service {
	return &Service{
		Keyring:        keyring,
		TTY:            ttyProvider,
		RepoSelector:   repoFilter,
		ClientProvider: clientProvider,
		Git:            gitClient,
	}
}

func (s *Service) Login(ctx context.Context) error {
	// read api token from user
	token, err := s.TTY.Read("Enter your GitHub API token: ")
	if err != nil {
		return err
	}

	if strings.TrimSpace(token) == "" {
		return fmt.Errorf("the provided token cannot be empty")
	}

	// test user token
	client := s.ClientProvider(token)
	userLogin, err := client.GetAuthenticatedUser(ctx)
	if err != nil {
		return err
	}

	fmt.Println("Successfully authenticated as user:", userLogin)

	// if token is valid, store it
	err = s.Keyring.Set(token)
	if err != nil {
		return err
	}
	fmt.Println("Token successfully stored in system keyring.")
	return nil
}

func (s *Service) Logout(ctx context.Context) error {
	err := s.Keyring.Delete()
	if err != nil {
		return err
	}
	fmt.Println("Token successfully deleted from system keyring.")
	return nil
}

func (s *Service) Show(ctx context.Context) error {
	token, err := s.Keyring.Get()
	if err != nil {
		return err
	}
	fmt.Println(token)
	return nil
}

type CloneFlags struct {
	Owner           string
	Repo            string
	OutDir          string
	Shallow         bool
	RepoFile        string
	IncludeArchived bool
}

func (s *Service) Clone(ctx context.Context, flags CloneFlags) error {
	var (
		defaultGitHubUser = os.Getenv("GG_GITHUB_USER")
		outDirEnv         = os.Getenv("GG_CLONE_DIR")
		outDir            = cmp.Or(flags.OutDir, outDirEnv)
	)
	if outDir == "" {
		return fmt.Errorf("must specify clone directory")
	}

	token, err := s.Keyring.Get()
	if err != nil {
		return err
	}

	client := s.ClientProvider(token)
	repos, err := client.FindRepos(
		ctx, gg.FindRepoOpts{
			RepoSelector:      s.RepoSelector,
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

	return client.Clone(ctx, s.Git, repos, outDir, flags.Shallow)
}
