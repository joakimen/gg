package main

import (
	"flag"
	"fmt"
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

func (m *Main) Run(osArgs []string) error {
	var cmd string
	var args []string
	if len(osArgs) > 0 {
		cmd, args = osArgs[0], osArgs[1:]
	}

	outWriter := os.Stderr
	cfg, err := Load(args, outWriter)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	info := cfg.InfoFn

	switch cmd {
	case "help":
		info(flag.ErrHelp)
		m.Usage()
		return nil
	case "version":
		return (&VersionCommand{}).Run(cfg)
	default:
		return (&CloneCommand{}).Run(cfg)
	}
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
