package cli

import (
	"context"
	"log/slog"
	"os"

	"github.com/joakimen/gg/fuzzy"
	"github.com/joakimen/gg/git"
	"github.com/joakimen/gg/github"
	"github.com/joakimen/gg/githubapi"
	"github.com/joakimen/gg/keyring"
	"github.com/joakimen/gg/prompt"
	"github.com/spf13/cobra"
)

func Run() error {
	rootCmd := newRootCmd()
	return rootCmd.ExecuteContext(context.Background())
}

func newRootCmd() *cobra.Command {
	githubService := github.NewService(
		keyring.New("github"),
		prompt.ReadPassword,
		githubapi.TokenClientProvider,
		git.NewClient(),
		fuzzy.NewProvider(),
	)

	var opts struct {
		Debug bool
	}

	rootCmd := &cobra.Command{
		Use:          "gg",
		Short:        "Convenience cli for everyday things",
		SilenceUsage: true,
		PersistentPreRun: func(_ *cobra.Command, _ []string) {
			if opts.Debug {
				enableDebugLogging()
			}
		},
	}

	rootCmd.PersistentFlags().BoolVar(&opts.Debug, "debug", false, "Enable debug logging")

	githubCmd := &cobra.Command{
		Use:   "github",
		Short: "Convenience wrapper for github stuff",
	}

	githubCmd.AddCommand(
		newGitHubLoginCmd(githubService),
		newGitHubLogoutCmd(githubService),
		newGitHubShowCmd(githubService),
		newGitHubCloneCmd(githubService),
	)

	versionCmd := newVersionCmd()

	rootCmd.AddCommand(
		githubCmd,
		versionCmd,
	)

	return rootCmd
}

func enableDebugLogging() {
	logOut := os.Stdout
	logHandlerOpts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	logHandler := slog.NewTextHandler(logOut, logHandlerOpts)
	logger := slog.New(logHandler)
	slog.SetDefault(logger)
	slog.Debug("debug logging enabled")
}
