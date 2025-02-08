package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/alecthomas/kong"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cli := CLI{
		Globals: Globals{
			Debug: false,
		},
	}
	ctx := kong.Parse(&cli,
		kong.Name(cliName),
		kong.Description(cliDesc),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}))

	if cli.Globals.Debug {
		logHandlerOpts := &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}
		logHandler := slog.NewTextHandler(os.Stderr, logHandlerOpts)
		logger := slog.New(logHandler)
		slog.SetDefault(logger)
	}

	slog.Debug("done parsing args", "args", cli)

	err := ctx.Run(&cli.Globals)
	ctx.FatalIfErrorf(err)
	return err
}
