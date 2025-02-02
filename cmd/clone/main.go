package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
)

var (
	version = "(development build)"
)

func main() {
	if err := (&Main{}).Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "failed to run:", err)
		os.Exit(1)
	}
}

func (m *Main) Run(args []string) error {
	cfg, err := Load(args)

	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if cfg.DebugLogging {
		logHandlerOpts := &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}
		logHandler := slog.NewTextHandler(os.Stderr, logHandlerOpts)
		logger := slog.New(logHandler)
		slog.SetDefault(logger)
	}

	if len(args) > 0 {
		switch args[0] {
		case "help":
			fmt.Fprintln(os.Stderr, flag.ErrHelp)
			m.Usage()
			return nil
		case "version":
			fmt.Println(version)
			return nil
		}
	}
	return (&CloneCommand{}).Run(cfg)
}

type Main struct{}

func (m *Main) Usage() {
	fmt.Fprintln(os.Stderr, `
clone is a tool for interactively cloning GitHub repositories.

Usage:

	clone <command> [arguments]

The commands are:

	version      prints the binary version
	help         display this help screen
`[1:])
}
