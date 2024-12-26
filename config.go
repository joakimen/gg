package main

import (
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
)

type Config struct {
	Owner           string
	Repo            string
	IncludeArchived bool
	Limit           int
	CloneDir        string
	Verbose         bool
}

type Flags struct {
	Owner           string
	Repo            string
	IncludeArchived bool
	Limit           int
	CloneDir        string
	Verbose         bool
}

type Envs struct {
	CloneDir string
}

func LoadConfig() (Config, error) {
	cfg := Config{
		Limit: 100,
	}

	flags := Flags{}
	flag.StringVar(&flags.Owner, "o", cfg.Owner, "owner of the repository to clone")
	flag.StringVar(&flags.Repo, "r", cfg.Repo, "name of the repository to clone")
	flag.BoolVar(&flags.IncludeArchived, "a", cfg.IncludeArchived, "include archived repositories")
	flag.IntVar(&flags.Limit, "l", cfg.Limit, "limit the number of repositories to search for")
	flag.StringVar(&flags.CloneDir, "d", cfg.CloneDir, "directory to clone the repositories into")
	flag.BoolVar(&flags.Verbose, "v", cfg.Verbose, "verbose output")
	flag.Parse()

	envs := Envs{
		CloneDir: os.Getenv("CLONE_DIR"),
	}

	if flags.Owner != "" {
		cfg.Owner = flags.Owner
	}

	if flags.Repo != "" {
		cfg.Repo = flags.Repo
	}

	if flags.IncludeArchived {
		cfg.IncludeArchived = flags.IncludeArchived
	}

	if flags.Limit != 0 {
		cfg.Limit = flags.Limit
	}

	if envs.CloneDir != "" {
		cfg.CloneDir = envs.CloneDir
	} else if flags.CloneDir != "" {
		cfg.CloneDir = flags.CloneDir
	} else {
		return Config{}, fmt.Errorf("clone directory not specified, set either the CLONE_DIR environment variable or use the -d flag")
	}

	if flags.Verbose {
		cfg.Verbose = flags.Verbose
	}

	return cfg, nil
}

func ConfigureLogger(w io.Writer, cfg Config) *slog.Logger {
	var handlerOpts slog.HandlerOptions
	if cfg.Verbose {
		handlerOpts.Level = slog.LevelInfo
	} else {
		handlerOpts.Level = slog.LevelWarn
	}
	logger := slog.New(slog.NewTextHandler(w, &handlerOpts))
	return logger
}
