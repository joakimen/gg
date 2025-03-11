package cli

import (
	"context"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

func Run() error {
	rootCmd := newRootCmd()
	return rootCmd.ExecuteContext(context.Background())
}

func newRootCmd() *cobra.Command {
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

	rootCmd.AddCommand(
		newVersionCmd(),
		newGitHubCmd(),
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
