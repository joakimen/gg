package root

import (
	"log/slog"
	"os"

	"github.com/joakimen/gg/pkg/cmd/github"
	"github.com/joakimen/gg/pkg/cmd/root/version"
	"github.com/spf13/cobra"
)

type Opts struct {
	Debug bool
}

func NewCmd() *cobra.Command {
	var opts Opts
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

	rootCmd.AddCommand(
		version.NewCmd(),
		github.NewCmd(),
	)

	return rootCmd
}
