package main

import (
	"errors"
	"flag"
	"log/slog"
	"os"
)

const (
	DefaultRepoListLimit = 100
)

type Config struct {
	Owner           string
	Repo            string
	IncludeArchived bool
	Limit           int
	CloneDir        string
	Verbose         bool
	RepoFile        string
	DebugLogging    bool
}

type Flags struct {
	Owner           string
	Repo            string
	IncludeArchived bool
	Limit           int
	CloneDir        string
	Verbose         bool
	RepoFile        string
	DebugLogging    bool
}

type Envs struct {
	CloneDir string
}

func Load(args []string) (Config, error) {

	slog.Debug("loading config", "args", args)

	// initialize cfg with defaults
	cfg := Config{
		Limit: DefaultRepoListLimit,
	}

	// parse flags
	flags := Flags{}
	fs := flag.NewFlagSet("clone", flag.ContinueOnError)
	fs.StringVar(&flags.Owner, "o", cfg.Owner, "owner of the repository to clone")
	fs.StringVar(&flags.Repo, "r", cfg.Repo, "name of the repository to clone")
	fs.BoolVar(&flags.IncludeArchived, "a", cfg.IncludeArchived, "include archived repositories")
	fs.IntVar(&flags.Limit, "l", cfg.Limit, "limit the number of repositories to search for")
	fs.StringVar(&flags.CloneDir, "d", cfg.CloneDir, "directory to clone the repositories into")
	fs.StringVar(&flags.RepoFile, "f", cfg.RepoFile, "File containing the list of repositories to clone")
	fs.BoolVar(&flags.Verbose, "v", cfg.Verbose, "verbose output")
	fs.BoolVar(&flags.DebugLogging, "debug", false, "enable debug logging")
	if err := fs.Parse(args); err != nil {
		return Config{}, err
	}

	// parse env vars
	envs := Envs{
		CloneDir: os.Getenv("CLONE_DIR"),
	}

	// conclude configurations
	cfg.Owner = flags.Owner
	cfg.Repo = flags.Repo
	cfg.IncludeArchived = flags.IncludeArchived
	cfg.Verbose = flags.Verbose
	cfg.RepoFile = flags.RepoFile
	cfg.Limit = flags.Limit
	cfg.DebugLogging = flags.DebugLogging

	switch {
	case flags.CloneDir != "":
		cfg.CloneDir = flags.CloneDir
	case envs.CloneDir != "":
		cfg.CloneDir = envs.CloneDir
	default:
		return Config{},
			errors.New("clone directory not specified, set either the CLONE_DIR environment variable or use the -d flag")
	}

	slog.Debug("done loading config", "config", cfg)

	// return finalized config
	return cfg, nil
}
