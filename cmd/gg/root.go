package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/joakimen/gg/cmd/gg/commands/github"
	"github.com/spf13/cobra"
)

// overwritten with ldflags during ci.
var version = "(development build)"

type RootOpts struct {
	Debug bool
}

func newRootCmd() *cobra.Command {
	var opts RootOpts
	rootCmd := &cobra.Command{
		Use:   "gg",
		Short: "Convenience cli for everyday things",
		PersistentPreRun: func(_ *cobra.Command, _ []string) {
			if opts.Debug {
				logOut := os.Stdout
				logHandlerOpts := &slog.HandlerOptions{
					Level: slog.LevelDebug,
				}
				logHandler := slog.NewTextHandler(logOut, logHandlerOpts)
				logger := slog.New(logHandler)
				slog.SetDefault(logger)
				slog.Debug("debug logging enabled")
			}
		},
	}

	rootCmd.PersistentFlags().BoolVar(&opts.Debug, "debug", false, "Enable debug logging")

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version information",
		RunE: func(_ *cobra.Command, _ []string) error {
			fmt.Println(version)
			return nil
		},
	}

	githubCmd := newGitHubCmd()

	rootCmd.AddCommand(
		versionCmd,
		githubCmd,
	)

	return rootCmd
}

func newGitHubCmd() *cobra.Command {
	githubCmd := &cobra.Command{
		Use:   "github",
		Short: "Convenience wrapper for github stuff",
	}

	githubCmd.AddCommand(
		github.NewLoginCmd(),
		github.NewShowCmd(),
		github.NewLogoutCmd(),
		github.NewCloneCmd(),
	)

	return githubCmd
}
